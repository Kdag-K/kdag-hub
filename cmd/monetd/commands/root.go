//Package commands implements the CLI commands for monetd
package commands

import (
	"fmt"

	"github.com/Kdag-K/kdag-hub/src/common"

	"github.com/Kdag-K/kdag-hub/cmd/monetd/commands/config"
	"github.com/Kdag-K/kdag-hub/cmd/monetd/commands/keys"
	"github.com/Kdag-K/kdag-hub/src/configuration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

/*******************************************************************************
RootCmd
*******************************************************************************/

//RootCmd is the root command for monetd
var RootCmd = &cobra.Command{
	Use:   "monetd",
	Short: "monet daemon",
	Long: `
Monetd is the daemon component of the Monet Toolchain; a distributed
smart-contract platform based on the Ethereum Virtual Machine and Babble 
consensus.  
	
See the documentation at https://docs.monet.network/ for further information.
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
	RootCmd.AddCommand(
		keys.KeysCmd,
		config.ConfigCmd,
		newRunCmd(),
		versionCmd,
	)

	// set global flags
	RootCmd.PersistentFlags().BoolP("verbose", "v", configuration.Global.Verbose, "verbose output")

	// do not print usage when error occurs
	RootCmd.SilenceUsage = true
}

/*******************************************************************************
HELPERS
*******************************************************************************/

// Read config into Viper. CLI flags have precedence over the toml file.
func readConfig(cmd *cobra.Command) error {

	// Register flags with viper
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	// first unmarshal to read from cli flags
	if err := viper.Unmarshal(configuration.Global); err != nil {
		return err
	}

	// Read from configuration file if there is one.
	viper.SetConfigName("monetd")                       // name of config file (without extension)
	viper.AddConfigPath(configuration.Global.ConfigDir) // search config directory

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		common.DebugMessage(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		common.DebugMessage(fmt.Sprintf("No config file monetd.toml found in %s\n", configuration.Global.DataDir))
	} else {
		return err
	}

	// second unmarshal to read from config file
	return viper.Unmarshal(configuration.Global)
}
