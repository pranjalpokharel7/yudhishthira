package wallet

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"

	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
)

// might use elliptic curve cryptography if this ends up becoming too slow
type Wallet struct {
	privateKey rsa.PrivateKey
	PublicKey  rsa.PublicKey
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

func (wallet *Wallet) GenerateAddress() (string, error) {
	pubKeyHash, err := pubKeyHashRipeMD160(&wallet.PublicKey)
	if err != nil {
		return "", err
	}

	// not including version in hash, let's see for now
	checksum := deriveChecksum(pubKeyHash)
	fullHash := append(pubKeyHash, checksum...)
	address := base58.Encode(fullHash)
	return address, nil
}

func deriveChecksum(pubKeyHash []byte) []byte {
	hashPrimary := sha256.Sum256(pubKeyHash)
	hashFinal := sha256.Sum256(hashPrimary[:])

	return hashFinal[:CHECKSUM_SIZE]
}
