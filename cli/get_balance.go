package cli

import (
	"blockchain_demo/blockchain"
	"blockchain_demo/utils"
	"fmt"
)

func (cli *CLI) getBalance(address string) {
	bc := blockchain.NewBlockChain()
	defer bc.DB.Close()
	UTXOSet := &blockchain.UTXOSet{
		BlockChain: bc,
	}
	balance := 0
	pubKeyHash := utils.Base58Decode([]byte(address))
	UTXOs := UTXOSet.FindUTXO(pubKeyHash[1 : len(pubKeyHash)-4])

	for _, UTXO := range UTXOs {
		balance += UTXO.Value
	}
	fmt.Printf("%s's found is %d\n", address, balance)
}
