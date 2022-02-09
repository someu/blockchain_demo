package cli

import (
	"blockchain_demo/blockchain"
	"blockchain_demo/wallet"
	"fmt"
	"log"
)

func (cli *CLI) createBlockchain(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("invalid address")
	}
	bc := blockchain.CreateBlockchain(address)
	defer bc.DB.Close()
	fmt.Println("Done!")
}
