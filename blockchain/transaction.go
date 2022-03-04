package blockchain

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/pranjalpokharel7/yudhishthira/wallet"
)

type Tx struct {
	TxID       HexByte `json:"txID"`       // hash of this transaction
	UTXOID     HexByte `json:"UTXOID"`     // reference to the hash last transaction the item was a part of
	Signature  HexByte `json:"signature"`  // signature of seller i.e. we need proof that transaction was indeed confirmed by the seller
	ItemHash   HexByte `json:"itemHash"`   // hash of the item involved in transaction
	SellerHash HexByte `json:"sellerHash"` // pubkey hash of the seller
	BuyerHash  HexByte `json:"buyerHash"`  // pubkey hash of the seller
	Amount     uint64  `json:"amount"`     // amount invloved in transaction

	// remove these fields for merkel pls
	InputCount  int
	OutputCount int
}

func (tx Tx) SerializeTxToGOB() ([]byte, error) {
	var encoded bytes.Buffer
	err := gob.NewEncoder(&encoded).Encode(tx)
	return encoded.Bytes(), err // if err in encoding then nil is returned anyway
}

func DeserializeTxFromGOB(serializedTx []byte) (*Tx, error) {
	var tx Tx
	err := gob.NewDecoder(bytes.NewReader(serializedTx)).Decode(&tx)
	return &tx, err
}

func (tx *Tx) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("------------ Transaction: %x ------------", tx.TxID))
	lines = append(lines, fmt.Sprintf("UTXOID: %x", tx.UTXOID))
	lines = append(lines, fmt.Sprintf("Signature: %x", tx.Signature))
	lines = append(lines, fmt.Sprintf("Item Hash: %x", tx.ItemHash))
	lines = append(lines, fmt.Sprintf("Seller Hash: %x", tx.SellerHash))
	lines = append(lines, fmt.Sprintf("Buyer Hash: %x", tx.BuyerHash))
	lines = append(lines, fmt.Sprintf("Amount: %d", tx.Amount))
	return strings.Join(lines, "\n")
}

func (tx *Tx) deepCopy() Tx {
	var txCopy Tx
	txCopy.Amount = tx.Amount
	copy(txCopy.BuyerHash, tx.BuyerHash)
	copy(txCopy.SellerHash, tx.SellerHash)
	copy(txCopy.ItemHash, tx.ItemHash)
	copy(txCopy.UTXOID, tx.UTXOID)
	copy(txCopy.Signature, tx.Signature)
	copy(txCopy.TxID, tx.TxID)
	return txCopy
}

func (tx *Tx) CalculateTxHash() ([]byte, error) {
	var hash [32]byte
	txCopy := tx.deepCopy() // TODO: is this simply shallow copy?
	txCopy.TxID = []byte{}  // TODO: might remove this here, we won't be initializing hash beforehand anyway
	txCopySerialized, err := txCopy.SerializeTxToGOB()
	if err != nil {
		return nil, err
	}
	hash = sha256.Sum256(txCopySerialized)
	return hash[:], nil
}

func CoinBaseTransaction(address string, itemHash []byte, basePrice uint64, chain *BlockChain) (*Tx, error) {
	// check if the address is valid
	pubKeyHash, err := wallet.PubKeyHashFromAddress(address)
	if err != nil {
		return nil, err
	}

	// check if the address has sufficient funds for coinbase transactions
	hasFunds := wallet.CheckSufficientFunds(pubKeyHash) // TODO: complete this method later, for now just returns true
	if !hasFunds {
		return nil, errors.New("the address owner does not have sufficient funds for introducing items into the chain")
	}

	// coinbase transactions have seller hash nil, previous linked output nil
	coinBaseTx := Tx{ItemHash: itemHash, BuyerHash: pubKeyHash, Amount: basePrice, SellerHash: nil}
	coinBaseTx.TxID, err = coinBaseTx.CalculateTxHash()
	if err != nil {
		return nil, err
	}
	return &coinBaseTx, nil
}

func (tx *Tx) IsCoinbase() bool {
	return tx.SellerHash == nil && tx.UTXOID == nil
}

func NewTransaction(addrFrom string, addrTo string, item string, amount uint64, chain *BlockChain) (*Tx, error) {
	itemHash, err := hex.DecodeString(item)
	if err != nil {
		return nil, err
	}

	sellerPubKeyHash, err := wallet.PubKeyHashFromAddress(addrFrom)
	if err != nil {
		return nil, err
	}

	sellerUTXOItems, err := chain.FindItemsOwned(sellerPubKeyHash)
	if err != nil {
		return nil, err
	}

	tx := &Tx{ItemHash: itemHash, SellerHash: sellerPubKeyHash}
	itemOwned := false

	for itemHash, prevItemTransaction := range sellerUTXOItems {
		if itemHash == item {
			// this means the seller does own the item, second variable in map is the transaction corresponding to the item
			// TODO: need to do this
			tx.UTXOID = prevItemTransaction.TxID
			itemOwned = true
		}
	}

	if !itemOwned {
		return nil, errors.New("seller does not own the item")
	}

	buyerPubKeyHash, err := wallet.PubKeyHashFromAddress(addrTo)
	if err != nil {
		return nil, err
	}
	tx.BuyerHash = buyerPubKeyHash
	tx.Amount = amount

	// will not sign transaction here since we do not pass private key into this function
	// also will not hash transaction here, final hash will be calculated after signing transaction only
	return tx, nil
}

func sign(privKey *rsa.PrivateKey, txHash []byte) ([]byte, error) {
	signature, err := rsa.SignPSS(rand.Reader, privKey, crypto.SHA256, txHash, nil)
	return signature, err
}

func (tx *Tx) SignTransaction(wlt *wallet.Wallet, prevUTXO []byte) error {
	sellerPrivKey := wlt.PrivateKey

	if tx.IsCoinbase() {
		signature, err := sign(&sellerPrivKey, tx.TxID)
		if err != nil {
			return err
		}
		tx.Signature = signature
		return nil
	}

	// why do I need to calculate hash again? think, if not, directly use txID
	transactionHash, err := tx.CalculateTxHash()
	if err != nil {
		return err
	}
	// TODO: this check might be redundant, check later
	// assumed the tx has already been linked with prev transaction
	if !bytes.Equal(prevUTXO, tx.UTXOID) {
		return errors.New("previous transaction/item is not the seller's to spend")
	}
	signature, err := rsa.SignPSS(rand.Reader, &sellerPrivKey, crypto.SHA256, transactionHash, nil)
	if err != nil {
		return err
	}
	tx.Signature = signature
	return nil
}

// if we don't get any errors from verify signature then our signature is valid
func (tx *Tx) VerifySignature(sellerPubKey *rsa.PublicKey) error {
	return rsa.VerifyPSS(sellerPubKey, crypto.SHA256, tx.TxID, tx.Signature, nil)
}

// TODO: all these functions below to be implemented
func MinerReward() {
	// TODO: set mining difficulty based on transaction amount? to prevent money laundering lol
}
