package transactions

import "math/big"

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