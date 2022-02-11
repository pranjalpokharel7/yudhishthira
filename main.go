package main

import (
	"fmt"

	"github.com/pranjalpokharel7/yudhishthira/blockchain"
)

func main() {
	fmt.Println("This is where it begins...")

	var genBlock, b1 blockchain.Block
	var bChain blockchain.BlockChain
	bChain.Difficulty = 3

	genBlock.CreateGenesisBlock(1)
	bChain.AddGenesisBlock(&genBlock)

	b1.CreateBlock(0)
	bChain.AddToBlockchain(&b1)
	bChain.ProofOfWork(&b1)
	bChain.PrintChain()

}
