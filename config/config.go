package config

import "math"

const (
	DBFile               = "/Users/zhupeng/Code/git/blockchain_go/blockchain.db" // 数据库文件
	WalletFile           = "/Users/zhupeng/Code/git/blockchain_go/wallet.dat"    // 钱包文件
	UTXOBucketName       = "chainstate"                                          // UTXO集 bucket 名
	BlockChainBucketName = "blocks"                                              // 区块链 bucket 名
	TargetBits           = 16                                                    // POW计算难度
	MaxNonce             = math.MaxInt64                                         // POW计算最大nonce
	Protocol             = "tcp"
	NodeVersion          = 1
	CommandLength        = 12
)
