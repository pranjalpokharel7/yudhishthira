package blockchain

import (
	"crypto/sha256"

	"github.com/pranjalpokharel7/pdb/transaction"
)

type BlockHeader struct {
	nonce        uint64            // unsigned representation for now, might allocate 64 bits later, upgrade to 64 bits if version field is removed
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
func (blk *Block) CalculateHash(nonce uint64) [sha256.Size]byte {
	var temp_slice []byte = make([]byte, nonce, blk.header.timestamp) // might need to rename this later
	//blk.header.blockHash = sha256.Sum256(append(blk.header.previousHash[:], temp_slice...))

	blockHash := sha256.Sum256(append(blk.header.previousHash[:], temp_slice...))
	return blockHash
}

func (blk *Block) InitBlock(timestamp uint64) {
	blk.header.timestamp = timestamp
	blk.header.previousHash = sha256.Sum256([]byte("This is the previous hash..."))
	blk.header.nonce = 0
}
