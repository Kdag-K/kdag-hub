package network

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CLI flags.
var forceNetwork = false
var useExisting = false
var startNodes = false

type copyRecord struct {
	from string
	to   string
}

func newStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [network]",
		Short: "start a docker network",
		Long: `
dogye network start

Starts a network. Does not start individual nodes.
		`,
		Args: cobra.ExactArgs(1),
		RunE: networkStart,
	}
	
	addStartFlags(cmd)
	
	return cmd
}

func addStartFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&forceNetwork, "force-network", forceNetwork, "force network down if already exists")
	cmd.Flags().BoolVar(&useExisting, "use-existing", useExisting, "use existing network if already exists")
	cmd.Flags().BoolVar(&startNodes, "start-nodes", startNodes, "start nodes")
	viper.BindPFlags(cmd.Flags())
}

func networkStart(cmd *cobra.Command, args []string) error {
	return nil
}