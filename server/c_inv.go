package server

import (
	"blockchain_demo/blockchain"
	"blockchain_demo/utils"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

func sendInv(address, kind string, items [][]byte) {
	inventory := inv{nodeAddress, kind, items}
	payload := utils.GobEncode(inventory)
	request := append(commandToBytes("inv"), payload...)

	sendData(address, request)
}

func handleInv(request []byte, bc *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload inv

	buff.Write(extractCommand(request))
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items

		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Equal(b, blockHash) {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if mempool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}
