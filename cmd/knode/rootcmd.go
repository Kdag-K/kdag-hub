package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		return nil
	},
}

func init() {
	//
}

// Read config into Viper. CLI flags have precedence over the toml file.
func readConfig(cmd *cobra.Command) error {
	// Register flags with viper
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err

	}

	// first unmarshal to read from cli flags
	if err := viper.Unmarshal("knode"); err != nil {
		return err
	}

	// Read from configuration file if there is one.
	viper.SetConfigName("knode") // name of config file (without extension)
	viper.AddConfigPath("knode") // search config directory

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		//
	} else {
		return err
	}

	// second unmarshal to read from config file
	return viper.Unmarshal("knode")
}
