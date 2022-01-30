package blockchain

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/pranjalpokharel7/yudhishthira/transaction"
)

type BlockHeader struct {
	version      uint32          // which version of the coin the block is a part of, might be removed
	nonce        uint64          // unsigned representation for now, might allocate 64 bits later, upgrade to 64 bits if version field is removed
	timestamp    uint64          // unix date time, string representation now, might convert to uint64 if time zones are not taken into consideration
	previousHash [HASH_SIZE]byte // hash of previous block
	blockHash    [HASH_SIZE]byte // hash of the current block
}

// the body of the block only contains the transactions
// might need a few more fields
// including the consensus number in the blockheader
type Block struct {
	header       BlockHeader
	transactions []transaction.Tx
}

// will have to jsonify the block
func (blk *Block) PrintBlock() {
	fmt.Printf("Block Hash: %x \n", blk.header.blockHash)
	fmt.Printf("Previous Block Hash: %x \n", blk.header.previousHash)
	fmt.Println("Minted On: ", blk.header.timestamp)
	fmt.Println("Nonce: ", blk.header.nonce)
}

// define how the block is printed
func (blk Block) String() string {
	return "This is a block!"
}

// function to calculate the hash of a block goes here
// might remove the interface format and/or change it to purely functional
// TODO: this function is too expensive, might not use slices after all, F
func (blk *Block) CalculateHash() {
	nonceSlice := make([]byte, blk.header.nonce)         // might need to rename this later
	timestampSlice := make([]byte, blk.header.timestamp) // might need to rename this later
	blk.header.blockHash = sha256.Sum256(append(blk.header.previousHash[:], append(timestampSlice, nonceSlice...)...))
}

func (blk *Block) CreateBlock(nonce uint64) {
	blk.header.timestamp = uint64(time.Now().Unix())
	blk.header.nonce = nonce
}

func (blk *Block) LinkPreviousHash(prevBlock *Block) {
	blk.header.previousHash = prevBlock.header.blockHash
}

func (blk *Block) CreateGenesisBlock(nonce uint64) {
	blk.header.timestamp = uint64(time.Now().Unix())
	blk.header.nonce = nonce
	blk.header.blockHash = sha256.Sum256([]byte("Genesis Block"))
}
