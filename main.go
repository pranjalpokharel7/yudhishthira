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
	var w wallet.Wallet
	w.GenerateKeyPair()
	err := w.GenerateAddress()
	if err != nil {
		panic(err)
	}
	pubKeyBytes, _ := wallet.PublicKeyToBytes(&w.PublicKey)
	fmt.Printf("Wallet Public Key: %x\n", pubKeyBytes)
	fmt.Printf("Wallet Address: %s\n", w.Address)
	w.SaveWalletToFile("./k0.keystore")
	// p2p.StartServer("3000")
	// cli.RunCLI()

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
	// block.MineBlock(chain, &wlt1)

	// chain.AddBlock(block)
	chain.PrintChain()

	api.StartServer(&wlt0, chain)
}
