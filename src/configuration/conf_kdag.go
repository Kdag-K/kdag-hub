package configuration

import (
	"fmt"
	"time"

	"github.com/Kdag-K/kdag-hub/src/common"
)

var (
	defaultNodeAddr        = fmt.Sprintf("%s:%d", common.GetNodeIP(), 1337)
	defaultHeartbeat       = 200 * time.Millisecond
	defaultTCPTimeout      = 1000 * time.Millisecond
	defaultCacheSize       = 50000
	defaultSyncLimit       = 1000
	defaultBootstrap       = true
	defaultMaxPool         = 2
	defaultMaintenanceMode = false
	defaultSuspendLimit    = 300
)

// KdagConf contains the configuration for the Kdag node used by knode.
// It only presents a subset of the options Kdag can accept, because knode
// forces some configurations values. In particular, the --fast-sync and
// --store flags are disabled because knode does not support the FastSync
// protocol, and it requires a persistant database.
type KdagConf struct {

	// BindAddr is the local address:port where this node gossips with other
	// nodes. By default, this is "0.0.0.0", meaning Kdag will bind to all
	// addresses on the local machine and will advertise the private IPv4
	// address to the rest of the cluster. However, in some cases, there may be
	// a routable address that cannot be bound. Use AdvertiseAddr to enable
	// gossiping a different address to support this. If this address is not
	// routable, the node will be in a constant flapping state as other nodes
	// will treat the non-routability as a failure
	BindAddr string `mapstructure:"listen"`

	// AdvertiseAddr is used to change the address that we advertise to other
	// nodes in the cluster
	AdvertiseAddr string `mapstructure:"advertise"`

	// Gossip heartbeat
	Heartbeat time.Duration `mapstructure:"heartbeat"`

	// TCP timeout
	TCPTimeout time.Duration `mapstructure:"timeout"`

	// Max number of items in caches
	CacheSize int `mapstructure:"cache-size"`

	// Max number of Event in SyncResponse
	SyncLimit int `mapstructure:"sync-limit"`

	// Max number of connections in net pool
	MaxPool int `mapstructure:"max-pool"`

	// Bootstrap from database
	Bootstrap bool `mapstructure:"bootstrap"`

	// MaintenanceMode when set to true causes Kdag to initialise in a
	// suspended state. I.e. it does not start gossipping
	MaintenanceMode bool `mapstructure:"maintenance-mode"`

	// SuspendLimit is the number of undetermined-events produced since the last
	// run, that will cause the node to be automaitically suspended.
	SuspendLimit int `mapstructure:"suspend-limit"`

	// Moniker is a friendly name to indentify this peer
	Moniker string `mapstructure:"moniker"`
}

// DefaultBabbleConfig returns the default configuration for a Kdag node
func DefaultBabbleConfig() *KdagConf {
	return &KdagConf{
		BindAddr:        defaultNodeAddr,
		Heartbeat:       defaultHeartbeat,
		TCPTimeout:      defaultTCPTimeout,
		CacheSize:       defaultCacheSize,
		SyncLimit:       defaultSyncLimit,
		MaxPool:         defaultMaxPool,
		Bootstrap:       defaultBootstrap,
		MaintenanceMode: defaultMaintenanceMode,
		SuspendLimit:    defaultSuspendLimit,
	}
}
