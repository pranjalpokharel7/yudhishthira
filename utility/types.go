package utility

import (
	"encoding/hex"
	"encoding/json"
)

// custom byte type for marshaling
type HexByte []byte

func (hb HexByte) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(hb))
}
