package wallet

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
)

// might use elliptic curve cryptography if this ends up becoming too slow
type Wallet struct {
	privateKey rsa.PrivateKey
	PublicKey  rsa.PublicKey
	Address    []byte
}

func PublicKeyToBytes(pubKey *rsa.PublicKey) ([]byte, error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	return pubKeyBytes, nil
}

// instead of using rand.Reader maybe ask the user for passphrase/key of words
func (wallet *Wallet) GenerateKeyPair() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, RSA_KEY_SIZE)
	if err != nil {
		return err
	}
	wallet.privateKey = *privateKey
	wallet.PublicKey = privateKey.PublicKey
	return nil
}

func pubKeyHashRipeMD160(pubKey *rsa.PublicKey) ([]byte, error) {
	// derive public key hash
	pubKeyBytes, err := PublicKeyToBytes(pubKey)
	if err != nil {
		return nil, err
	}
	pubKeyHash := sha256.Sum256(pubKeyBytes)

	// pass through ripemd160 hash
	ripeMDHasher := ripemd160.New()
	_, err = ripeMDHasher.Write(pubKeyHash[:])
	if err != nil {
		return nil, err
	}
	pubKeyRipMD := ripeMDHasher.Sum(nil)
	return pubKeyRipMD, nil
}

func (wallet *Wallet) GenerateAddress() error {
	pubKeyHash, err := pubKeyHashRipeMD160(&wallet.PublicKey)
	if err != nil {
		return err
	}

	// not including version in hash, let's see for now
	checksum := deriveChecksum(pubKeyHash)
	fullHash := append(pubKeyHash, checksum...)
	address := base58.Encode(fullHash)
	wallet.Address = []byte(address)
	return nil
}

func deriveChecksum(pubKeyHash []byte) []byte {
	hashPrimary := sha256.Sum256(pubKeyHash)
	hashFinal := sha256.Sum256(hashPrimary[:])
	return hashFinal[:CHECKSUM_SIZE]
}

func (wallet *Wallet) LoadWalletFromFile(walletFile string) error {
	if _, err := os.Stat(WALLET_FILE); os.IsNotExist(err) {
		return err
	}
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		return err
	}
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(wallet)
	return err // err == nil or error
}

func (wallet *Wallet) SaveWalletToFile(walletFile string) error {
	var wBuffer bytes.Buffer
	err := gob.NewEncoder(&wBuffer).Encode(wallet)
	if err != nil {
		return err
	}
	// 0644 for permissions
	err = ioutil.WriteFile(walletFile, wBuffer.Bytes(), 0644)
	return err // if nil then nil is returned, else error is returned
}

func GenerateWallet(walletFile string) error {
	var wlt Wallet
	err := wlt.GenerateKeyPair()
	if err != nil {
		return err
	}
	err = wlt.GenerateAddress()
	if err != nil {
		return err
	}
	err = wlt.SaveWalletToFile(walletFile)
	if err != nil {
		return err
	}
	fmt.Printf("Wallet generated and saved to %s\n", walletFile)
	fmt.Printf("Your address is %s", wlt.Address)
	return nil
}
