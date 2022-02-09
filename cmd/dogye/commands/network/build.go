package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/Kdag-K/kdag-hub/src/genesis"
	bpeers "github.com/Kdag-K/kdag-hub/src/peers"

	"github.com/Kdag-K/kdag-hub/cmd/dogye/configuration"
	"github.com/Kdag-K/kdag-hub/src/common"
	monetconfig "github.com/Kdag-K/kdag-hub/src/configuration"
	"github.com/Kdag-K/kdag-hub/src/files"
	"github.com/spf13/cobra"
)

func newBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [network_name]",
		Short: "build knode configuration files based on a dogye network file",
		Args:  cobra.ExactArgs(1),
		RunE:  networkBuild,
	}
	return cmd
}

func networkBuild(cmd *cobra.Command, args []string) error {
	return buildNetwork(strings.TrimSpace(args[0]))
}

// buildNetwork builds the network. It is called directly from the "new" command
// as well.
func buildNetwork(networkName string) error {
	if !common.CheckMoniker(networkName) {
		return errors.New("network name, " + networkName + ", is invalid")
	}

	// Check all the files and directories we expect actually exist
	thisNetworkDir := filepath.Join(configuration.DogyeConfigDir, dogyeNetworksDir, networkName)
	if !files.CheckIfExists(thisNetworkDir) {
		return errors.New("cannot find the configuration folder, " + thisNetworkDir + " for " + networkName)
	}

	networkTomlFile := filepath.Join(thisNetworkDir, networkTomlFileName)
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
		return fmt.Errorf("Failed to parse the toml file at '%s': %v", networkTomlFile, err)
	}

	common.DebugMessage("Building network " + networkName)

	err = generateKnodeConfig(&conf, thisNetworkDir)
	if err != nil {
		common.ErrorMessage("Error writing peers json file")
		return err
	}

	return nil
}

func generateKnodeConfig(conf *Config, thisNetworkDir string) error {
	var peers []*bpeers.Peer
	var alloc = make(genesis.Alloc)

	for _, n := range conf.Nodes {

		netaddr := n.NetAddr
		if !strings.Contains(netaddr, ":") {
			netaddr += ":" + monetconfig.DefaultGossipPort
		}

		rec := genesis.AllocRecord{Moniker: n.Moniker, Balance: n.Tokens}
		alloc[n.Address] = &rec

		if !n.Validator || n.NonNode {
			continue
		}

		peers = append(peers, bpeers.NewPeer(n.PubKeyHex, netaddr, n.Moniker))
	}

	err := generateKdagFiles(thisNetworkDir, peers)
	if err != nil {
		return err
	}

	err = genesis.GenerateGenesisJSON(thisNetworkDir,
		"",
		peers,
		&alloc,
		monetconfig.DefaultContractAddress,
		monetconfig.DefaultControllerContractAddress,
	)
	if err != nil {
		return err
	}

	return err
}

func generateKdagFiles(configDir string, peers []*bpeers.Peer) error {
	peersJSONOut, err := json.MarshalIndent(peers, "", "\t")
	if err != nil {
		return err
	}

	// write peers.json
	jsonFileName := filepath.Join(configDir, monetconfig.PeersJSON)
	err = files.WriteToFile(jsonFileName, string(peersJSONOut), files.OverwriteSilently)
	if err != nil {
		return err
	}

	// Write peers.genesis.json
	jsonFileName = filepath.Join(configDir, monetconfig.PeersGenesisJSON)
	err = files.WriteToFile(jsonFileName, string(peersJSONOut), files.OverwriteSilently)
	if err != nil {
		return err
	}
	return nil
}
