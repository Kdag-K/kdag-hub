package keys

import (
	"fmt"
	
	"github.com/Kdag-K/kdag-hub/src/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var _newPasswordFile string

// newUpdateCmd returns the command that changes the passphrase of a keyfile
func newUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [moniker]",
		Short: "change the passphrase on a keyfile",
		Long: `
Change the passphrase on a keyfile.

If --passfile is not specified, the user will be prompted to enter the current
passphrase manually. Likewise, if --new-passfile is not specified, the user will
be prompted to input and confirm the new password.
		`,
		Args: cobra.ExactArgs(1),
		RunE: update,
	}
	
	addUpdateFlags(cmd)
	
	return cmd
}

func addUpdateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&_newPasswordFile, "new-passfile", "", "the file containing the new passphrase")
	viper.BindPFlags(cmd.Flags())
}

func update(cmd *cobra.Command, args []string) error {
	moniker := args[0]
	
	err := crypto.UpdateKeyByMoniker(_keystore, moniker, _passwordFile, _newPasswordFile)
	if err != nil {
		fmt.Println(err)
	}
	
	return nil
}