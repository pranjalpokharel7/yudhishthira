package api

import (
	"encoding/hex"

	"github.com/pranjalpokharel7/yudhishthira/blockchain"
)

type ErrorJSON struct {
	ErrorMsg string `json:"error"`
}

// TODO: add binding validation
type NewTxFormInput struct {
	Destination string `json:"destination" binding:"required"`
	ItemHash    string `json:"item_hash" binding:"required"`
	Amount      uint64 `json:"amount" binding:"required"`
}

type CoinBaseTxFormInput struct {
	ItemHash string `json:"item_hash" binding:"required"`
	Amount   uint64 `json:"amount" binding:"required"`
}

type TokenSignModel struct {
	Token string `json:"token"`
}

type TokenVerifyModel struct {
	OriginalToken string `json:"token"`
	SignedToken   string `json:"signed_token"`
	PublicKey     string `json:"public_key"`
}

type TransactionsModel struct {
	TxID       string `json:"txID"`
	UTXOID     string `json:"UTXOID"`
	Signature  string `json:"signature"`
	ItemHash   string `json:"itemHash"`
	SellerHash string `json:"sellerHash"`
	BuyerHash  string `json:"buyerHash"`
	Amount     uint64 `json:"amount"`
	Timestamp  uint64 `json:"timestamp"`
}

func ModelToTx(txModel TransactionsModel) (*blockchain.Tx, error) {
	var tx blockchain.Tx
	var err error

	// parse txid
	tx.TxID, err = hex.DecodeString(txModel.TxID)
	if err != nil {
		return nil, err
	}

	tx.UTXOID, err = hex.DecodeString(txModel.UTXOID)
	if err != nil {
		return nil, err
	}

	tx.Signature, err = hex.DecodeString(txModel.Signature)
	if err != nil {
		return nil, err
	}

	tx.BuyerHash, err = hex.DecodeString(txModel.BuyerHash)
	if err != nil {
		return nil, err
	}

	tx.SellerHash, err = hex.DecodeString(txModel.SellerHash)
	if err != nil {
		return nil, err
	}

	tx.ItemHash, err = hex.DecodeString(txModel.ItemHash)
	if err != nil {
		return nil, err
	}

	tx.Amount = txModel.Amount
	tx.Timestamp = txModel.Timestamp

	return &tx, nil
}
