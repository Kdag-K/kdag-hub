package commands

import (
	"github.com/Kdag-K/kdag-hub/cmd/dogye/commands/keys"
	"github.com/Kdag-K/kdag-hub/cmd/dogye/commands/network"
	"github.com/Kdag-K/kdag-hub/cmd/dogye/commands/parse"
	"github.com/Kdag-K/kdag-hub/cmd/dogye/commands/transactions"
	"github.com/Kdag-K/kdag-hub/src/common"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//RootCmd is the root command for dogye
var RootCmd = &cobra.Command{
	Use:   "dogye",
	Short: "Dogye",
	Long: `Dogye
	
Dogye is the swiss army knife of advanced tools for the Monet Hub. For most users,
you should not need to use this command. The inbuild commands in monetd will suffice for
most use cases.`,
}

func init() {
	
	RootCmd.AddCommand(
		// keys.KeysCmd,
		// network.NetworkCmd,
		// transactions.TransCmd,
		// parse.ParseCmd,
		VersionCmd,
	)
	//do not print usage when error occurs
	RootCmd.SilenceUsage = true
	
	RootCmd.PersistentFlags().BoolVarP(&common.VerboseLogging, "verbose", "v", false, "verbose messages")
	
	viper.BindPFlags(RootCmd.Flags())
}