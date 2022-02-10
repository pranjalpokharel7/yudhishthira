package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/pranjalpokharel7/yudhishthira/blockchain"
	"github.com/pranjalpokharel7/yudhishthira/p2p"
)

func main() {
	fmt.Println("This is where it begins...")

	// example usage
	var genBlock, b1 blockchain.Block
	var bChain blockchain.BlockChain
	bChain.Difficulty = 3

	genBlock.CreateGenesisBlock(1)
	bChain.AddGenesisBlock(&genBlock)

	b1.CreateBlock(0)
	bChain.ProofOfWork(&b1)
	bChain.AddToBlockchain(&b1)
	// bChain.PrintChain()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parse options from the command line
	listenF := flag.Int("l", 0, "wait for incoming connections")
	targetF := flag.String("d", "", "target peer to dial")
	insecureF := flag.Bool("insecure", false, "use an unencrypted connection")
	seedF := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()
	h, _ := p2p.MakeBasicHost(*listenF, *insecureF, *seedF)

	if *targetF == "" {
		p2p.StartListener(ctx, h, *listenF, *insecureF)
		// Run until canceled.
		<-ctx.Done()
	} else {
		p2p.RunSender(ctx, h, *targetF)
	}
}
