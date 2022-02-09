package cli

import (
	"blockchain_demo/blockchain"
	"blockchain_demo/transaction"
	"blockchain_demo/wallet"
	"fmt"
	"log"
)

func (cli *CLI) Send(from, to string, amount int) {
	if !wallet.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !wallet.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}
	bc := blockchain.NewBlockChain()
	defer bc.DB.Close()
	UTXOSet := blockchain.UTXOSet{
		BlockChain: bc,
	}
	tx := blockchain.NewUTXOTransaction(from, to, amount, &UTXOSet)
	cbTx := transaction.NewCoinbaseTx(from, "reward")
	block := bc.MineBlock([]*transaction.Transaction{tx, cbTx})
	// fix update
	// fix reindexutxo
	UTXOSet.Update(*block)
	fmt.Println("Done!")
}
