package blockchain

import (
	"blockchain_demo/config"
	"blockchain_demo/transaction"
	"blockchain_demo/wallet"
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
			outs := transaction.DeserializeTXOutputs(v)
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

func (u UTXOSet) FindUTXO(pubKeyHash []byte) []transaction.TXOutput {
	var UTXOs []transaction.TXOutput
	db := u.BlockChain.DB

	err := db.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(config.UTXOBucketName))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := transaction.DeserializeTXOutputs(v)
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
					updateOutputs := transaction.TXOutputs{}
					outs := transaction.DeserializeTXOutputs(b.Get(vin.Txid))

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
			newOutputs := transaction.TXOutputs{
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

func NewUTXOTransaction(from string, to string, amount int, UTXOSet *UTXOSet) *transaction.Transaction {
	var inputs []transaction.TXInput
	var outputs []transaction.TXOutput

	wallets, err := wallet.NewWallets()
	if err != nil {
		log.Panic(err)
	}
	w := wallets.GetWallet(from)
	pubKeyHash := wallet.HashPubKey(w.PublicKey)

	acc, spendableOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("no enough funds")
	}

	for txid, outs := range spendableOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		for _, out := range outs {
			inputs = append(inputs, transaction.TXInput{
				Txid:      txID,
				Vout:      out,
				Signature: nil,
				PubKey:    w.PublicKey,
			})
		}
	}
	outputs = append(outputs, transaction.NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, transaction.NewTXOutput(acc-amount, from))
	}
	tx := transaction.Transaction{
		ID:   nil,
		Vin:  inputs,
		Vout: outputs,
	}
	tx.ID = tx.Hash()
	UTXOSet.BlockChain.SignTransaction(&tx, w.PrivateKey)

	return &tx
}
