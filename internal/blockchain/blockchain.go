package blockchain

import "strings"

type Blockchain struct {
	blocks []Block
}

func New() *Blockchain {
	return &Blockchain{
		blocks: []Block{
			createGenesisBlock(),
		},
	}
}

func (bc *Blockchain) AddBlock(trans []Transaction) error {
	blockLen := len(bc.blocks) - 1
	prevHash := bc.blocks[blockLen].Hash
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
