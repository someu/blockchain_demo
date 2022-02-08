package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  reindexutxo - Rebuilds the UTXO set")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT -mine - Send AMOUNT of coins from FROM address to TO. Mine on the same node, when -mine is set.")
	fmt.Println("  startnode -miner ADDRESS - Start a node with ID specified in NODE_ID env. var. -miner enables mining")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	listAddressCmd := flag.NewFlagSet("listaddress", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)

	createBlockChainAddress := createBlockChainCmd.String("address", "", "address")
	getBalanceAddress := getBalanceCmd.String("address", "", "addresss")
	sendFrom := sendCmd.String("from", "", "from")
	sendTo := sendCmd.String("to", "", "to")
	sendAmount := sendCmd.Int("amount", 0, "amount")

	subCmdArgs := os.Args[2:]
	var err error
	switch os.Args[1] {
	case "createblockchain":
		err = createBlockChainCmd.Parse(subCmdArgs)
	case "createwallet":
		err = createWalletCmd.Parse(subCmdArgs)
	case "getbalance":
		err = getBalanceCmd.Parse(subCmdArgs)
	case "listaddresses":
		err = listAddressCmd.Parse(subCmdArgs)
	case "printchain":
		err = printChainCmd.Parse(subCmdArgs)
	case "send":
		err = sendCmd.Parse(subCmdArgs)
	case "reindexutxo":
		err = reindexUTXOCmd.Parse(subCmdArgs)
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if err != nil {
		log.Panic(err)
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockChainAddress == "" {
			createBlockChainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockChainAddress)
	}
	if createWalletCmd.Parsed() {
		cli.createWallet()
	}
	if listAddressCmd.Parsed() {
		cli.listAddresses()
	}
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}
	if printChainCmd.Parsed() {
		cli.PrintChain()
	}
	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		cli.Send(*sendFrom, *sendTo, *sendAmount)
	}
	if reindexUTXOCmd.Parsed() {
		cli.ReindexUTXO()
	}

}

func New() *CLI {
	return &CLI{}
}
