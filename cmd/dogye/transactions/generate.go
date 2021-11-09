package transactions

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addGenerateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&networkName, "network", "n", "", "network name")
	cmd.Flags().StringVar(&ips, "ips", "", "ips.dat file path")
	cmd.Flags().StringVar(&faucet, "faucet", faucet, "faucet account moniker")
	cmd.Flags().IntVar(&totalTransactions, "count", totalTransactions, "number of tranactions to generate")
	cmd.Flags().IntVar(&surplusCredit, "surplus", surplusCredit, "additional credit to allocate each account from the faucet above the bare minimum")
	viper.BindPFlags(cmd.Flags())
}