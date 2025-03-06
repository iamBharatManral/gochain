package blockchain

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var gb Block
var once sync.Once

type Serializer interface {
	Serialize() string
}

type Header struct {
	Index        uint
	Timestamp    time.Duration
	PreviousHash string
	Hash         string
}

type Data struct {
	Transactions []Transaction
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
	serializedHeader := []byte(fmt.Sprintf("%s%s%s", h.Index, h.Timestamp, h.PreviousHash))
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
		initialTranscation := Transaction{
			Sender:   "Creator",
			Receiver: "System",
			Amount:   100000,
		}
		data := []Transaction{initialTranscation}
		header := Header{
			Index:        0,
			Timestamp:    time.Duration(time.Now().UnixMilli()),
			PreviousHash: "",
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

func createNewBlock(index uint, prevHash string, trans []Transaction) Block {
	header := Header{
		Index:        index,
		PreviousHash: prevHash,
		Timestamp:    time.Duration(time.Now().UnixMilli()),
	}
	block := Block{
		Header: header,
		Data: Data{
			Transactions: trans,
		},
	}
	header.Hash = GenerateHash(block)
	block.Header = header
	return block
}

func (b Block) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Block %d:\n", b.Index))
	sb.WriteString(fmt.Sprintf("  Timestamp: %d\n", b.Timestamp))
	sb.WriteString(fmt.Sprintf("  Previous Hash: %s\n", b.PreviousHash))
	sb.WriteString(fmt.Sprintf("  Hash: %s\n", b.Hash))
	sb.WriteString("  Transactions:\n")

	for _, tx := range b.Data.Transactions {
		sb.WriteString(fmt.Sprintf("    %s -> %s: %.2f\n", tx.Sender, tx.Receiver, tx.Amount))
	}

	return sb.String()
}
