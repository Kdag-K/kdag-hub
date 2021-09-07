package crypto

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	
	"github.com/Kdag-K/evm/src/crypto"
	"github.com/Kdag-K/kdag-hub/src/common"
	"github.com/Kdag-K/kdag-hub/src/configuration"
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

// InspectKey inspects an encrypted keyfile.
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
// InspectKeyByMoniker is a wrapper around InspectKey to add moniker support.
func InspectKeyByMoniker(keystore string, moniker string, PasswordFile string, showPrivate bool, outputJSON bool) error {
	fp := filepath.Join(keystore, moniker+".json")
	
	if !files.CheckIfExists(fp) {
		return errors.New("cannot find keyfile for that moniker")
	}
	
	return InspectKey(fp, PasswordFile, showPrivate, outputJSON)
}

// GetPrivateKeyString decrypts a keystore and returns the private key as a string.
func GetPrivateKeyString(keyfilePath string, passwordFile string) (string, error) {
	
	privKey, err := crypto.GetPrivateKey(keyfilePath, passwordFile)
	if err != nil {
		return "", err
	}
	
	return hex.EncodeToString(eth_crypto.FromECDSA(privKey)), nil
}

// GetPrivateKey decrypts a keystore and returns the private key.
func GetPrivateKey(keyfilepath string, PasswordFile string) (*ecdsa.PrivateKey, error) {
	
	// Read key from file.
	keyjson, err := ioutil.ReadFile(keyfilepath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the keyfile at '%s': %v", keyfilepath, err)
	}
	
	// Decrypt key with passphrase.
	passphrase, err := crypto.GetPassphrase(PasswordFile, false)
	if err != nil {
		return nil, err
	}
	
	key, err := keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		return nil, fmt.Errorf("Error decrypting key: %v", err)
	}
	
	return key.PrivateKey, nil
	
}

/*
GenerateKeyfile generates an Ethereum keyfile and writes it.

keyfilepath: path to write the new keyfile to.
passwordFile: plain text file containing the passphrase to use for the
              keyfile.

privateKeyfile: the path to an unencrypted private key. If specified, this
                function does not generate a new keyfile, it instead
                generates a keyfile from the unencrypted private key.

outputJSON: controls whether the output to stdio is in JSON format or not.
            The function returns a key object which can be used to retrieve
            public or private keys or the address.
*/
func GenerateKeyfile(keyfilepath, passwordFile, privateKeyfile string, outputJSON bool) (*keystore.Key, error) {
	const dirPerm = 0700
	const filePerm = 0600
	if keyfilepath == "" {
		keyfilepath = configuration.DefaultKeyfile
	}
	if _, err := os.Stat(keyfilepath); err == nil {
		return nil, fmt.Errorf("Keyfile already exists at %s", keyfilepath)
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("Error checking if keyfile exists: %v", err)
	}
	
	var privateKey *ecdsa.PrivateKey
	var err error
	
	if file := privateKeyfile; file != "" {
		// Load private key from file.
		common.DebugMessage("Loading Private Key: ", file)
		privateKey, err = eth_crypto.LoadECDSA(file)
		if err != nil {
			return nil, fmt.Errorf("Can't load private key: %v", err)
		}
	} else {
		// If not loaded, generate random.
		privateKey, err = eth_crypto.GenerateKey()
		if err != nil {
			return nil, fmt.Errorf("Failed to generate random private key: %v", err)
		}
	}
	
	// Create the keyfile object with a random UUID
	key := WrapKey(privateKey)
	
	// Encrypt key with passphrase.
	passphrase, err := crypto.GetPassphrase(passwordFile, true)
	if err != nil {
		return nil, err
	}
	
	keyjson, err := EncryptKey(key, passphrase, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return nil, fmt.Errorf("Error encrypting key: %v", err)
	}
	
	// Load private key from file.
	common.DebugMessage("Destination: ", keyfilepath)
	
	// Store the file to disk.
	if err := os.MkdirAll(filepath.Dir(keyfilepath), dirPerm); err != nil {
		return nil, fmt.Errorf("Could not create directory %s: %v", filepath.Dir(keyfilepath), err)
	}
	if err := ioutil.WriteFile(keyfilepath, keyjson, filePerm); err != nil {
		return nil, fmt.Errorf("Failed to write keyfile to %s: %v", keyfilepath, err)
	}
	
	// Output some information.
	out := outputGenerate{
		Address: key.Address.Hex(),
	}
	
	if outputJSON {
		common.MustPrintJSON(out)
	} else {
		fmt.Println("Address:", out.Address)
	}
	
	return key, nil
}