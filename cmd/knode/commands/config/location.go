package config

import (
	"fmt"
	"path/filepath"
	
	"github.com/Kdag-K/kdag-hub/src/configuration"
	"github.com/spf13/cobra"
)

// newLocationCmd shows the config file path
func newLocationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "location",
		Short: "show default configuration files",
		Long:  "Show the default locations of knode configuration files.",
		RunE:  locationConfig,
	}
	
	return cmd
}

func locationConfig(cmd *cobra.Command, args []string) error {
	
	fmt.Println("Knode Config        : " + filepath.Join(configuration.DefaultConfigDir(), configuration.KnodeTomlFile))
	
	fmt.Println("Kdag Peers         : " + filepath.Join(configuration.DefaultConfigDir(), configuration.KdagDir, configuration.PeersJSON))
	fmt.Println("Kdag Genesis Peers : " + filepath.Join(configuration.DefaultConfigDir(), configuration.KdagDir, configuration.PeersGenesisJSON))
	fmt.Println("Kdag Private Key   : " + filepath.Join(configuration.DefaultConfigDir(), configuration.KdagDir, configuration.KdagPrivKey))
	
	fmt.Println("EVM-Lite Genesis     : " + filepath.Join(configuration.DefaultConfigDir(), configuration.EthDir, configuration.GenesisJSON))
	
	fmt.Println("Kdag Database      : " + filepath.Join(configuration.DefaultDataDir(), configuration.KdagDB))
	fmt.Println("EVM Database    : " + filepath.Join(configuration.DefaultDataDir(), configuration.EthDB))
	
	fmt.Println("Keystore        : " + configuration.DefaultKeystoreDir())
	
	return nil
}