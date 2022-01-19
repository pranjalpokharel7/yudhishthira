package main

import (
	"fmt"

	"github.com/pranjalpokharel7/yudhishthira/merkel"
	"github.com/pranjalpokharel7/yudhishthira/transaction"
)

func main() {
	fmt.Println("This is where it begins...")

	transactions := make([]transaction.Tx, 8)

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
		tree.GetRoot().Print()
	}
}
