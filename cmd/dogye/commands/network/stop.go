package network

import (
	"github.com/Kdag-K/kdag-hub/src/docker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var stopAndDelete = false

func addStopFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&stopAndDelete, "remove", stopAndDelete, "stop and remove node")
	viper.BindPFlags(cmd.Flags())
}

func networkStop(cmd *cobra.Command, args []string) error {

	if len(args) == 1 { // Network
		return docker.StopNetwork(args[0], stopAndDelete)
	}

	return docker.StopNode(args[0], args[1], stopAndDelete)

}
