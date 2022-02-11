package transaction

type Tx struct {
	InputCount  int
	OutputCount int

	ItemHash   []byte
	SellerHash []byte
	BuyerHash  []byte
	Amount     uint64
}

// TODO: serialize this part
func (tx Tx) Serialize() []byte {
	var data []byte
	data = tx.ItemHash

	return data
}
