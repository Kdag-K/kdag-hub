package config

import (
	"github.com/Kdag-K/kdag-hub/src/common"
	"github.com/Kdag-K/kdag-hub/src/configuration"
	"github.com/spf13/cobra"
)

var (
	_keystore     = configuration.DefaultKeystoreDir()
	_configDir    = configuration.DefaultConfigDir()
	_keyParam     = ""
	_addressParam = common.GetNodeIP()
	_passwordFile string
	_force        = false
)


// ConfigCmd implements the config CLI subcommand.
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "manage configuration",
	Long: `
Manage Knode configuration.

* config build - creates the configuration for a single-node network, based on
                 one of the keys in <keystore>. This is a quick and easy way to
                 get started with Knode.

* config pull -  fetches the configuration from a running node. This is used to
                 join an existing network.

For more complex scenarios, please refer to 'giverny', which is a specialised
Knode configuration tool.
`,
	TraverseChildren: true,
}