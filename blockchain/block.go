package blockchain

import "github.com/pranjalpokharel7/pdb/transaction"

type BlockHeader struct {
	index        uint64 // index of the block in the chain
	version      uint16 // which version of the coin the block is a part of
	previousHash string // hash of previous block
	blockHash    string // hash of the current block
	timestamp    uint64 // unix date time
	nonce        uint32 // unsigned representation for now, might allocate 64 bits later
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
	blk.header.blockHash = ""
}
