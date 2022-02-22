package main

import "github.com/boltdb/bolt"

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
