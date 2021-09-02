package genesis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/Kdag-K/kdag-hub/src/configuration"
	"github.com/Kdag-K/kdag-hub/src/crypto"
)

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

// buildAlloc builds the alloc structure of the genesis file
func buildAlloc(accountsDir string) (Alloc, error) {
	var alloc = make(Alloc)

	tfiles, err := ioutil.ReadDir(accountsDir)
	if err != nil {
		return alloc, err
	}

	for _, f := range tfiles {
		if filepath.Ext(f.Name()) != ".json" {
			continue
		}

		path := filepath.Join(accountsDir, f.Name())

		// Read key from file.
		keyjson, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("Failed to read the keyfile at '%s': %v", path, err)
		}

		k := new(crypto.EncryptedKeyJSONKnode)
		if err := json.Unmarshal(keyjson, k); err != nil {
			return nil, err
		}

		moniker := strings.TrimSuffix(f.Name(), ".json")
		balance := configuration.DefaultAccountBalance
		addr := k.Address

		rec := AllocRecord{Moniker: moniker, Balance: balance}
		alloc[addr] = &rec
	}

	return alloc, nil
}
