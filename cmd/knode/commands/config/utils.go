package config
import (
	"fmt"
	"path/filepath"
	
	"github.com/Kdag-K/kdag-hub/src/common"
	"github.com/Kdag-K/kdag-hub/src/configuration"
	"github.com/Kdag-K/kdag-hub/src/files"
)

// CreateMonetConfigFolders creates the standard directory layout for a monet
// configuration folder
func CreateMonetConfigFolders(configDir string) error {
	return files.CreateDirsIfNotExists([]string{
		configDir,
		filepath.Join(configDir, configuration.KdagDir),
		filepath.Join(configDir, configuration.EthDir),
		filepath.Join(configDir, configuration.EthDir, configuration.POADir),
	})
}