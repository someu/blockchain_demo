package cli

import (
	"blockchain_demo/blockchain"
	"fmt"
)

func (cli *CLI) ReindexUTXO() {
	bc := blockchain.NewBlockChain()
	defer bc.DB.Close()
	utxo := blockchain.UTXOSet{
		BlockChain: bc,
	}
	utxo.Reindex()
	fmt.Println("done!")
}
