package server

import (
	"blockchain_demo/blockchain"
	"blockchain_demo/utils"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type block struct {
	AddrFrom string
	Block    []byte
}

func sendBlock(addr string, b *blockchain.Block) {
	data := block{nodeAddress, b.Serialize()}
	payload := utils.GobEncode(data)
	request := append(commandToBytes("block"), payload...)

	sendData(addr, request)
}

func handleBlock(request []byte, bc *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload block

	buff.Write(extractCommand(request))
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := blockchain.DeserializeBlock(blockData)

	fmt.Println("Recevied a new block!")
	bc.AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hash)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := blockchain.UTXOSet{
			BlockChain: bc,
		}
		UTXOSet.Reindex()
	}
}
