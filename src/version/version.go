//Package version provides version information for the application
package version

import (
	"fmt"
	
	geth "github.com/ethereum/go-ethereum/params"
	kdag  "github.com/Kdag-K/kdag/src/version"
	evm "github.com/Kdag-K/evm/src/version"
)

var (
	//Version is the full version string
	Version = "0.1.1"
	
	// GitCommit is set with --ldflags "-X main.gitCommit=$(git rev-parse HEAD)"
	GitCommit string
	// GitBranch is set with --ldflags "-X main.gitBranch=$(git symbolic-ref --short HEAD)"
	GitBranch string
)

func init() {
	// branch is only of interest if it is not the master branch
	if GitBranch != "" && GitBranch != "master" {
		Version += "-" + GitBranch
	}
	
	if GitCommit != "" {
		Version += "-" + GitCommit[:8]
	}
}

//FullVersion outputs version information for Knode, EVM, Kdag and Geth
func FullVersion() string {
	return fmt.Sprintln("knode Version: "+Version) +
		fmt.Sprintln("     EVM Version: "+evm.Version) +
		fmt.Sprintln("     Kdag Version: "+kdag.Version) +
		fmt.Sprintln("     Geth Version: "+geth.Version)
}