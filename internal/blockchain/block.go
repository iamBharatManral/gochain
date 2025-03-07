package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/iamBharatManral/gochain/internal/merkle"
	"github.com/iamBharatManral/gochain/internal/transaction"
)

var gb Block
var once sync.Once

var baseTarget = uint32(math.Pow(2, 256) - 1)

const baseDifficulty = 5

const genesisMessage = "gochain"

type Header struct {
	Index        uint
	Timestamp    time.Duration
	PreviousHash string
	Hash         string
	MerkelRoot   string
	Nounce       uint
	Difficulty   uint32
}

type Data struct {
	Transactions []transaction.Transaction
}

type Block struct {
	Header
	Data
}

func (b Block) Validate() bool {
	currentHash := b.Hash
	calculatedHash := GenerateHash(b)
	return currentHash == calculatedHash
}

func (h Header) Serialize() string {
	serializedHeader := []byte(fmt.Sprintf("%s%s%s%s%s%s", h.Index, h.Timestamp, h.PreviousHash, h.MerkelRoot, h.Nounce, h.Difficulty))
	return fmt.Sprintf("%s", serializedHeader)
}

func (d Data) Serialize() string {
	var data string
	for _, t := range d.Transactions {
		data += t.Serialize()
	}
	return data
}

func (b Block) Serialize() string {
	return fmt.Sprintf("%s%s", b.Header.Serialize(), b.Data.Serialize())
}

func createGenesisBlock() Block {
	once.Do(func() {
		initialTranscation := transaction.Transaction{
			Sender:   "Creator",
			Receiver: "System",
			Amount:   100000,
		}
		data := []transaction.Transaction{initialTranscation}
		header := Header{
			Index:        0,
			Timestamp:    time.Duration(time.Now().UnixMilli()),
			PreviousHash: "",
			MerkelRoot:   fmt.Sprintf("%s", sha256.Sum256([]byte(genesisMessage))),
			Nounce:       0,
			Difficulty:   baseDifficulty,
		}
		block := Block{
			Header: header,
			Data: Data{
				Transactions: data,
			},
		}

		block.Hash = GenerateHash(block)
		gb = block

	})
	return gb
}

func createNewBlock(index uint, prevHash string, trans []transaction.Transaction) Block {
	header := Header{
		Index:        index,
		PreviousHash: prevHash,
		Timestamp:    time.Duration(time.Now().UnixMilli()),
		Difficulty:   baseDifficulty,
		MerkelRoot:   merkle.NewMerkleTree(trans).Root.Hash,
	}
	block := Block{
		Header: header,
		Data: Data{
			Transactions: trans,
		},
	}
	header.Nounce = mineBlock(block)
	block.Header = header
	header.Hash = GenerateHash(block)
	block.Header = header
	return block
}

func mineBlock(b Block) uint {
	var nounce uint
	for {
		header := b.Header
		header.Nounce = nounce
		b.Header = header
		hashed := GenerateHash(b)
		if hashMeetDifficulty(hashed) {
			b.Header.Nounce = nounce
			break
		}
		nounce++
	}
	return nounce
}

func hashMeetDifficulty(hash string) bool {
	hashInt := new(big.Int)
	hashInt.SetString(hash, 16)
	target := calculateTarget()
	return hashInt.Cmp(target) == -1
}

func calculateTarget() *big.Int {
	baseTarget := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
	target := new(big.Int).Rsh(baseTarget, baseDifficulty)
	return target
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

func (b Block) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Block %d:\n", b.Index))
	sb.WriteString(fmt.Sprintf("  Timestamp: %d\n", b.Timestamp))
	sb.WriteString(fmt.Sprintf("  Previous Hash: %s\n", b.PreviousHash))
	sb.WriteString(fmt.Sprintf("  Hash: %s\n", b.Hash))
	sb.WriteString(fmt.Sprintf("  MerkleRoot: %x\n", b.MerkelRoot))
	sb.WriteString(fmt.Sprintf("  Difficulty: %d\n", b.Difficulty))
	sb.WriteString(fmt.Sprintf("  Nounce: %d\n", b.Nounce))
	sb.WriteString("  Transactions:\n")

	for _, tx := range b.Data.Transactions {
		sb.WriteString(fmt.Sprintf("    %s -> %s: %.2f\n", tx.Sender, tx.Receiver, tx.Amount))
	}

	return sb.String()
}
