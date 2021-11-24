package transactions

import (
	"github.com/Kdag-K/kdag-hub/cmd/dogye/configuration"
	knodeconfig "github.com/Kdag-K/kdag-hub/src/configuration"
	"github.com/spf13/cobra"
)
var (
	_keystore = knodeconfig.DefaultKeystoreDir()
	_dogye  = configuration.DogyeConfigDir
)

//TODO duplicates the definition in networks package.
//Probably better to publish them and use them directly.
const (
	dogyeNetworksDir     = "networks"
	dogyeKeystoreDir     = "keystore"
	dogyeTransactionsDir = "trans"
	networkTomlFileName    = "network.toml"
)

//TransCmd implements the transactions subcommand
var TransCmd = &cobra.Command{
	Use:   "transactions",
	Short: "dogye transactions",
	Long: `Server
	
The dogye transaction command is used to generate sets of transactions for
testing networks.`,
	
	TraverseChildren: true,
}

func init() {
	//Subcommands
	TransCmd.AddCommand(
		newGenerateCmd(),
		newSoloCmd(),
	)
	
	TransCmd.PersistentFlags().StringVarP(&_keystore, "keystore", "k", _keystore, "keystore directory")
	TransCmd.PersistentFlags().StringVarP(&_dogye, "dir", "d", _dogye, "dogye directory")
	
}