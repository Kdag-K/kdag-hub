package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// vercmd displays the version of evml being used
var vercmd = &cobra.Command{
	Use:   "version",
	Short: "show version info",
	Long: `knode Version information
           The version command outputs the version number for kdag-hub, EVM, kdag and Geth.

           If you compile your own tools, the suffices are the GIT branch and the GIT commit hash.`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Print("knode version : xxx")
	},
}
