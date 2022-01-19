package blockchain

import (
	"crypto/sha256"
	"fmt"

	"github.com/pranjalpokharel7/yudhishthira/transaction"
)

type BlockHeader struct {
	version      uint32            // which version of the coin the block is a part of, might be removed
	nonce        uint32            // unsigned representation for now, might allocate 64 bits later, upgrade to 64 bits if version field is removed
	timestamp    uint64            // unix date time, string representation now, might convert to uint64 if time zones are not taken into consideration
	previousHash [sha256.Size]byte // hash of previous block
	blockHash    [sha256.Size]byte // hash of the current block
}

// the body of the block only contains the transactions
// might need a few more fields
// including the consensus number in the blockheader
type Block struct {
	header       BlockHeader
	transactions []transaction.Tx
}

// function to calculate the hash of a block goes here
// might remove the interface format and/or change it to purely functional
func (blk *Block) CalculateHash() {
	var temp_slice []byte = make([]byte, blk.header.nonce, blk.header.nonce) // might need to rename this later

	// need to extract previousHash to SetBlockValues parameter, design pattern to be decided
	blk.header.previousHash = sha256.Sum256([]byte("This is the previous hash..."))

	blk.header.blockHash = sha256.Sum256(append(blk.header.previousHash[:], temp_slice...))

	// TODO: remove this printf later
	fmt.Printf("%x\n", blk.header.blockHash)
}

func (blk *Block) SetBlockValues(timestamp uint64, nonce uint32) {
	blk.header.timestamp = timestamp
	blk.header.nonce = nonce
}
