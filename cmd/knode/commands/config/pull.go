package config

import (
	"fmt"
	"path/filepath"
	
	"github.com/Kdag-K/evm/src/keystore"
	"github.com/Kdag-K/kdag-hub/src/common"
	"github.com/Kdag-K/kdag-hub/src/configuration"
	"github.com/Kdag-K/kdag-hub/src/files"
	"github.com/spf13/cobra"
)

func pullConfig(cmd *cobra.Command, args []string) error {
	peerAddr := args[0]
	
	var err error
	
	if _keyParam == "" {
		if _keyParam, err = getDefaultKey(_keystore); err != nil {
			return err
		}
	}
	
	common.InfoMessage("Moniker: ", _keyParam)
	
	// Helpful debugging output
	common.MessageWithType(common.MsgDebug, "Pulling from         : ", peerAddr)
	common.MessageWithType(common.MsgDebug, "Using Network Address: ", _addressParam)
	common.MessageWithType(common.MsgDebug, "Using Key            : ", _keyParam)
	common.MessageWithType(common.MsgDebug, "Using Password File  : ", _passwordFile)
	
	// Retrieve the keyfile corresponding to moniker
	privateKey, err := keystore.GetKey(_keystore, _keyParam, _passwordFile)
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
	
	rootURL := "http://" + peerAddr
	
	filesList := []*downloadItem{
		{URL: rootURL + "/genesispeers",
			Dest: filepath.Join(_configDir, configuration.KdagDir, configuration.PeersGenesisJSON)},
		{URL: rootURL + "/peers",
			Dest: filepath.Join(_configDir, configuration.KdagDir, configuration.PeersJSON)},
		{URL: rootURL + "/genesis",
			Dest: filepath.Join(_configDir, configuration.EthDir, configuration.GenesisJSON)},
	}
	
	for _, item := range filesList {
		err := files.DownloadFile(
			item.URL,
			item.Dest,
			!_force)
		if err != nil {
			common.ErrorMessage(fmt.Sprintf("Error downloading %s", item.URL))
			return err
		}
		common.DebugMessage("Downloaded ", item.Dest)
	}
	
	// Write TOML file for Knode based on global config object
	err = configuration.DumpGlobalTOML(
		_configDir,
		configuration.KnodeTomlFile,
		!_force)
	if err != nil {
		return err
	}
	
	ShowIPWarnings()
	
	return nil
}