package network

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/Kdag-K/kdag-hub/cmd/dogye/configuration"
	"github.com/Kdag-K/kdag-hub/src/common"
	knodeconfig "github.com/Kdag-K/kdag-hub/src/configuration"
	"github.com/Kdag-K/kdag-hub/src/files"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

func newAWSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aws [network] [output-path]",
		Short: "write aws configuration files",
		Long: `
dogye network aws

Writes AWS configuration.
		`,
		Args: cobra.ExactArgs(2),
		RunE: networkAWS,
	}

	return cmd
}

func networkAWS(cmd *cobra.Command, args []string) error {
	network := args[0]
	outPath := args[1]

	if !files.CheckIfExists(outPath) {
		return errors.New("cannot find the output folder, " + outPath + " for " + network)
	}

	if err := buildNetworkConfig(network, outPath); err != nil {
		return err
	}

	return nil
}

func buildNetworkConfig(networkName string, outPath string) error {

	// Set some paths
	thisNetworkDir := filepath.Join(configuration.DogyeConfigDir, dogyeNetworksDir, networkName)
	networkTomlFile := filepath.Join(thisNetworkDir, networkTomlFileName)

	// Check expected config exists
	if !files.CheckIfExists(thisNetworkDir) {
		return errors.New("cannot find the configuration folder, " + thisNetworkDir + " for " + networkName)
	}

	if !files.CheckIfExists(networkTomlFile) {
		return errors.New("cannot find the configuration file: " + networkTomlFile)
	}

	var conf = Config{}

	tomlbytes, err := ioutil.ReadFile(networkTomlFile)
	if err != nil {
		return fmt.Errorf("Failed to read the toml file at '%s': %v", networkTomlFile, err)
	}

	err = toml.Unmarshal(tomlbytes, &conf)
	if err != nil {
		return nil
	}

	common.DebugMessage("Configuring Network ", conf.Docker.Name)

	err = exportAWSConfigs(&conf, outPath)
	if err != nil {
		return err
	}

	return nil
}

func exportAWSConfigs(conf *Config, outPath string) error {

	// Configure some paths
	networkDir := filepath.Join(configuration.DogyeConfigDir, dogyeNetworksDir, conf.Network.Name)
	err := files.CreateDirsIfNotExists([]string{outPath})
	if err != nil {
		return err
	}

	for _, n := range conf.Nodes { // loop around nodes
		if !n.NonNode {
			if err := exportAWSNodeConfig(networkDir, outPath, &n); err != nil {
				return err
			}
		}
	}

	return nil
}

func exportAWSNodeConfig(networkDir, outPath string, n *node) error {

	netaddr := n.NetAddr
	if !strings.Contains(netaddr, ":") {
		netaddr += ":" + knodeconfig.DefaultGossipPort
	}
	// Build output files

	if n.Moniker != "" { // Should not be blank here, but safety first

		knodeDir := filepath.Join(outPath, n.Moniker)
		configDir := filepath.Join(knodeDir, knodeconfig.ConfigDir)
		knodeConfigDir := filepath.Join(configDir, knodeconfig.KdagDir)
		ethConfigDir := filepath.Join(configDir, knodeconfig.EthDir)
		keystoreDir := filepath.Join(knodeDir, knodeconfig.KeyStoreDir)

		common.DebugMessage("Creating config in " + configDir)

		err := files.CreateDirsIfNotExists([]string{
			knodeConfigDir,
			ethConfigDir,
			keystoreDir,
		})
		if err != nil {
			return err
		}
		//copying record.

		// Write a node description file containing all of the parameters needed
		// to start a container. Saves having to load and parse network.toml for
		//  every node
		nodeConfigFile := filepath.Join(outPath, n.Moniker+".toml")
		nodeConfig := dockerNodeConfig{
			Moniker: n.Moniker,
			NetAddr: strings.Split(netaddr, ":")[0],
		}

		tomlBytes, err := toml.Marshal(nodeConfig)
		if err != nil {
			return err
		}

		err = files.WriteToFile(nodeConfigFile, string(tomlBytes), 0)
		if err != nil {
			return err
		}

		// edit knode.toml and set knode.listen appropriately.

	}
	return nil
}
