package blockchain

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/iamBharatManral/gochain/internal/block"
	"github.com/iamBharatManral/gochain/internal/persistence"
	"github.com/iamBharatManral/gochain/internal/transaction"
)

type Blockchain struct {
	Blocks  []block.Block
	Storage *persistence.Storage
}

func New(storage persistence.Storage) *Blockchain {
	return &Blockchain{
		Blocks: []block.Block{
			block.CreateGenesisBlock(),
		},
		Storage: &storage,
	}
}

func (bc *Blockchain) AddBlock(trans []transaction.Transaction) error {
	if err := transaction.ValidateTransactions(trans); err != nil {
		return err
	}

	blockLen := len(bc.Blocks)
	prevHash := bc.Blocks[blockLen-1].Hash

	newBlock := block.CreateNewBlock(uint(blockLen), prevHash, trans)

	if err := block.ValidateBlock(newBlock, bc.Blocks[blockLen-1], uint(blockLen)); err != nil {
		return err
	}

	bc.Blocks = append(bc.Blocks, newBlock)
	return nil
}

func (bc *Blockchain) Validate() error {
	var err strings.Builder
	for idx, b := range bc.Blocks {
		if !b.Validate() {
			err.WriteString("block hash does not match with its own previously calculate hash\n")
		}

		if idx == 0 {
			if err := bc.validateGenesisBlock(); err != nil {
				return err
			}
			continue
		}

		if prevBlockIndex, curBlockIndex := bc.Blocks[idx-1].Index, b.Index; !doesBlockIndexCorrectlyIncrement(prevBlockIndex, curBlockIndex) {
			err.WriteString(fmt.Sprintf("block index is not correct, expected: %d, got: %d\n", prevBlockIndex, curBlockIndex))
		}

		if prevBlockHash, prevHashInCurrentBlock := bc.Blocks[idx-1].Hash, b.PreviousHash; !doesHashMatches(prevBlockHash, prevHashInCurrentBlock) {
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
	curGenBlock := bc.Blocks[0]
	if ok := reflect.DeepEqual(curGenBlock, block.CreateGenesisBlock()); !ok {
		return errors.New("")
	}
	return nil
}

func (bc *Blockchain) String() string {
	var sb strings.Builder
	sb.WriteString("Blockchain:\n")
	for _, b := range bc.Blocks {
		sb.WriteString(b.String() + "\n")
	}
	return sb.String()
}
