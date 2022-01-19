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
	tree, _ = merkel.CreateMerkelTree(transactions, tree)

	for i := 0; i < 15; i++ {
		tx := transaction.Tx{
			InputCount: i,
		}

		fmt.Println(tree.VerifyTransaction(tx))
	}
}
