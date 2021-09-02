package crypto

import "github.com/ethereum/go-ethereum/accounts/keystore"

const ethereumVersion = 3

type encryptedKeyJSONV3 struct {
	Address string              `json:"address"`
	Crypto  keystore.CryptoJSON `json:"crypto"`
	ID      string              `json:"id"`
	Version int                 `json:"version"`
}

// EncryptedKeyJSONKnode is an extension of a regular Ethereum keyfile with an
// added public key. It makes our lives easier when working with Babble. We
// could change the Version number, but then other non-Knode tools, would not be
// able to decrypt keys
type EncryptedKeyJSONKnode struct {
	Address   string              `json:"address"`
	PublicKey string              `json:"pub"`
	Crypto    keystore.CryptoJSON `json:"crypto"`
	ID        string              `json:"id"`
	Version   int                 `json:"version"`
}
