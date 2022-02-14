package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const (
	dbFile              string = "blockchain.db"
	blocksBucket        string = "blocks"
	genesisCoinbaseData string = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

type Blockchain struct {
	tip []byte // hash value of "tip" Block
	db  *bolt.DB
}

// MineBlock mines a new Block with transactions
func (bc *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte

	// get lastHash value from db
	_ = bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		lastHash = bucket.Get([]byte("l"))
		return nil
	})

	// create new Block with lastHash and insert it to db
	newBlock := NewBlock(transactions, lastHash)
	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))

		err := bucket.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}
		err = bucket.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			return err
		}

		bc.tip = newBlock.Hash
		return nil
	})
	if err != nil {
		log.Panic(err.Error())
	}
}

// FindUnspentTransactions returns a slice of all unspent txs for address
func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction
	spentTxOuts := make(map[string][]int)
	iter := bc.Iterator()

	for {
		block := iter.Next()
		for _, tx := range block.Transactions {
			txid := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				if spentTxOuts[txid] != nil {
					for _, spentOut := range spentTxOuts[txid] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlockedWith(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTxOuts[inTxID] = append(spentTxOuts[inTxID], in.Vout)
					}
				}
			}
		}
		// exit condition: reach the genesis block
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTxs
}

// FindUnspentTxOuts returns a slice of all unspent TxOutputs for address
func (bc *Blockchain) FindUnspentTxOuts(address string) []TxOutput {
	var unspentTxOuts []TxOutput
	unspentTxs := bc.FindUnspentTransactions(address)
	for _, tx := range unspentTxs {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				unspentTxOuts = append(unspentTxOuts, out)
			}
		}
	}
	return unspentTxOuts
}

// FindSpendableOutputs find candidate TxOuts to make 'amount' of coins for another transaction
// return - accumulated coins, txid: txOutIdxs
func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentTxOuts := make(map[string][]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txid := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentTxOuts[txid] = append(unspentTxOuts[txid], outIdx)
			}
			if accumulated >= amount {
				break Work
			}
		}
	}
	return accumulated, unspentTxOuts
}

// Iterator creates new BlockchainIterator for the Blockchain
func (bc *Blockchain) Iterator() (iter *BlockchainIterator) {
	iter = &BlockchainIterator{bc.tip, bc.db}
	return
}

type BlockchainIterator struct {
	currentHash []byte // hash value of Block at current position
	db          *bolt.DB
}

// Next returns next Block of the Blockchain
func (iter *BlockchainIterator) Next() *Block {
	var block *Block

	// get serialized Block with currentHash and deserialize it
	_ = iter.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		serializedBlock := bucket.Get(iter.currentHash)
		block = DeserializeBlock(serializedBlock)

		return nil
	})
	// update currentHash to current Block's PrevBlockHash
	iter.currentHash = block.PrevBlockHash
	return block
}

// NewGenesisBlock creates "Genesis Block" of the blockchain
// "Genesis Block" means the first block of the blockchain
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

func dbExists() bool {
	_, err := os.Stat(dbFile)
	return !os.IsNotExist(err)
}

// NewBlockchain creates a new blockchain from "tip" Block
func NewBlockchain(address string) *Blockchain {
	if !dbExists() {
		fmt.Println("No existing Blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err.Error())
	}

	_ = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		tip = bucket.Get([]byte("l"))

		return nil
	})

	bc := Blockchain{tip, db}
	return &bc
}

// CreateBlockchain creates a new Blockchain with genesis Block
func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err.Error())
	}

	err = db.Update(func(tx *bolt.Tx) error {
		coinbase := NewCoinbaseTx(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(coinbase)
		bucket, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			return err
		}

		// serialized Block is mapped with Hash value
		err = bucket.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			return err
		}

		// key "l" stores a Hash value of the last Block in the Blockchain
		err = bucket.Put([]byte("l"), genesis.Hash)
		if err != nil {
			return err
		}
		tip = genesis.Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}
	return &bc
}
