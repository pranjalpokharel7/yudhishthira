package wallet

import (
	"crypto/rand"
	"crypto/rsa"
)

// might use elliptic curve cryptography if this ends up becoming too slow
type Wallet struct {
	privateKey *rsa.PrivateKey
	PublicKey  rsa.PublicKey
}

func (wallet *Wallet) generateKeyPair() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, RSA_KEY_SIZE)
	if err != nil {
		return err
	}

	wallet.privateKey = privateKey
	wallet.PublicKey = privateKey.PublicKey
	return nil
}
