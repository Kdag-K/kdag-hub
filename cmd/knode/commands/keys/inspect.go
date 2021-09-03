package keys

import (
	"fmt"
	
	"github.com/Kdag-K/kdag-hub/src/crypto"
	"github.com/spf13/cobra"
)

func inspect(cmd *cobra.Command, args []string) error {
	moniker := args[0]
	
	err := crypto.InspectKeyByMoniker(_keystore, moniker, _passwordFile, _showPrivate, _outputJSON)
	if err != nil {
		fmt.Println(err)
	}
	
	return nil