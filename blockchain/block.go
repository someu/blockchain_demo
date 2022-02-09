package blockchain

import (
	"blockchain_demo/transaction"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp     int64
	Transactions  []*transaction.Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	if err := encoder.Encode(b); err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func (b *Block) HashTransactions() []byte {
	var txHashs [][]byte
	var txHash [32]byte
	for _, tx := range b.Transactions {
		txHashs = append(txHashs, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashs, []byte{}))
	return txHash[:]
}

func NewBlock(transactions []*transaction.Transaction, prevBlockHash []byte) *Block {
	b := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
	}
	p := NewProofOfWork(b)
	nonce, hash := p.Run()
	b.Nonce = nonce
	b.Hash = hash
	return b
}

func NewGenesisBlock(coinbase *transaction.Transaction) *Block {
	return NewBlock([]*transaction.Transaction{coinbase}, []byte{})
}

func DeserializeBlock(d []byte) *Block {
	var b Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	if err := decoder.Decode(&b); err != nil {
		log.Panic(err)
	}
	return &b
}
