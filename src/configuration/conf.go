// Package configuration holds shared configuration structs for kdag-hub & knode, EVM and Kdag.
package configuration

import (
	"path/filepath"

	evm_conf "github.com/Kdag-K/evm/src/config"
	"github.com/Kdag-K/kdag-hub/src/common"
	kdag_conf "github.com/Kdag-K/kdag/src/config"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	// Global is a global Config object used by commands in cmd/ to manipulate
	// configuration options.
	Global = DefaultConfig()
)


// Config contains the configuration for Knode node.
type Config struct {

	// Top level options use an anonymous struct
	BaseConfig `mapstructure:",squash"`

	// Options for evm
	Eth *EthConf `mapstructure:"eth"`

	// Options for Kdag
	Kdag *KdagConf `mapstructure:"kdag"`
}

// DefaultConfig returns the default configuration for a Knode node.
func DefaultConfig() *Config {
	return &Config{
		BaseConfig: DefaultBaseConfig(),
		Eth:        DefaultEthConfig(),
		Kdag:       DefaultKdagConfig(),
	}
}

// LogLevel returns a logrus-style log-level based on the verbose option.
func (c *Config) LogLevel() string {
	if c.Verbose {
		return "debug"
	}
	return "info"
}

// Logger returns a new prefixed logrus Entry with custom formatting
func (c *Config) Logger(prefix string) *logrus.Entry {
	if c.logger == nil {
		c.logger = logrus.New()
		c.logger.Level = common.LogLevel(c.LogLevel())
		c.logger.Formatter = new(prefixed.TextFormatter)
	}
	return c.logger.WithField("prefix", prefix)
}

// ToEVMConfig extracts evm configuration and returns a config object as
// used by the evm library.
func (c *Config) ToEVMConfig() *evm_conf.Config {
	evmConf := evm_conf.DefaultConfig()

	evmConf.LogLevel = c.LogLevel()
	evmConf.EthAPIAddr = c.APIAddr
	evmConf.Genesis = filepath.Join(c.ConfigDir, EthDir, GenesisJSON)
	evmConf.DbFile = filepath.Join(c.DataDir, EthDB)
	evmConf.Cache = c.Eth.Cache
	evmConf.MinGasPrice = c.Eth.MinGasPrice

	return evmConf
}

// ToKdagConfig extracts the kdag configuration and returns a config object
// as used by the Kdag library. It enforces the values of Store and
// EnableFastSync to true and false respectively.
func (c *Config) ToKdagConfig() *kdag_conf.Config {
	kdagConfig := kdag_conf.NewDefaultConf()
  kdagConfig.DataDir = filepath.Join(c.ConfigDir, KdagDir)
	kdagConfig.DatabaseDir = filepath.Join(c.DataDir, KdagDB)
	kdagConfig.LogLevel = c.LogLevel()
	kdagConfig.BindAddr = c.Kdag.BindAddr
	kdagConfig.AdvertiseAddr = c.Kdag.AdvertiseAddr
	kdagConfig.MaxPool = c.Kdag.MaxPool
	kdagConfig.HeartbeatTimeout = c.Kdag.Heartbeat
	kdagConfig.TCPTimeout = c.Kdag.TCPTimeout
	kdagConfig.CacheSize = c.Kdag.CacheSize
	kdagConfig.SyncLimit = c.Kdag.SyncLimit
	kdagConfig.Bootstrap = c.Kdag.Bootstrap
	kdagConfig.Moniker = c.Kdag.Moniker
	kdagConfig.MaintenanceMode = c.Kdag.MaintenanceMode
	kdagConfig.SuspendLimit = c.Kdag.SuspendLimit

	// Force Kdag to use persistant storage.
	kdagConfig.Store = true

	// Force FastSync = false because evm does not support Snapshot/Restore
	// yet.
	kdagConfig.EnableFastSync = false

	// An empty ServiceAddr tells Kdag not to start an API server. The API
	// handlers are still registered with the DefaultServeMux, so they will be
	// served by the evm API server automatically.
	kdagConfig.ServiceAddr = ""

	return kdagConfig
}

/*******************************************************************************
BASE CONFIG
*******************************************************************************/

// BaseConfig contains the top level configuration for an EVM-Kdag node
type BaseConfig struct {
	// ConfigDir contains static configuration files
	ConfigDir string `mapstructure:"config"`

	// DataDir contains kdag and eth databases
	DataDir string `mapstructure:"data"`

	// Verbose
	Verbose bool `mapstructure:"verbose"`

	// IP/PORT of API
	APIAddr string `mapstructure:"api-listen"`

	logger *logrus.Logger
}

// DefaultBaseConfig returns the default top-level configuration for EVM-Kdag
func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		ConfigDir: DefaultConfigDir(),
		DataDir:   DefaultDataDir(),
		APIAddr:   DefaultAPIAddr,
	}
}
