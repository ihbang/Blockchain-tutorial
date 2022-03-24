package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

type Block struct {
	Timestamp     int64          // Create timestamp
	Transactions  []*Transaction // Transactions
	PrevBlockHash []byte         // Hash value of previous Block
	Hash          []byte         // Hash value of this Block
	Nonce         int            // Nonce used to generate Hash
}

// NewBlock creates new Block
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce
	return block
}

// Serialize a block to a slice of bytes
func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	if err := encoder.Encode(b); err != nil {
		log.Println(fmt.Errorf("Block serialization failed"))
	}
	return result.Bytes()
}

// DeserializeBlock decodes a slice of bytes to a Block
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))

	if err := decoder.Decode(&block); err != nil {
		log.Println(fmt.Errorf("Block deserialization failed"))
	}

	return &block
}

// HashTransactions calculates sha256 hash value of all txs in the Block
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}
