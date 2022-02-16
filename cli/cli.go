package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/pranjalpokharel7/yudhishthira/object"
	"github.com/pranjalpokharel7/yudhishthira/utility"
	"github.com/pranjalpokharel7/yudhishthira/wallet"
)

func CommandLineHelp() {
	fmt.Println("Available Commands:")
	fmt.Println("\t genwallet --file filename - Generate wallet and store in filename")
	fmt.Println("\t objhash --obj filename - Generate object hash for object details stored in filename (csv for now, see ./object/dummy_object.csv)")
}

func RunCLI() {
	if len(os.Args) < 2 {
		CommandLineHelp()
		os.Exit(1)
	}

	genWallet := flag.NewFlagSet("genwallet", flag.ExitOnError)
	checkObjectHash := flag.NewFlagSet("objhash", flag.ExitOnError)

	walletFileLocation := genWallet.String("file", wallet.WALLET_FILE, "The location to store the wallet file")
	objectFileLocation := checkObjectHash.String("obj", "", "The location of the file where the object data is stored")

	switch os.Args[1] {
	case "genwallet":
		err := genWallet.Parse(os.Args[2:])
		utility.ErrThenPanic(err)
	case "objhash":
		err := checkObjectHash.Parse(os.Args[2:])
		utility.ErrThenPanic(err)
	}

	if genWallet.Parsed() {
		wallet.GenerateWallet(*walletFileLocation)
	}

	if checkObjectHash.Parsed() {
		var obj object.Object
		err := obj.LoadCSVData(*objectFileLocation)
		utility.ErrThenPanic(err)
		fmt.Printf("Object Hash: %x\n", obj.HashObject())
	}
}
