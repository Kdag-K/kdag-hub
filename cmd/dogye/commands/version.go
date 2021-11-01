package commands

import (
	"fmt"
	
	knode "github.com/Kdag-K/kdag-hub/src/version"
	"github.com/spf13/cobra"
)

// VersionCmd displays the version of evml being used
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version info",
	Long: `Giverny Version information
	
The version command outputs the version number for Knode, EVM-Lite, Kdag and
Geth.

If you compile your own tools, the suffices are the GIT branch and the GIT
commit hash.
`,
	Run: func(cmd *cobra.Command, args []string) {
		
		fmt.Print(knode.FullVersion())
	},
}