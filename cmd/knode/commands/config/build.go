package config

import (
	"encoding/json"
	"path/filepath"
	
	"github.com/Kdag-K/kdag-hub/src/configuration"
	"github.com/Kdag-K/kdag-hub/src/files"
	"github.com/Kdag-K/kdag/src/peers"
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
// peers.genesis.json in the babble directory
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