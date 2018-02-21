package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Block contains data that will be written to the blockchain
// and represents each case when you took your pulse rate
type Block struct {
	Index    int64     `json:"index"`
	Time     time.Time `json:"time"`
	Hash     string    `json:"hash"`
	PrevHash string    `json:"prevHash"`
	Data     string    `json:"data"`
}

// NewGenesisBlock creates new
func NewGenesisBlock(data string) *Block {
	genesisBlock := &Block{0, time.Now(), "", "", data}
	genesisBlock.Hash = genesisBlock.calculateHash()

	return genesisBlock
}

// NewBlock generates a new Block
func NewBlock(parent *Block, data string) (*Block, error) {
	newBlock := &Block{
		Index:    parent.Index + 1,
		Time:     time.Now(),
		Data:     data,
		PrevHash: parent.Hash,
	}

	newBlock.Hash = newBlock.calculateHash()

	return newBlock, nil
}

// String concatenates Index, Time, Data, PrevHash of the Block
func (b *Block) String() string {
	return string(b.Index) + string(b.Time.String()) + b.Data + b.PrevHash
}

// calculateHash returns the SHA256 hash as a string
func (b *Block) calculateHash() string {
	h := sha256.New()
	h.Write([]byte(b.String()))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

func (b *Block) isValid(parent *Block) bool {
	// check if Index has correct value
	if parent.Index+1 != b.Index {
		return false
	}

	// check if PrevHash is correct
	if parent.Hash != b.PrevHash {
		return false
	}

	// check if current Hash is correct
	if b.calculateHash() != b.Hash {
		return false
	}

	return true
}
