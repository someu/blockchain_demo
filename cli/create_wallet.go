package cli

import (
	"blockchain_demo/blockchain"
	"fmt"
)

func (cli *CLI) createWallet() {
	wallets, _ := blockchain.NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Println(address)
}
