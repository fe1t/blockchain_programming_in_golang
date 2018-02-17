package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const (
	maxNonce   = math.MaxInt64
	targetBits = 4
)

// ProofOfWork structure
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// PrepareData to get bytes stream
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevHash,
			pow.block.HashTransactions(),
			IntToByteStream(pow.block.Timestamp),
			IntToByteStream(int64(targetBits)),
			IntToByteStream(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

// Run to get nonce, hash for the block
func (pow *ProofOfWork) Run() (int, []byte) {
	var (
		hashInt big.Int
		hash    [32]byte
	)
	nonce := 0
	fmt.Printf("Minging the new block")
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		}
		nonce++
	}
	fmt.Printf("\n\n")
	return nonce, hash[:]
}

// NewProofOfWork creates new PoW for the current block
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{b, target}
	return pow
}

// Validate checks if nonce is valid for the hash
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.target) == -1
}