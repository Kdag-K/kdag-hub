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

// JSONGenesisFile is the structure that a Genesis file gets parsed into.
type JSONGenesisFile struct {
	Alloc      *Alloc `json:"alloc"`
	Poa        *POA   `json:"poa"`
	Controller *POA   `json:"controller"`
}

// MinimalPeerRecord is used where only an Address and Moniker are required.
// The standard Peer datatypes us PubKeyHex not address.
type MinimalPeerRecord struct {
	Address string
	Moniker string
}
