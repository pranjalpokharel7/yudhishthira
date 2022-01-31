package blockchain

import (
	"errors"
)

// technically block chain is just a chain of blocks
// a single int field to determine whether the chain is main chain or test chain
type BlockChain struct {
	chainType CHAIN_TYPE
	Blocks    []Block // store blocks or references to blocks?
}

func (blockchain *BlockChain) AddToBlockchain(block *Block) error {
	// maybe not create block here?
	if len(blockchain.Blocks) == 0 {
		return errors.New("Genesis block not created!")
	}
	previousBlock := blockchain.Blocks[len(blockchain.Blocks)-1]
	block.LinkPreviousHash(&previousBlock)
	block.CalculateHash()
	blockchain.Blocks = append(blockchain.Blocks, *block)
	return nil
}

func (blockchain *BlockChain) PrintChain() {
	for _, value := range blockchain.Blocks {
		value.PrintBlock()
	}
}

func (blockchain *BlockChain) AddGenesisBlock(genesisBlock *Block) error {
	if len(blockchain.Blocks) != 0 {
		return errors.New("Genesis block already added to the chain!")
	}
	blockchain.Blocks = append(blockchain.Blocks, *genesisBlock)
	return nil
}

// load block chain from JSON format
func loadBlockChain() {}

// broadcast your changes to the chain from this function
func updateBlockChain() {}
