package main

import (
	"fmt"

	"github.com/pranjalpokharel7/yudhishthira/merkel"
	"github.com/pranjalpokharel7/yudhishthira/transaction"
)

func main() {
	fmt.Println("This is where it begins...")
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
