package crypto

import (
	"encoding/hex"
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
)

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

// EncryptKey encrypts a key using the specified scrypt parameters into a json
// blob that can be decrypted later on.
func EncryptKey(key *keystore.Key, auth string, scryptN, scryptP int) ([]byte, error) {
	keyBytes := math.PaddedBigBytes(key.PrivateKey.D, 32)
	cryptoStruct, err := keystore.EncryptDataV3(keyBytes, []byte(auth), scryptN, scryptP)
	if err != nil {
		return nil, err
	}
	encryptedKeyJSONKnode := EncryptedKeyJSONKnode{
		hex.EncodeToString(key.Address[:]),
		hex.EncodeToString(crypto.FromECDSAPub(&key.PrivateKey.PublicKey)),
		cryptoStruct,
		key.Id.String(),
		ethereumVersion,
	}
	return json.Marshal(encryptedKeyJSONKnode)
}
