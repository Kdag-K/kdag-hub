package commands

import (
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/sf13/cobra"
)

// rootcmd is the root command for knode.
var RootCmd = &cobra.Command{
	Use:   "knode",
	Short: "knode daemon",
	Long: `
knode is the component of the kdag-hub Toolchain; a distributed
smart-contract platform based on the Ethereum Virtual Machine and Kdag
consensus.

See the documentation at kdag official website for further information.
`,
	TraverseChildren: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := readConfig(cmd); err != nil {
			return err
		}

		if configuration.Global.Verbose {
			common.VerboseLogging = true
		}

		return nil
	},
}

func init() {
	//
}

// Read config into Viper. CLI flags have precedence over the toml file.
func readConfig(cmd *cobra.Command) error {
	//
}
