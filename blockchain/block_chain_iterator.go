package blockchain

import (
	"blockchain_demo/config"

	"github.com/boltdb/bolt"
)

type BlcokChainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (i *BlcokChainIterator) Next() *Block {
	var block *Block
	i.db.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(config.BlockChainBucketName))
		if b == nil {
			return nil
		}
		d := b.Get(i.currentHash)
		if d == nil {
			return nil
		}
		block = DeserializeBlock(d)
		return nil
	})
	if block == nil {
		i.currentHash = []byte{}
	} else {
		i.currentHash = block.PrevBlockHash
	}

	return block
}
