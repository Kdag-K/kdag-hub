package network

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Kdag-K/kdag-hub/cmd/dogye/configuration"
	"github.com/Kdag-K/kdag-hub/src/common"
	mconfiguration "github.com/Kdag-K/kdag-hub/src/configuration"
	"github.com/Kdag-K/kdag-hub/src/files"

	"github.com/spf13/cobra"
)

func newLocationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "location [network_name]",
		Short: "show the location of the configuration files",
		Long: `
dogye network location
		`,
		Args: cobra.ArbitraryArgs,
		RunE: networkLocation,
	}

	return cmd
}

func networkLocation(cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		fmt.Println(configuration.DogyeConfigDir)
		return nil
	}

	networkName = strings.TrimSpace(args[0])

	if !common.CheckMoniker(networkName) {
		return errors.New("the network name, " + networkName + ", is invalid")
	}

	if !files.CheckIfExists(filepath.Join(configuration.DogyeConfigDir, dogyeNetworksDir, networkName)) {
		return errors.New("the network, " + networkName + " has not been created")
	}

	common.InfoMessage("Network                 : " + networkName)

	common.InfoMessage("Dogye Config Dir      : " + configuration.DogyeConfigDir)
	common.InfoMessage("Dogye Networks Dir    : " +
		filepath.Join(configuration.DogyeConfigDir, dogyeNetworksDir, networkName))
	common.InfoMessage("Dogye KeyStore Dir    : " +
		filepath.Join(configuration.DogyeConfigDir, dogyeNetworksDir, networkName, dogyeKeystoreDir))
	common.InfoMessage("Peers JSON              : " +
		filepath.Join(configuration.DogyeConfigDir, dogyeNetworksDir, networkName, mconfiguration.PeersJSON))
	common.InfoMessage("Peers Genesis JSON      : " +
		filepath.Join(configuration.DogyeConfigDir, dogyeNetworksDir, networkName, mconfiguration.PeersGenesisJSON))
	common.InfoMessage("Genesis JSON            : " +
		filepath.Join(configuration.DogyeConfigDir, dogyeNetworksDir, networkName, mconfiguration.GenesisJSON))
	common.InfoMessage("Knode TOML             : " +
		filepath.Join(configuration.DogyeConfigDir, dogyeNetworksDir, networkName, mconfiguration.KnodeTomlFile))
	common.InfoMessage("Network TOML            : " +
		filepath.Join(configuration.DogyeConfigDir, dogyeNetworksDir, networkName, networkTomlFileName))

	return nil
}
