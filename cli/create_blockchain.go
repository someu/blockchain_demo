package cli

import (
	"blockchain_demo/blockchain"
	"fmt"
	"log"
)

func (cli *CLI) createBlockchain(address string) {
	if !blockchain.ValidateAddress(address) {
		log.Panic("invalid address")
	}
	bc := blockchain.CreateBlockchain(address)
	defer bc.DB.Close()
	fmt.Println("Done!")
}
