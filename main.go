package main

import (
	"fmt"
	"time"

	"github.com/pranjalpokharel7/yudhishthira/blockchain"
)

func main() {
	fmt.Println("This is where it begins...")

	// example usage
	var blk1 blockchain.Block
	blk1.InitBlock(uint64(time.Now().Unix()))

	var chain blockchain.BlockChain
	chain.Difficulty = 2 // might reduce to 1, takes too much time even now
	chain.ProofOfWork(blk1)
}
