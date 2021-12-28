package main

import (
	"fmt"

	"github.com/pranjalpokharel7/pdb/transaction"
)

func main() {
	fmt.Println("This is where it begins...")

	// example usage
	var tr transaction.Tx
	tr.InputCount = 3
	fmt.Println(tr.InputCount)
}
