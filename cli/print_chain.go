package cli

import (
	"blockchain_demo/blockchain"
	"fmt"
	"strconv"
)

func (cli *CLI) PrintChain() {
	bc := blockchain.NewBlockChain()
	defer bc.DB.Close()
	iter := bc.Iterator()
	for {
		block := iter.Next()
		if block == nil {
			return
		}
		fmt.Printf("Prev: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("PoW : %s\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx.String())
		}
		fmt.Println()
	}
}
