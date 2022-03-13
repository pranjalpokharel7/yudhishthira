package main

import (
	"fmt"

	"github.com/pranjalpokharel7/yudhishthira/api"
	"github.com/pranjalpokharel7/yudhishthira/blockchain"
	"github.com/pranjalpokharel7/yudhishthira/wallet"
)

// D1Ji4eMTc9nDRfVrZFHg7FRcbueJyKqFG

func main() {
	fmt.Println("This is where it begins...")

	chain := blockchain.InitBlockChain()

	var wlt1 wallet.Wallet
	var wlt0 wallet.Wallet
	wlt1.LoadWalletFromFile("./mykeys.keystore")
	wlt0.LoadWalletFromFile("./k0.keystore")
	fmt.Println(string(wlt1.Address))
	fmt.Println(string(wlt0.Address))

	// pubKeyHash, err := wallet.PubKeyHashFromAddress(string(wlt1.Address))
	// utility.ErrThenPanic(err)
	// fmt.Printf("%x", pubKeyHash)

	// cli.RunCLI()

	// block := blockchain.CreateBlock()
	// block.MineBlock(chain, &wlt)

	// chain.AddBlock(block)
	// chain.PrintChain()

	api.StartServer(&wlt1, chain)
}
