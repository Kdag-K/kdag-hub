package transactions

import (
	"math/big"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type soloAccount struct {
	Moniker      string
	Address      string
	NextNonce    int
	Credits      *big.Int
	Debits       *big.Int
	Delta        *big.Int
	Transactions []soloTransaction
}

type soloTransaction struct {
	To     string
	ToName string
	Nonce  int
	Amount *big.Int
}

//CLI params
var accounts string
var outputfile = "trans.json"
var maxTransValue = 10
var roundRobin = false


func addSoloFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&faucet, "faucet", faucet, "faucet account moniker")
	cmd.Flags().StringVar(&accounts, "accounts", accounts, "comma separated account list")
	cmd.Flags().StringVar(&outputfile, "output", outputfile, "output file")
	
	cmd.Flags().BoolVar(&roundRobin, "round-robin", roundRobin, "set sender accounts round robin")
	
	cmd.Flags().IntVar(&totalTransactions, "count", totalTransactions, "number of tranactions to solo")
	cmd.Flags().IntVar(&surplusCredit, "surplus", surplusCredit, "additional credit to allocate each account from the faucet above the bare minimum")
	cmd.Flags().IntVar(&maxTransValue, "max-trans-value", maxTransValue, "maximum transaction value")

	viper.BindPFlags(cmd.Flags())
}