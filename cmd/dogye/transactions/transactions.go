package transactions

import (
	"github.com/Kdag-K/kdag-hub/cmd/dogye/configuration"
	knodeconfig "github.com/Kdag-K/kdag-hub/src/configuration"
	"github.com/spf13/cobra"
)
var (
	_keystore = knodeconfig.DefaultKeystoreDir()
	_dogye  = configuration.DogyeConfigDir
)

//TODO duplicates the definition in networks package.
//Probably better to publish them and use them directly.
const (
	dogyeNetworksDir     = "networks"
	dogyeKeystoreDir     = "keystore"
	dogyeTransactionsDir = "trans"
	networkTomlFileName    = "network.toml"
)