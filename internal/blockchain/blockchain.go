package blockchain

import "strings"

type Blockchain struct {
	blocks []block
}

func New() *Blockchain {
	return &Blockchain{
		blocks: []block{
			createGenesisBlock(),
		},
	}
}

func (bc *Blockchain) AddBlock(trans []transaction) error {
	blockLen := len(bc.blocks) - 1
	prevHash := bc.blocks[blockLen].previousHash
	block := createNewBlock(uint(blockLen)+1, prevHash, trans)
	bc.blocks = append(bc.blocks, block)
	return nil
}

func (bc Blockchain) String() string {
	var sb strings.Builder
	sb.WriteString("Blockchain:\n")
	for _, b := range bc.blocks {
		sb.WriteString(b.String() + "\n")
	}
	return sb.String()
}
