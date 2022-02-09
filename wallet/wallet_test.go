package wallet

import (
	"testing"
)

func TestValidateAddress(t *testing.T) {
	addr := NewWallet().GetAddress()
	if !ValidateAddress(string(addr)) {
		t.Error("invalid address")
	}
}
