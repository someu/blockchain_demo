package cli_test

import (
	"blockchain_demo/cli"
	"testing"
)

func TestPrintChain(t *testing.T) {
	cli := cli.New()
	cli.PrintChain()
}
