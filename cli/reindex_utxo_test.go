package cli_test

import (
	"blockchain_demo/cli"
	"testing"
)

func TestReindexUTXO(t *testing.T) {
	cli := cli.New()
	cli.ReindexUTXO()
}
