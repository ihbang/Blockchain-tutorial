package main

import (
	"bytes"
)

type TxInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	Pubkey    []byte
}

// UsesKey checks if a public key in the input is same as pubKeyHash
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.Pubkey)
	return bytes.Equal(lockingHash, pubKeyHash)
}
