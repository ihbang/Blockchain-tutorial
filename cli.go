package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {
	bc *Blockchain
}

func (cli CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  help: Print usage")
	fmt.Println("  addblock -data DATA: Add new block to blockchain with DATA")
	fmt.Println("  printchain: Print all the blocks of the blockchain")
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
		_             = flag.NewFlagSet("help", flag.ExitOnError)
		addBlockCmd   = flag.NewFlagSet("addblock", flag.ExitOnError)
		printChainCmd = flag.NewFlagSet("printchain", flag.ExitOnError)
		addBlockData  = addBlockCmd.String("data", "", "Block data")
	)

	switch os.Args[1] {
	case "help":
		cli.printUsage()
		os.Exit(0)
	case "addblock":
		err = addBlockCmd.Parse(os.Args[2:])
	case "printchain":
		err = printChainCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if err != nil {
		log.Panic(err)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("Success!")
}

func (cli CLI) printChain() {
	iter := cli.bc.Iterator()

	for {
		block := iter.Next()
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
