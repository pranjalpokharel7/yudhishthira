package transaction

type Tx struct {
	InputCount  int
	OutputCount int

	ItemHash   []byte
	SellerHash []byte
	BuyerHash  []byte
	Amount     uint64
}
