package blockchain

import "fmt"

// Chain is a chain of Blocks
type Chain []*Block

// NewBlockchain creates a new Blockchain
func NewBlockchain() Chain {
	return Chain{NewGenesisBlock("")}
}

// Append appends new block to chain
func (bc Chain) Append(b *Block) error {
	if !b.isValid(bc[len(bc)-1]) {
		return fmt.Errorf("Invalid block")
	}

	bc.replace(append(bc, b))

	return nil
}

func (bc Chain) replace(chain Chain) {
	if len(bc) < len(chain) {
		bc = chain
	}
}
