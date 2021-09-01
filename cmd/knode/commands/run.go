package commands

import (
	"fmt"
	evers "github.com/Kdag-K/evm/src/version
	"github.com/Kdag-K/kdag-hub/src/genesis"
	"github.com/Kdag-K/evm/src/engine"
	"github.com/Kdag-K/kdag-hub/src/common"
	"github.com/Kdag-K/kdag-hub/src/configuration"
	"github.com/spf13/cobra"
)

// newRunCmd returns the command that starts the daemon.
func newRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run a node",
		Long:  `Run a node.`,

		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			common.DebugMessage(fmt.Sprintf("Base Config: %+v", configuration.Global.BaseConfig))
			common.DebugMessage(fmt.Sprintf("Kdag Config: %+v", configuration.Global.Kdag))
			common.DebugMessage(fmt.Sprintf("Eth Config: %+v", configuration.Global.Eth))
			return nil
		},

		RunE: runKnode,
	}

	bindFlags(cmd)

	return cmd
}

func bindFlags(cmd *cobra.Command) {
	// Config and data directories
	cmd.Flags().StringP("config", "c", configuration.Global.ConfigDir, "configuration directory")
	cmd.Flags().StringP("data", "d", configuration.Global.DataDir, "data directory")

	// EVM and Kdag share the same API address
	cmd.Flags().String("api-listen", configuration.Global.APIAddr, "IP:PORT of HTTP API service")

	// Kdag config
	cmd.Flags().String("kdag.listen", configuration.Global.Kdag.BindAddr, "bind IP:PORT of Kdag node")
	cmd.Flags().String("kdag.advertise", configuration.Global.Kdag.AdvertiseAddr, "advertise IP:PORT of Kdag node")
	cmd.Flags().Duration("kdag.heartbeat", configuration.Global.Kdag.Heartbeat, "heartbeat timer milliseconds (time between gossips)")
	cmd.Flags().Duration("kdag.timeout", configuration.Global.Kdag.TCPTimeout, "TCP timeout milliseconds")
	cmd.Flags().Int("kdag.cache-size", configuration.Global.Kdag.CacheSize, "number of items in LRU caches")
	cmd.Flags().Int("kdag.sync-limit", configuration.Global.Kdag.SyncLimit, "max number of Events per sync")
	cmd.Flags().Int("kdag.max-pool", configuration.Global.Kdag.MaxPool, "max number of pool connections")
	cmd.Flags().Bool("kdag.bootstrap", configuration.Global.Kdag.Bootstrap, "bootstrap Kdag from database")
	cmd.Flags().String("kdag.moniker", configuration.Global.Kdag.Moniker, "friendly name")
	cmd.Flags().Bool("kdag.maintenance-mode", configuration.Global.Kdag.MaintenanceMode, "start kdag in suspended (non-gossipping) state")
	cmd.Flags().Int("kdag.suspend-limit", configuration.Global.Kdag.SuspendLimit, "number of undetermined-events since last run that will trigger automatic suspension")

	// Eth config
	cmd.Flags().Int("eth.cache", configuration.Global.Eth.Cache, "megabytes of memory allocated to internal caching (min 16MB / database forced)")
	cmd.Flags().String("eth.min-gas-price", configuration.Global.Eth.MinGasPrice,
		"minimum gasprice of transactions submitted through this node (ex 1K, 1M, 1G, etc.)")
}

// Run the EVM / Kdag  engine
func runKnode(cmd *cobra.Command, args []string) error {
	return nil
}