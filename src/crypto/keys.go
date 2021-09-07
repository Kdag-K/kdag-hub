package crypto

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	
	"github.com/Kdag-K/evm/src/crypto"
	"github.com/Kdag-K/kdag-hub/src/common"
	"github.com/Kdag-K/kdag-hub/src/files"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	eth_crypto "github.com/ethereum/go-ethereum/crypto"
)

type outputGenerate struct {
	Address      string
	AddressEIP55 string
}

type outputInspect struct {
	Address    string
	PublicKey  string
	PrivateKey string
}

// InspectKey inspects an encrypted keyfile
func InspectKey(keyfilepath string, PasswordFile string, showPrivate bool, outputJSON bool) error {

	// Read key from file.
	keyjson, err := ioutil.ReadFile(keyfilepath)
	if err != nil {
		return fmt.Errorf("Failed to read the keyfile at '%s': %v", keyfilepath, err)
	}

	// Decrypt key with passphrase.
	passphrase, err := crypto.GetPassphrase(PasswordFile, false)
	if err != nil {
		return err
	}

	key, err := keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		return fmt.Errorf("Error decrypting key: %v", err)
	}

	// Output all relevant information we can retrieve.
	out := outputInspect{
		Address: key.Address.Hex(),
		PublicKey: hex.EncodeToString(
			eth_crypto.FromECDSAPub(&key.PrivateKey.PublicKey)),
	}
	if showPrivate {
		out.PrivateKey = hex.EncodeToString(eth_crypto.FromECDSA(key.PrivateKey))
	}

	if outputJSON {
		common.MustPrintJSON(out)
	} else {
		fmt.Println("Address:       ", out.Address)
		fmt.Println("Public key:    ", out.PublicKey)
		if showPrivate {
			fmt.Println("Private key:   ", out.PrivateKey)
		}
	}

	return nil
}
// InspectKeyByMoniker is a wrapper around InspectKey to add moniker support
func InspectKeyByMoniker(keystore string, moniker string, PasswordFile string, showPrivate bool, outputJSON bool) error {
	fp := filepath.Join(keystore, moniker+".json")
	
	if !files.CheckIfExists(fp) {
		return errors.New("cannot find keyfile for that moniker")
	}
	
	return InspectKey(fp, PasswordFile, showPrivate, outputJSON)
}

// GetPrivateKeyString decrypts a keystore and returns the private key as a
// string
func GetPrivateKeyString(keyfilePath string, passwordFile string) (string, error) {
	
	privKey, err := crypto.GetPrivateKey(keyfilePath, passwordFile)
	if err != nil {
		return "", err
	}
	
	return hex.EncodeToString(eth_crypto.FromECDSA(privKey)), nil
}