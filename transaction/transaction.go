package transaction

import (
	"encoding/json"
	"fmt"
)

type Tx struct {
	InputCount  int `json:"inputCount"`
	OutputCount int `json:"outputCount"`

	ItemHash   []byte `json:"itemHash"`
	SellerHash []byte `json:"sellerHash"`
	BuyerHash  []byte `json:"buyerHash"`
	Amount     uint64 `json:"amount"`
}

func (tx *Tx) CalculateHash() []byte {
	var hash []byte
	return hash
}

func (tx Tx) Serialize() ([]byte, error) {
	jsonData, err := json.Marshal(tx)

	if err != nil {
		fmt.Println("Transaction error:", err)
	}

	return jsonData, nil
}

func Deserialize(data []byte) (*Tx, error) {
	var tx *Tx
	var unmarshalInterface map[string]interface{}

	if err := json.Unmarshal(data, &unmarshalInterface); err != nil {
		fmt.Println(err)
		return nil, err
	}

	for k, v := range unmarshalInterface {
		if k == "inputCount" {
			tx.InputCount = int(v.(float64))
		} else if k == "outputCount" {
			tx.OutputCount = int(v.(float64))
		} else if k == "itemHash" {
			if v != nil {
				tx.ItemHash = v.([]byte)
			}
		} else if k == "sellerHash" {
			if v != nil {
				tx.SellerHash = v.([]byte)
			}
		} else if k == "buyerHash" {
			if v != nil {
				tx.BuyerHash = v.([]byte)
			}
		} else if k == "amount" {
			tx.Amount = uint64(v.(float64))
		}
	}

	return tx, nil
}
