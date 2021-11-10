package transactions
import "math/big"

// CLI Parameters
var networkName = ""
var ips = ""
var faucet = "Faucet"
var totalTransactions = 20
var surplusCredit = 1000000

type ipmapping map[string]string

type account struct {
	Moniker   string
	Tokens    string
	PubKeyHex string
	Address   string
}

type transaction struct {
	From   int
	To     int
	Amount *big.Int
}