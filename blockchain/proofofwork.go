package blockchain

import (
	"blockchain_demo/utils"
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

const targetBits = 16
const maxNonce = math.MaxInt64

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	t := big.NewInt(1)
	t.Lsh(t, uint(256-targetBits))
	return &ProofOfWork{
		block:  b,
		target: t,
	}
}

func (p *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			p.block.PrevBlockHash,
			p.block.HashTransactions(),
			utils.IntToHex(uint64(p.block.Timestamp)),
			utils.IntToHex(uint64(targetBits)),
			utils.IntToHex(uint64(nonce)),
		},
		[]byte{},
	)
	return data
}

func (p *ProofOfWork) Run() (int, []byte) {
	var nonce = 0
	var hash [32]byte
	var hashInt big.Int
	for ; nonce < maxNonce; nonce++ {
		d := p.prepareData(nonce)
		hash = sha256.Sum256(d)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(p.target) == -1 {
			return nonce, hash[:]
		}
	}
	return nonce, hash[:]
}

func (p *ProofOfWork) Validate() bool {
	var hash [32]byte
	var hashInt big.Int
	d := p.prepareData(p.block.Nonce)
	hash = sha256.Sum256(d)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(p.target) == -1

}
