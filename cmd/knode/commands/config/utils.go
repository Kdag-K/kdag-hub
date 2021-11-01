package config

import (
	"fmt"
	"path/filepath"
	
	"github.com/Kdag-K/kdag-hub/src/common"
	"github.com/Kdag-K/kdag-hub/src/configuration"
	"github.com/Kdag-K/kdag-hub/src/files"
)

// CreateKnodeConfigFolders creates the standard directory layout for a knode
// configuration folder
func CreateKnodeConfigFolders(configDir string) error {
	return files.CreateDirsIfNotExists([]string{
		configDir,
		filepath.Join(configDir, configuration.KdagDir),
		filepath.Join(configDir, configuration.EthDir),
		filepath.Join(configDir, configuration.EthDir, configuration.POADir),
	})
}

// ShowIPWarnings outputs warnings if IP addresses are local and propably not
// reachable from the outside.
func ShowIPWarnings() {
	api := configuration.Global.APIAddr
	listen := configuration.Global.Kdag.BindAddr
	advertise := configuration.Global.Kdag.AdvertiseAddr
	
	if common.CheckIP(api, true) {
		common.MessageWithType(common.MsgWarning, fmt.Sprintf("Knode service API address in knode.toml may be internal: %s", api))
	}
	
	if advertise != "" && common.CheckIP(advertise, false) {
		common.MessageWithType(common.MsgWarning, fmt.Sprintf("kdag.advertise address in knode.toml may be internal: %s \n", listen))
	} else if common.CheckIP(listen, false) {
		common.MessageWithType(
			common.MsgWarning,
			fmt.Sprintf("kdag.listen address in knode.toml may be internal: %s. Consider setting an advertise address.", listen),
		)
	}
}