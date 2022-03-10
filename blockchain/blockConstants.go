package blockchain

import (
	"encoding/hex"
	"encoding/json"
)

// all blockchain related constants go here
type CHAIN_TYPE int

// different chains for different purposes
const (
	MAIN_CHAIN CHAIN_TYPE = iota
	TEST_CHAIN
)

const (
	MAX_ITERATIONS_POW = 100000 // will migrate to unlimited iterations once optimized to use goroutines for mining
	DB_PATH            = "./db"
	LAST_HASH          = "lh"
	GENESIS_STRING     = "BBC News (Thursday, March 10, 2022 1:33:39 PM) - Ukraine war: No progress on ceasefire after Kyiv-Moscow talks"
	GENESIS_TIMESTAMP  = 1646919219
)

// custom byte type for marshaling
type HexByte []byte

func (hb HexByte) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(hb))
}
