package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct{}

func (cli CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  help: Print usage")
	fmt.Println("  printchain: Print all the blocks of the blockchain")
	fmt.Println("  getbalance -address ADDRESS: Get a balance of ADDRESS")
	fmt.Println("  createblockchain -address ADDRESS: Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT: Send AMOUNT of coins from FROM address to TO")
}

func (cli CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	var err error
	cli.validateArgs()

	var (
		_                       = flag.NewFlagSet("help", flag.ExitOnError)
		printChainCmd           = flag.NewFlagSet("printchain", flag.ExitOnError)
		getBalanceCmd           = flag.NewFlagSet("getbalance", flag.ExitOnError)
		createBlockchainCmd     = flag.NewFlagSet("createblockchain", flag.ExitOnError)
		sendCmd                 = flag.NewFlagSet("send", flag.ExitOnError)
		getBalanceAddress       = getBalanceCmd.String("address", "", "The address to get balance for")
		createBlockchainAddress = createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
		sendFrom                = sendCmd.String("from", "", "Source wallet address")
		sendTo                  = sendCmd.String("to", "", "Destination wallet address")
		sendAmount              = sendCmd.Int("amount", 0, "Amount of coins to send")
	)

	switch os.Args[1] {
	case "help":
		cli.printUsage()
		os.Exit(0)
	case "printchain":
		err = printChainCmd.Parse(os.Args[2:])
	case "getbalance":
		err = getBalanceCmd.Parse(os.Args[2:])
	case "createblockchain":
		err = createBlockchainCmd.Parse(os.Args[2:])
	case "send":
		err = sendCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if err != nil {
		log.Panic(err)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

}

func (cli CLI) printChain() {
	// TODO: Fix this
	bc := NewBlockchain("")
	defer bc.db.Close()

	iter := bc.Iterator()

	for {
		block := iter.Next()
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

// getBalance prints total balance of address
func (cli CLI) getBalance(address string) {
	bc := NewBlockchain(address)
	defer bc.db.Close()

	balance := 0
	UTXOs := bc.FindUnspentTxOuts(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

// createBlockchain creates a new blockchain and give reward to address
func (cli *CLI) createBlockchain(address string) {
	bc := CreateBlockchain(address)
	bc.db.Close()
	fmt.Println("Done!")
}

func (cli *CLI) send(from, to string, amount int) {
	bc := NewBlockchain(from)
	defer bc.db.Close()

	tx := NewTransaction(from, to, amount, bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("Success!")
}
