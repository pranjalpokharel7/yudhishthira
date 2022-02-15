package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/pranjalpokharel7/yudhishthira/wallet"
)

func CommandLineHelp() {
	fmt.Println("Available Commands:")
	fmt.Println("\tgenwallet --file filename - Generate wallet and store in filename")
}

func RunCLI() {
	if len(os.Args) < 2 {
		CommandLineHelp()
		os.Exit(1)
	}

	genWallet := flag.NewFlagSet("genwallet", flag.ExitOnError)
	walletFileAddress := genWallet.String("file", wallet.WALLET_FILE, "The location to store the wallet file")

	switch os.Args[1] {
	case "genwallet":
		err := genWallet.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
	}

	if genWallet.Parsed() {
		wallet.GenerateWallet(*walletFileAddress)
	}
}
