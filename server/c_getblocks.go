package server

import (
	"blockchain_demo/blockchain"
	"blockchain_demo/utils"
	"bytes"
	"encoding/gob"
	"log"
)

type getblocks struct {
	AddrFrom string
}

func sendGetBlocks(address string) {
	payload := utils.GobEncode(getblocks{nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func handleGetBlocks(request []byte, bc *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload getblocks

	buff.Write(extractCommand(request))
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, "block", blocks)
}
