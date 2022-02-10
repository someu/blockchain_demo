package server

import (
	"blockchain_demo/blockchain"
	"blockchain_demo/transaction"
	"blockchain_demo/utils"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type tx struct {
	AddrFrom    string
	Transaction []byte
}

func sendTx(addr string, tnx *transaction.Transaction) {
	data := tx{nodeAddress, tnx.Serialize()}
	payload := utils.GobEncode(data)
	request := append(commandToBytes("tx"), payload...)

	sendData(addr, request)
}

func handleTx(request []byte, bc *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload tx

	buff.Write(extractCommand(request))
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Transaction
	tx := transaction.DeserializeTransaction(txData)
	mempool[hex.EncodeToString(tx.ID)] = tx

	if nodeAddress == knownNodes[0] {
		for _, node := range knownNodes {
			if node != nodeAddress && node != payload.AddrFrom {
				sendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		if len(mempool) >= 2 && len(miningAddress) > 0 {
		MineTransactions:
			var txs []*transaction.Transaction

			for id := range mempool {
				tx := mempool[id]
				if bc.VerifyTransaction(tx) {
					txs = append(txs, &tx)
				}
			}

			if len(txs) == 0 {
				fmt.Println("All transactions are invalid! Waiting for new ones...")
				return
			}

			cbTx := transaction.NewCoinbaseTx(miningAddress, "")
			txs = append(txs, cbTx)

			newBlock := bc.MineBlock(txs)
			UTXOSet := blockchain.UTXOSet{
				BlockChain: bc,
			}
			UTXOSet.Reindex()

			fmt.Println("New block is mined!")

			for _, tx := range txs {
				txID := hex.EncodeToString(tx.ID)
				delete(mempool, txID)
			}

			for _, node := range knownNodes {
				if node != nodeAddress {
					sendInv(node, "block", [][]byte{newBlock.Hash})
				}
			}

			if len(mempool) > 0 {
				goto MineTransactions
			}
		}
	}
}
