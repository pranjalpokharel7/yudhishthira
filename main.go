package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/pranjalpokharel7/yudhishthira/api"
	"github.com/pranjalpokharel7/yudhishthira/blockchain"
	"github.com/pranjalpokharel7/yudhishthira/wallet"
)

// D1Ji4eMTc9nDRfVrZFHg7FRcbueJyKqFG

func main() {
	fmt.Println("This is where it begins...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

	chain := blockchain.InitBlockChain()
	chain.PrintChain()

	api.StartServer(&wlt1, chain)
}
