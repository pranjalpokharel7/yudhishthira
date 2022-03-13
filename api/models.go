package api

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
