package blockchain

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// technically block chain is just a chain of blocks
// a single int field to determine whether the chain is main chain or test chain
type BlockChain struct {
	Difficulty     uint8
	miningInterval uint64 // time after which a block is compulsorily added to the chain
	chainType      CHAIN_TYPE
	Blocks         []Block
}

func containsLeadingZeroes(hash [32]byte, difficulty uint8) bool {
	var hexRepresentation string = hex.EncodeToString(hash[:])
	var leadingZeroes string = strings.Repeat("0", int(difficulty))
	// fmt.Println(hexRepresentation[0:difficulty])
	if hexRepresentation[0:difficulty] != leadingZeroes {
		return false
	}
	return true
}

func (bc *BlockChain) ProofOfWork(blk Block) {
	for i := 0; i < MAX_ITERATIONS_POW; i++ { // arbitrary 1000 to prevent potential endless loop
		hash := blk.CalculateHash(uint64(i))
		fmt.Println(hex.EncodeToString(hash[:]))
		if containsLeadingZeroes(hash, bc.Difficulty) {
			fmt.Println("-----------")
			fmt.Println("Final Calculated Hash After PoW")
			fmt.Println(hex.EncodeToString(hash[:]))

			// modify blocks properties
			blk.header.blockHash = hash
			blk.header.nonce = uint64(i) // TODO: change everything to int?

			break
		}
	}
}

// only has the main chain for now
func createGenesisBlock() {}

// load block chain from JSON format
func loadBlockChain() {}

// broadcast your changes to the chain from this function
func updateBlockChain() {}
