package blockchain

import "bytes"

type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

func (in TXInput) UsesKey(pubKeyHash []byte) bool {
	return bytes.Equal(pubKeyHash, HashPubKey(in.PubKey))
}
