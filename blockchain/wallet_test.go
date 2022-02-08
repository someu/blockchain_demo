package blockchain_test

import (
	"blockchain_demo/blockchain"
	"testing"
)

func TestValidateAddress(t *testing.T) {
	addr := blockchain.NewWallet().GetAddress()
	if !blockchain.ValidateAddress(string(addr)) {
		t.Error("invalid address")
	}
}
