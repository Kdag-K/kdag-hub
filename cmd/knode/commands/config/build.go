package config

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	
	"github.com/Kdag-K/evm/src/keystore"
	"github.com/Kdag-K/kdag-hub/src/common"
	"github.com/Kdag-K/kdag-hub/src/configuration"
	"github.com/Kdag-K/kdag-hub/src/files"
	"github.com/Kdag-K/kdag-hub/src/genesis"
	"github.com/Kdag-K/kdag/src/peers"
	eth_crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addBuildFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&_keystore, "keystore", _keystore, "keystore directory")
	cmd.Flags().StringVar(&_configDir, "config", _configDir, "output directory")
	cmd.Flags().StringVar(&_addressParam, "address", _addressParam, "IP/hostname of this node")
	cmd.Flags().StringVar(&_passwordFile, "passfile", "", "file containing the passphrase")
	viper.BindPFlags(cmd.Flags())
}


// dumpPeers takes a list of peers and dumps it into peers.json and
// peers.genesis.json in the kdag directory
func dumpPeers(configDir string, peers []*peers.Peer) error {
	peersJSONOut, err := json.MarshalIndent(peers, "", "\t")
	if err != nil {
		return err
	}
	
	// peers.json
	jsonFileName := filepath.Join(configDir, configuration.KdagDir, configuration.PeersJSON)
	files.WriteToFile(jsonFileName, string(peersJSONOut), files.OverwriteSilently)
	
	// peers.genesis.json
	jsonFileName = filepath.Join(configDir, configuration.KdagDir, configuration.PeersGenesisJSON)
	files.WriteToFile(jsonFileName, string(peersJSONOut), files.OverwriteSilently)
	
	return nil
}


func buildConfig(cmd *cobra.Command, args []string) error {
	moniker := args[0]
	
	address := fmt.Sprintf("%s:%s", _addressParam, configuration.DefaultGossipPort)
	
	// Some debug output confirming parameters
	common.DebugMessage("Building Config for   : ", moniker)
	common.DebugMessage("Using Network Address : ", address)
	common.DebugMessage("Using Password File   : ", _passwordFile)
	
	// set global config moniker
	configuration.Global.Kdag.Moniker = moniker
	
	// Retrieve the keyfile corresponding to moniker
	privateKey, err := keystore.GetKey(_keystore, moniker, _passwordFile)
	if err != nil {
		return err
	}
	
	// Create Directories if they don't exist
	CreateKnodeConfigFolders(_configDir)
	
	// Copy the key to kdag directory with appropriate permissions
	err = keystore.DumpPrivKey(
		filepath.Join(_configDir, configuration.KdagDir),
		privateKey)
	if err != nil {
		return err
	}
	
	pubKey := hex.EncodeToString(eth_crypto.FromECDSAPub(&privateKey.PublicKey))
	
	// Create a peer-set whith a single node
	peers := []*peers.Peer{
		peers.NewPeer(pubKey, address, moniker),
	}
	
	// Write peers.json and peers.genesis.json
	if err := dumpPeers(_configDir, peers); err != nil {
		return err
	}
	
	// Create the eth/genesis.json file
	err = genesis.GenerateGenesisJSON(
		filepath.Join(_configDir, configuration.EthDir),
		_keystore,
		peers,
		nil,
		configuration.DefaultContractAddress,
		configuration.DefaultControllerContractAddress)
	if err != nil {
		return err
	}
	
	// Write TOML file for knode based on global config object
	err = configuration.DumpGlobalTOML(
		_configDir,
		configuration.KnodeTomlFile,
		true)
	if err != nil {
		return err
	}
	
	return nil
}