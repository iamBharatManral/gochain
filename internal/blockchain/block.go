package blockchain

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

type Transaction struct {
	Sender   string
	Receiver string
	Amount   float64
}

type Header struct {
	Index        uint
	Timestamp    time.Duration
	PreviousHash string
	Hash         string
}

type Block struct {
	Header
	Data []Transaction
}

func createGenesisBlock() Block {
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
	header.Hash = generateHash(header, data)

	return Block{
		Header: header,
		Data:   data,
	}
}

func serializeTransaction(t Transaction) string {
	serializedData := []byte(fmt.Sprintf("%s%s%s", t.Sender, t.Receiver, t.Amount))
	return fmt.Sprintf("%s", serializedData)
}

func serializeHeader(h Header) string {
	serializedHeader := []byte(fmt.Sprintf("%s%s%s", h.Index, h.Timestamp, h.PreviousHash))
	return fmt.Sprintf("%s", serializedHeader)
}

func serializeData(d []Transaction) string {
	var data string
	for _, t := range d {
		data += serializeTransaction(t)
	}
	return data
}

func generateHash(h Header, d []Transaction) string {
	header, trans := serializeHeader(h), serializeData(d)
	serializedData := fmt.Sprintf("%s%s", header, trans)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(serializedData)))
}

func createNewBlock(index uint, prevHash string, trans []Transaction) Block {
	header := Header{
		Index:        index,
		PreviousHash: prevHash,
		Timestamp:    time.Duration(time.Now().UnixMilli()),
	}
	header.Hash = generateHash(header, trans)
	return Block{
		Data:   trans,
		Header: header,
	}
}

func (b Block) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Block %d:\n", b.Index))
	sb.WriteString(fmt.Sprintf("  Timestamp: %d\n", b.Timestamp))
	sb.WriteString(fmt.Sprintf("  Previous Hash: %s\n", b.PreviousHash))
	sb.WriteString(fmt.Sprintf("  Hash: %s\n", b.Hash))
	sb.WriteString("  Transactions:\n")

	for _, tx := range b.Data {
		sb.WriteString(fmt.Sprintf("    %s -> %s: %.2f\n", tx.Sender, tx.Receiver, tx.Amount))
	}

	return sb.String()
}
