package server

import (
	"blockchain_demo/blockchain"
	"blockchain_demo/utils"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"log"
)

type getdata struct {
	AddrFrom string
	Type     string
	ID       []byte
}

func sendGetData(address, kind string, id []byte) {
	payload := utils.GobEncode(getdata{nodeAddress, kind, id})
	request := append(commandToBytes("getdata"), payload...)

	sendData(address, request)
}

func handleGetData(request []byte, bc *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload getdata

	buff.Write(extractCommand(request))
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}

		sendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := mempool[txID]

		sendTx(payload.AddrFrom, &tx)
	}
}
