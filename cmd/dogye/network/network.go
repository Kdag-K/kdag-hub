package network

import (
	"path/filepath"
	
	"github.com/Kdag-K/kdag-hub/cmd/dogye/configuration"
	"github.com/Kdag-K/kdag-hub/src/files"
)

const (
	dogyeNetworksDir  = "networks"
	dogyeKeystoreDir  = "keystore"
	dogyeDockerDir    = "docker"
	dogyeTmpDir       = ".tmp"
	defaultTokens       = "1234567890000000000000"
	networkTomlFileName = "network.toml"
)

var (
	networkName = "network0"
)


func createdogyeRootNetworkFolders() error {
	
	files.CreateDirsIfNotExists([]string{
		configuration.DogyeConfigDir,
		filepath.Join(configuration.DogyeConfigDir, dogyeNetworksDir),
	})
	
	return nil
}