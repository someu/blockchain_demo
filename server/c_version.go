package server

import (
	"blockchain_demo/blockchain"
	"blockchain_demo/config"
	"blockchain_demo/utils"
	"bytes"
	"encoding/gob"
	"log"
)

type version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

func sendVersion(addr string, bc *blockchain.BlockChain) {
	bh := bc.GetBestHeight()
	payload := utils.GobEncode(version{
		Version:    config.NodeVersion,
		BestHeight: bh,
		AddrFrom:   nodeAddress},
	)

	request := append(commandToBytes("version"), payload...)

	sendData(addr, request)
}

func handleVersion(request []byte, bc *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload version

	buff.Write(extractCommand(request))
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	}

	if !nodeIsKnown(payload.AddrFrom) {
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
}
