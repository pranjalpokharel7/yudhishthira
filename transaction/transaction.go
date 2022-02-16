package transaction

type Tx struct {
	InputCount  int `json:"inputCount"`
	OutputCount int `json:"outputCount"`

	ItemHash   []byte `json:"itemHash"`
	SellerHash []byte `json:"sellerHash"`
	BuyerHash  []byte `json:"buyerHash"`
	Amount     uint64 `json:"amount"`
}

// TODO: serialize this part
func (tx Tx) Serialize() ([]byte, error) {
	return nil, nil
}
