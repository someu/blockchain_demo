package blockchain

import (
	"blockchain_demo/config"
	"encoding/hex"
	"log"

	"github.com/boltdb/bolt"
)

type UTXOSet struct {
	BlockChain *BlockChain
}

func (u UTXOSet) Reindex() {
	db := u.BlockChain.DB
	bucketName := []byte(config.UTXOBucketName)

	err := db.Update(func(t *bolt.Tx) error {
		err := t.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			return err
		}
		_, err = t.CreateBucket(bucketName)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	utxo := u.BlockChain.FindUTXO()

	err = db.Update(func(t *bolt.Tx) error {
		b := t.Bucket(bucketName)

		for txID, outs := range utxo {
			key, err := hex.DecodeString(txID)
			if err != nil {
				return err
			}
			err = b.Put(key, outs.Serialize())
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0

	db := u.BlockChain.DB

	err := db.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(config.UTXOBucketName))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := DeserializeTXOutputs(v)
			for outIndex, out := range outs.Outputs {
				if out.IsLockedwithKey(pubKeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIndex)
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return accumulated, unspentOutputs
}

func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	db := u.BlockChain.DB

	err := db.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(config.UTXOBucketName))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := DeserializeTXOutputs(v)
			for _, out := range outs.Outputs {
				if out.IsLockedwithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

func (u UTXOSet) Update(block Block) {
	db := u.BlockChain.DB

	err := db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(config.UTXOBucketName))

		for _, tx := range block.Transactions {
			if !tx.IsCoinbase() {
				for _, vin := range tx.Vin {
					updateOutputs := TXOutputs{}
					outs := DeserializeTXOutputs(b.Get(vin.Txid))

					for outIdx, out := range outs.Outputs {
						if outIdx != vin.Vout {
							updateOutputs.Outputs = append(updateOutputs.Outputs, out)
						}
					}
					if len(outs.Outputs) == 0 {
						err := b.Delete(vin.Txid)
						if err != nil {
							return err
						}
					} else {
						err := b.Put(vin.Txid, updateOutputs.Serialize())
						if err != nil {
							return err
						}
					}
				}
			}
			newOutputs := TXOutputs{
				Outputs: tx.Vout,
			}
			err := b.Put(tx.ID, newOutputs.Serialize())
			if err != nil {
				return err
			}

		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
