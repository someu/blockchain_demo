package blockchain

import (
	"blockchain_demo/config"
	"blockchain_demo/transaction"
	"blockchain_demo/wallet"
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

type BlockChain struct {
	tip []byte
	DB  *bolt.DB
}

func (bc *BlockChain) MineBlock(transactions []*transaction.Transaction) *Block {
	for _, tx := range transactions {
		if !bc.VerifyTransaction(*tx) {
			log.Panic("invalid transaction")
		}
	}

	var lastHash []byte
	err := bc.DB.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(config.BlockChainBucketName))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	newBlock := NewBlock(transactions, lastHash)
	err = bc.DB.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(config.BlockChainBucketName))
		if err := b.Put(newBlock.Hash, newBlock.Serialize()); err != nil {
			return err
		}
		if err := b.Put([]byte("l"), newBlock.Hash); err != nil {
			return err
		}
		bc.tip = newBlock.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return newBlock
}

func (bc *BlockChain) Iterator() *BlcokChainIterator {
	return &BlcokChainIterator{
		currentHash: bc.tip,
		db:          bc.DB,
	}
}

func (bc *BlockChain) FindUTXO() map[string]transaction.TXOutputs {
	var UTXOs = make(map[string]transaction.TXOutputs)

	spentTXOs := make(map[string][]int)
	var bci *BlcokChainIterator
	bci = bc.Iterator()
	for {
		block := bci.Next()
		if block == nil {
			break
		}
		for _, tx := range block.Transactions {
			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					inTxId := hex.EncodeToString(in.Txid)
					spentTXOs[inTxId] = append(spentTXOs[inTxId], in.Vout)
				}
			}
		}
	}

	bci = bc.Iterator()
	for {
		block := bci.Next()
		if block == nil {
			break
		}
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Spent:
			for outIndex, out := range tx.Vout {
				for _, spentOut := range spentTXOs[txID] {
					if spentOut == outIndex {
						break Spent
					}
				}
				outs := UTXOs[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXOs[txID] = outs
			}
		}
	}

	return UTXOs
}

func (bc *BlockChain) FindTransaction(ID []byte) *transaction.Transaction {
	bci := bc.Iterator()
	for {
		bc := bci.Next()
		if bc == nil {
			break
		}
		for _, tx := range bc.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return tx
			}
		}
	}
	return nil
}

func (bc *BlockChain) GetPrevTransactions(tx transaction.Transaction) map[string]transaction.Transaction {
	prevTXs := make(map[string]transaction.Transaction)

	for _, vin := range tx.Vin {
		prevTX := bc.FindTransaction(vin.Txid)
		if prevTX == nil {
			log.Panic("no transaction")
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = *prevTX
	}
	return prevTXs
}

func (bc *BlockChain) SignTransaction(tx *transaction.Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := bc.GetPrevTransactions(*tx)
	tx.Sign(privKey, prevTXs)
}

func (bc *BlockChain) VerifyTransaction(tx transaction.Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	prevTXs := bc.GetPrevTransactions(tx)
	return tx.Verify(prevTXs)
}

func NewBlockChain() *BlockChain {
	if !dbExist(config.DBFile) {
		log.Panic("db file not exist")
	}
	db, err := bolt.Open(config.DBFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	var tip []byte
	err = db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(config.BlockChainBucketName))
		if b == nil {
			return errors.New("invalid db")
		}
		tip = b.Get([]byte("l"))
		if tip == nil {
			return errors.New("invalid l hash")
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return &BlockChain{
		tip: tip,
		DB:  db,
	}
}

func CreateBlockchain(address string) *BlockChain {
	if dbExist(config.DBFile) {
		log.Panic("db already exist")
	}
	if address == "" || !wallet.ValidateAddress(address) {
		log.Panic("invalid genesis address")
	}
	db, err := bolt.Open(config.DBFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	var tip []byte
	err = db.Update(func(t *bolt.Tx) error {
		genesis := NewGenesisBlock(transaction.NewCoinbaseTx(address, "to genesis"))
		b, err := t.CreateBucket([]byte(config.BlockChainBucketName))
		if err != nil {
			return err
		}
		if err = b.Put(genesis.Hash, genesis.Serialize()); err != nil {
			return err
		}
		if err = b.Put([]byte("l"), genesis.Hash); err != nil {
			return err
		}
		tip = genesis.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return &BlockChain{
		tip: tip,
		DB:  db,
	}
}

func dbExist(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}
