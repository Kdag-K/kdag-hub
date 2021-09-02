package genesis

// AllocRecord is an object that contains information about a pre-funded acount.
type AllocRecord struct {
	Balance string `json:"balance"`
	Moniker string `json:"moniker"`
}

// Alloc is the section of a genesis file that contains the list of pre-funded
// accounts.
type Alloc map[string]*AllocRecord

// POA is the section of a genesis file that contains information about
// the POA smart-contract.
type POA struct {
	Address string            `json:"address"`
	Abi     string            `json:"abi"`
	Code    string            `json:"code"`
	Storage map[string]string `json:"storage,omitempty"`
}
