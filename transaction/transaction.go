package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

type TxOutput struct {
	ObjectHash      []byte // our object hash is the UTXO in this case
	PubKeyHashBuyer []byte // only the person with the private key of buyer hash can use this output, pubkeyhash is the buyer's digital signature
}
type TxInput struct {
	TxID      []byte // transaction the unspent output is a part of
	Signature []byte // hmm, need to tink
	PubKey    []byte // public key of seller I guess? can be used to calculate PubKeyHash for verification
}
type Tx struct {
	TxID   []byte
	Input  TxInput
	Output TxOutput
	Amount uint64

	// remove input count and output count after making adjustments in the merkel tree
	InputCount  int
	OutputCount int
}

func (tx *Tx) SerializeToGOB() ([]byte, error) {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		return nil, err
	}
	return encoded.Bytes(), nil
}

func (tx *Tx) CalculateTxHash() ([]byte, error) {
	var hash [32]byte
	txCopy := *tx
	txCopy.TxID = []byte{}
	txCopySerialized, err := txCopy.SerializeToGOB()
	if err != nil {
		return nil, err
	}
	hash = sha256.Sum256(txCopySerialized)
	return hash[:], nil
}

func CoinBaseTransaction(pubKey []byte) {
	// TODO: check if the wallet with public key has enough funds to introduce new items into the chain
}
