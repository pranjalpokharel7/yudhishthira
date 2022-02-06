package blockchain

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// technically block chain is just a chain of blocks
// a single int field to determine whether the chain is main chain or test chain
type BlockChain struct {
	Difficulty uint8
	// miningInterval uint64 // time after which a block is compulsorily added to the chain
	// chainType CHAIN_TYPE
	Blocks []Block
}

func containsLeadingZeroes(hash []byte, difficulty uint8) bool {
	var hexRepresentation string = hex.EncodeToString(hash[:])
	var leadingZeroes string = strings.Repeat("0", int(difficulty))
	return hexRepresentation[0:difficulty] == leadingZeroes
}

func (bc *BlockChain) ProofOfWork(blk *Block) {
	for i := uint64(0); i < MAX_ITERATIONS_POW; i++ { // arbitrary 1000 to prevent potential endless loop
		hash := CalculateHash(blk, i)
		if containsLeadingZeroes(hash, bc.Difficulty) {
			blk.BlockHash = hash
			blk.Nonce = i

			break
		}
	}
}

// this function should only be run after proof of work
func (blockchain *BlockChain) AddToBlockchain(block *Block) error {
	if len(blockchain.Blocks) == 0 {
		return errors.New("genesis block not created")
	}
	previousBlock := blockchain.Blocks[len(blockchain.Blocks)-1]
	block.LinkPreviousHash(&previousBlock)
	blockchain.Blocks = append(blockchain.Blocks, *block)
	return nil
}

func (blockchain *BlockChain) PrintChain() {
	for _, block := range blockchain.Blocks {
		blockJson, _ := block.MarshalBlockToJSON()
		fmt.Println(string(blockJson))
	}
}

func (blockchain *BlockChain) AddGenesisBlock(genesisBlock *Block) error {
	if len(blockchain.Blocks) != 0 {
		return errors.New("genesis block already added to the chain")
	}
	blockchain.Blocks = append(blockchain.Blocks, *genesisBlock)
	return nil
}
