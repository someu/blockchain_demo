package cli

import (
	"blockchain_demo/wallet"
	"fmt"
)

func (cli *CLI) createWallet() {
	wallets, _ := wallet.NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Println(address)
}
