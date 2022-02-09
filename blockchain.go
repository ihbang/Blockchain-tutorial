package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const (
	dbFile       string = "blockchain.db"
	blocksBucket string = "blocks"
)

type Blockchain struct {
	tip []byte // hash value of "tip" Block
	db  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte // hash value of Block at current position
	db          *bolt.DB
}

// AddBlock creates new Block with data and add Block to bc
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	// get lastHash value from db
	_ = bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		lastHash = bucket.Get([]byte("l"))
		return nil
	})

	// create new Block with lastHash and insert it to db
	newBlock := NewBlock(data, lastHash)
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

// NewGenesisBlock creates "Genesis Block" of the blockchain
// "Genesis Block" means the first block of the blockchain
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// NewBlockchain creates new blockchain
func NewBlockchain() *Blockchain {
	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err.Error())
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		// if there is no "blocks" bucket, create new one
		if bucket == nil {
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := NewGenesisBlock()
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
		} else {
			tip = bucket.Get([]byte("l"))
		}
		return nil
	})
	if err != nil {
		log.Panic(err.Error())
	}

	bc := Blockchain{tip, db}
	return &bc
}

// Iterator creates new BlockchainIterator for the Blockchain
func (bc *Blockchain) Iterator() (iter *BlockchainIterator) {
	iter = &BlockchainIterator{bc.tip, bc.db}
	return
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
