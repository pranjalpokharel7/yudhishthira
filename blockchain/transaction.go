package blockchain

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pranjalpokharel7/yudhishthira/wallet"
)

// TODO: timestamp of item -> when coinbase? necessary?
type Tx struct {
	TxID       HexByte `json:"txID"`       // hash of this transaction
	UTXOID     HexByte `json:"UTXOID"`     // reference to the hash last transaction the item was a part of
	Signature  HexByte `json:"signature"`  // signature of seller i.e. we need proof that transaction was indeed confirmed by the seller
	ItemHash   HexByte `json:"itemHash"`   // hash of the item involved in transaction
	SellerHash HexByte `json:"sellerHash"` // pubkey hash of the seller
	BuyerHash  HexByte `json:"buyerHash"`  // pubkey hash of the buyer
	Amount     uint64  `json:"amount"`     // amount invloved in transaction
	Timestamp  uint64  `json:"timestamp"`
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

func (tx Tx) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("------------ Transaction: %x ------------", tx.TxID))
	lines = append(lines, fmt.Sprintf("UTXOID: %x", tx.UTXOID))
	lines = append(lines, fmt.Sprintf("Signature: %x", tx.Signature))
	lines = append(lines, fmt.Sprintf("Item Hash: %x", tx.ItemHash))
	lines = append(lines, fmt.Sprintf("Seller Hash: %x", tx.SellerHash))
	lines = append(lines, fmt.Sprintf("Buyer Hash: %x", tx.BuyerHash))
	lines = append(lines, fmt.Sprintf("Amount: %d", tx.Amount))
	lines = append(lines, fmt.Sprintf("Timestamp: %d", tx.Timestamp))
	return strings.Join(lines, "\n")
}

func (tx *Tx) deepCopy() Tx {
	var txCopy Tx
	txCopy.Amount = tx.Amount
	txCopy.Timestamp = tx.Timestamp
	copy(txCopy.BuyerHash, tx.BuyerHash)
	copy(txCopy.SellerHash, tx.SellerHash)
	copy(txCopy.ItemHash, tx.ItemHash)
	copy(txCopy.UTXOID, tx.UTXOID)
	copy(txCopy.Signature, tx.Signature)
	copy(txCopy.TxID, tx.TxID)
	return txCopy
}

// can be used to verify hash as well, so we need to make a copy beforehand
func (tx *Tx) CalculateTxHash() ([]byte, error) {
	var hash [32]byte
	txCopy := tx.deepCopy()
	txCopy.TxID = []byte{}
	txCopy.Signature = []byte{} // since tx is only signed after hash is calculated, do not factor this into hash calculation
	txCopySerialized, err := txCopy.SerializeTxToGOB()
	if err != nil {
		return nil, err
	}
	hash = sha256.Sum256(txCopySerialized)
	return hash[:], nil
}

func CoinBaseTransaction(srcWallet *wallet.Wallet, itemHash []byte, basePrice uint64, chain *BlockChain) (*Tx, error) {
	// check if the item already exists in the chain before, if yes, can't enter existing item as new item
	itemExists, err := chain.FindItemExists(itemHash)
	if err != nil {
		return nil, err
	}
	if itemExists {
		return nil, errors.New("item already exists in the chain")
	}

	// check if the address is valid
	walletAddress := string(srcWallet.Address)
	pubKeyHash, err := wallet.PubKeyHashFromAddress(walletAddress)
	if err != nil {
		return nil, err
	}

	// check if the address has sufficient funds for coinbase transactions
	hasFunds, err := HasFundsForCoinbaseTx(walletAddress, chain)
	if !hasFunds {
		return nil, errors.New("the address owner does not have sufficient funds for introducing items into the chain")
	}
	if err != nil {
		return nil, err
	}

	// coinbase transactions have seller hash nil, previous linked output nil
	coinBaseTx := Tx{
		ItemHash:   itemHash,
		BuyerHash:  pubKeyHash,
		Amount:     basePrice,
		SellerHash: nil,
		UTXOID:     nil,
		Timestamp:  uint64(time.Now().Unix()),
	}

	// calculate transaciton hash
	txID, err := coinBaseTx.CalculateTxHash()
	if err != nil {
		return nil, err
	}
	coinBaseTx.TxID = txID

	// sign transaction
	coinBaseTx.SignTransaction(srcWallet)

	return &coinBaseTx, nil
}

func (tx *Tx) IsCoinbase() bool {
	return tx.SellerHash == nil && tx.UTXOID == nil
}

func NewTransaction(srcWallet *wallet.Wallet, destinationAddr string, itemHash []byte, basePrice uint64, chain *BlockChain) (*Tx, error) {
	// fetch last transaction the item was a part of
	lastBlockWithItem, txIndex, err := chain.GetLastBlockWithItem(itemHash)
	if err != nil {
		return nil, err
	}
	lastTxWithItem := lastBlockWithItem.TxMerkleTree.LeafNodes[txIndex].Transaction

	// check if the last transaction destination address is the current source address
	sellerPubKeyHash, err := wallet.PubKeyHashFromAddress(string(srcWallet.Address))
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(lastTxWithItem.BuyerHash, sellerPubKeyHash) {
		return nil, errors.New("the item does not belong to the source address to sell")
	}

	// create new transaction
	buyerPubKeyHash, err := wallet.PubKeyHashFromAddress(destinationAddr)
	if err != nil {
		return nil, err
	}
	newTx := Tx{
		ItemHash:   itemHash,
		SellerHash: sellerPubKeyHash,
		BuyerHash:  buyerPubKeyHash,
		Amount:     basePrice,
		UTXOID:     lastTxWithItem.TxID,
		Timestamp:  uint64(time.Now().Unix()),
	}

	// calculate transaction hash
	txID, err := newTx.CalculateTxHash()
	if err != nil {
		return nil, err
	}
	newTx.TxID = txID
	newTx.SignTransaction(srcWallet)

	return &newTx, nil
}

func sign(privKey *rsa.PrivateKey, txHash []byte) ([]byte, error) {
	signature, err := rsa.SignPSS(rand.Reader, privKey, crypto.SHA256, txHash, nil)
	return signature, err
}

func (tx *Tx) SignTransaction(wlt *wallet.Wallet) error {
	sellerPrivKey := wlt.PrivateKey

	if tx.IsCoinbase() {
		signature, err := sign(&sellerPrivKey, tx.TxID)
		if err != nil {
			return err
		}
		tx.Signature = signature
		return nil
	}

	signature, err := rsa.SignPSS(rand.Reader, &sellerPrivKey, crypto.SHA256, tx.TxID, nil)
	if err != nil {
		return err
	}
	tx.Signature = signature
	return nil
}

// if we don't get any errors from verify signature then our signature is valid
func VerifySignature(tx *Tx, sellerPubKey *rsa.PublicKey) error {
	return rsa.VerifyPSS(sellerPubKey, crypto.SHA256, tx.TxID, tx.Signature, nil)
}
