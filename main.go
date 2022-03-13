package main

import (
	"fmt"

	"github.com/pranjalpokharel7/yudhishthira/cli"
	"github.com/pranjalpokharel7/yudhishthira/p2p"
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

	p2p.StartServer("3000")
	cli.RunCLI()
}
