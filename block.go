package main

import (
	"time"
)

type Block struct {
	Timestamp     int64  // Create timestamp
	Data          []byte // Data
	PrevBlockHash []byte // Hash value of previous Block
	Hash          []byte // Hash value of this Block
	Nonce         int    // Nonce used to generate Hash
}

// NewBlock creates new Block
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}
