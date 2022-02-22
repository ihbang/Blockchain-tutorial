package main

import (
	"bytes"

	"github.com/btcsuite/btcutil/base58"
)

type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

// Lock signs the output with the address
func (out *TxOutput) Lock(address string) {
	pubKeyHash := base58.Decode(address)
	out.PubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
}

// IsLockedWith checks if the output can be used by the owner of the pubKey
func (out *TxOutput) IsLockedWith(pubKeyHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, pubKeyHash)
}

func NewTxOutput(value int, address string) *TxOutput {
	txOut := &TxOutput{value, nil}
	txOut.Lock(address)
	return txOut
}
