package cli

import (
	"blockchain_demo/blockchain"
	"fmt"
	"log"
)

func (cli *CLI) Send(from, to string, amount int) {
	if !blockchain.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !blockchain.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}
	bc := blockchain.NewBlockChain()
	defer bc.DB.Close()
	UTXOSet := blockchain.UTXOSet{
		BlockChain: bc,
	}
	tx := blockchain.NewUTXOTransaction(from, to, amount, &UTXOSet)
	cbTx := blockchain.NewCoinbaseTx(from, "reward")
	block := bc.MineBlock([]*blockchain.Transaction{tx, cbTx})
	// fix update
	// fix reindexutxo
	UTXOSet.Update(*block)
	fmt.Println("Done!")
}
