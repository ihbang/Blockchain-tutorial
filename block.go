package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
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

	block.Hash = hash
	block.Nonce = nonce
	return block
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	if err := encoder.Encode(b); err != nil {
		_ = fmt.Errorf("Block serialization failed\n")
	}
	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))

	if err := decoder.Decode(&block); err != nil {
		_ = fmt.Errorf("Block deserialization failed\n")
	}

	return &block
}
