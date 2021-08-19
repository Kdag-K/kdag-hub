package configuration

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// Directory Constants
const (
	// config
	ConfigDir = "knode-config"
	KdagDir   = "kdag"
	EthDir    = "eth"
	POADir    = "poa"

	// data
	DatabaseDir = "knode-data"

	// keystore
	KeyStoreDir = "keystore"
)

// Knode Configuration Directory
const (
	KnodeTomlDirDot  = ".knode"
	KnodeTomlDirCaps = "Knode"
)

// Filename constants
const (
	PeersJSON        = "peers.json"
	PeersGenesisJSON = "peers.genesis.json"
	GenesisJSON      = "genesis.json"
	KnodeTomlFile    = "knode.toml"
	EthDB            = "eth-db"
	KdagDB           = "kdag-db"
	WalletTomlFile   = "wallet.toml"
	ServerPIDFile    = "server.pid"
	KdagPrivKey      = "priv_key"
)

// Network Constants
const (
	DefaultGossipPort = "1337"
	DefaultAPIAddr    = ":8080"
)

//Keys constants
const (
	DefaultKeyfile        = "keyfile.json"
	DefaultPrivateKeyFile = "priv_key"
)

// Genesis Constants
const (
	DefaultAccountBalance            = "1234567890000000000000"
	DefaultContractAddress           = "abbaabbaabbaabbaabbaabbaabbaabbaabbaabba"
	DefaultControllerContractAddress = "aabbaabbaabbaabbaabbaabbaabbaabbaabbaabb"
	GenesisContract                  = "contract0.sol"
	GenesisABI                       = "contract0.abi"
	ControllerContract               = "contract1.sol"
	ControllerABI                    = "contract1.abi"
	CompileResultFile                = "compile.toml"
)

// DefaultConfigDir returns the full path of the config directory where static
// configuration files are stored.
func DefaultConfigDir() string {
	return filepath.Join(DefaultKnodeDir(), ConfigDir)
}

// DefaultDataDir returns the full path of the data directory where databases
// are stored.
func DefaultDataDir() string {
	return filepath.Join(DefaultKnodeDir(), DatabaseDir)
}

// DefaultKeystoreDir returns the full path of the keystore where encrypted
// keyfiles are stored.
func DefaultKeystoreDir() string {
	return filepath.Join(DefaultKnodeDir(), "keystore")
}

// DefaultKnodeDir returns a the full path for the default location Knode
// configuration files based on the underlying OS.
func DefaultKnodeDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "Knode")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "Knode")
		} else {
			return filepath.Join(home, ".knode")
		}
	}
	return ""
}

// Guess a sensible default location from OS and environment variables.
func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}
