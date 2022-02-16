package main

import (
	"fmt"

	"github.com/pranjalpokharel7/yudhishthira/blockchain"
	"github.com/pranjalpokharel7/yudhishthira/merkel"
	"github.com/pranjalpokharel7/yudhishthira/transaction"
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

	// var w wallet.Wallet
	// w.GenerateKeyPair()
	// err := w.GenerateAddress()
	// if err != nil {
	// 	panic(err)
	// }
	// pubKeyBytes, _ := wallet.PublicKeyToBytes(&w.PublicKey)
	// fmt.Printf("Wallet Public Key: %x\n", pubKeyBytes)
	// fmt.Printf("Wallet Address: %s\n", w.Address)

	transactions := make([]transaction.Tx, 5)

	for i := range transactions {
		transactions[i] = transaction.Tx{
			InputCount: i,
		}
	}

	var tree *merkel.MerkelTree
	var err error
	tree, err = merkel.CreateMerkelTree(transactions, tree)
	if err != nil {
		println("Error Occured")
	} else {
		// tree.GetRoot().Print()
	}

	data, err := tree.MarshalToJSON()

	if err != nil {

	} else {
		merkel.UnMarshalFromJSON(data)
	}
}
