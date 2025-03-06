package blockchain

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

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
	if err := validateTransactions(trans); err != nil {
		return err
	}

	blockLen := len(bc.blocks) - 1
	prevHash := bc.blocks[blockLen].Hash

	newBlock := createNewBlock(uint(blockLen)+1, prevHash, trans)

	if err := validateBlock(bc, newBlock, uint(blockLen)); err != nil {
		return err
	}

	bc.blocks = append(bc.blocks, newBlock)
	return nil
}

func validateBlock(bc *Blockchain, b Block, currentBlockLength uint) error {
	if b.Index != currentBlockLength+1 {
		return errors.New(fmt.Sprintf("new block index is not correct, expected: %d, got: %d", currentBlockLength+1, b.Index))
	}

	currentBlockHash := bc.blocks[currentBlockLength].Hash

	if b.PreviousHash != currentBlockHash {
		return errors.New(fmt.Sprintf("previous block hash does not match in new block, expected: %x, got: %x", currentBlockHash, b.PreviousHash))
	}

	return nil
}

func validateTransactions(ts []Transaction) error {
	for _, t := range ts {
		if err := validateTransaction(t); err != nil {
			return err
		}
	}
	return nil
}

func validateTransaction(t Transaction) error {
	var sb strings.Builder
	if t.Sender == "" {
		sb.WriteString(fmt.Sprintf("%s is empty\n", t.Sender))
	}
	if t.Receiver == "" {
		sb.WriteString(fmt.Sprintf("%s is empty\n", t.Receiver))
	}
	if t.Amount < 0 {
		sb.WriteString(fmt.Sprintf("%s is negative\n", t.Amount))
	}
	err := sb.String()
	if err == "" {
		return nil
	}
	return errors.New(err)

}

func (bc *Blockchain) Validate() error {
	var err strings.Builder
	for idx, b := range bc.blocks {
		if !b.Validate() {
			err.WriteString("block hash does not match with its own previously calculate hash\n")
		}

		if idx == 0 {
			if err := bc.validateGenesisBlock(); err != nil {
				return err
			}
			continue
		}

		if prevBlockIndex, curBlockIndex := bc.blocks[idx-1].Index, b.Index; !doesBlockIndexCorrectlyIncrement(prevBlockIndex, curBlockIndex) {
			err.WriteString(fmt.Sprintf("block index is not correct, expected: %d, got: %d\n", prevBlockIndex, curBlockIndex))
		}

		if prevBlockHash, prevHashInCurrentBlock := bc.blocks[idx-1].Hash, b.PreviousHash; !doesHashMatches(prevBlockHash, prevHashInCurrentBlock) {
			err.WriteString(fmt.Sprintf("previous hash does not match, expected: %x, got: %x\n", prevBlockHash, prevHashInCurrentBlock))
		}
	}
	if err.String() != "" {
		return errors.New(err.String())
	}
	return nil
}

func doesBlockIndexCorrectlyIncrement(prev, cur uint) bool {
	return prev+1 == cur
}

func doesHashMatches(hashOfPrevBlock, hashInCurBlock string) bool {
	return hashOfPrevBlock == hashInCurBlock
}

func (bc *Blockchain) validateGenesisBlock() error {
	curGenBlock := bc.blocks[0]
	if ok := reflect.DeepEqual(curGenBlock, createGenesisBlock()); !ok {
		return errors.New("")
	}
	return nil
}

func (bc *Blockchain) String() string {
	var sb strings.Builder
	sb.WriteString("Blockchain:\n")
	for _, b := range bc.blocks {
		sb.WriteString(b.String() + "\n")
	}
	return sb.String()
}
