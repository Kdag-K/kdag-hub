package config

import (
	"github.com/Kdag-K/kdag-hub/src/common"
	"github.com/Kdag-K/kdag-hub/src/configuration"
)

var (
	_keystore     = configuration.DefaultKeystoreDir()
	_configDir    = configuration.DefaultConfigDir()
	_keyParam     = ""
	_addressParam = common.GetNodeIP()
	_passwordFile string
	_force        = false
)