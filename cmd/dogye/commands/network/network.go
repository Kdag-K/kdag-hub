package network

import (
	"path/filepath"
	
	"github.com/Kdag-K/kdag-hub/cmd/dogye/configuration"
	"github.com/Kdag-K/kdag-hub/src/files"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	dogyeNetworksDir    = "networks"
	dogyeKeystoreDir    = "keystore"
	dogyeDockerDir      = "docker"
	dogyeTmpDir         = ".tmp"
	defaultTokens       = "1234567890000000000000"
	networkTomlFileName = "network.toml"
)

var (
	networkName = "network0"
)

// NetworkCmd is the CLI subcommand.
var NetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Advanced Network Configuration",
	Long: `Network
	
Advanced Network Config Manager. `,
	
	TraverseChildren: true,
}

func init() {
	
	// Subcommands
	NetworkCmd.AddCommand(
		newAWSCmd(),
		newStartCmd(),
		newPushCmd(),
		newStatusCmd(),
		newStopCmd(),
		newListCmd(),
		newAWSCmd(),
		newDumpCmd(),
		newLocationCmd(),
	)
	
	viper.BindPFlags(NetworkCmd.Flags())
	
	// make sure the giverny config folders exist.
	createdogyeRootNetworkFolders()
}

func createdogyeRootNetworkFolders() error {
	
	files.CreateDirsIfNotExists([]string{
		configuration.DogyeConfigDir,
		filepath.Join(configuration.DogyeConfigDir, dogyeNetworksDir),
	})
	
	return nil
}