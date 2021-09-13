package keys

import (
	"fmt"
	
	"github.com/Kdag-K/kdag-hub/src/crypto"
	"github.com/spf13/cobra"
)

// newNewCmd returns the command that creates a new keypair
func newNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [moniker]",
		Short: "create a new keyfile",
		Long: `
Generate a new key identified by [moniker].

The keyfile will be written to <keystore>/[moniker].json. If the --passfile flag
is not specified, the user will be prompted to enter the passphrase manually.
`,
		Args: cobra.ExactArgs(1),
		RunE: newKey,
	}
	
	return cmd
}

func newKey(cmd *cobra.Command, args []string) error {
	moniker := args[0]
	
	// key is returned, but we don't want to do anything with it.
	_, err := crypto.NewKeyfile(_keystore, moniker, _passwordFile)
	if err != nil {
		fmt.Println(err)
	}
	
	return nil
}