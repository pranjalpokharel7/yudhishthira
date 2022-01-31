package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/pranjalpokharel7/yudhishthira/transaction"
)

type BlockHeader struct {
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

// since this function does not modify the actual block properties, we remove the interface from it
func CalculateHash(blk *Block, nonce uint64) [HASH_SIZE]byte {
	var buf bytes.Buffer

	buf.Write(blk.header.blockHash[:])                    // write blockhash to buffer
	blockData := nonce ^ blk.header.timestamp             // XOR timestamp and nonce
	binary.LittleEndian.PutUint64(buf.Bytes(), blockData) // write XORed  uint64 data to buffer
	calculatedHash := sha256.Sum256(buf.Bytes())          // calculate hash

	return calculatedHash
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
	blk.header.blockHash = sha256.Sum256([]byte(GENESIS_STRING))
}
