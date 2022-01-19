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
	blk1.SetBlockValues(uint64(time.Now().Unix()), 4) // will include previous hash parameter in next commit
	blk1.CalculateHash()
}
