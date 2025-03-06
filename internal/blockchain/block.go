package blockchain

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

type transaction struct {
	sender   string
	receiver string
	amount   float64
}

type header struct {
	index        uint
	timestamp    time.Duration
	previousHash string
	hash         string
}

type block struct {
	header
	data []transaction
}

func createGenesisBlock() block {
	initialTranscation := transaction{
		sender:   "Creator",
		receiver: "System",
		amount:   100000,
	}
	data := []transaction{initialTranscation}
	header := header{
		index:        0,
		timestamp:    time.Duration(time.Now().UnixMilli()),
		previousHash: "",
	}
	header.hash = generateHash(header, data)

	return block{
		header: header,
		data:   data,
	}
}

func serializeTransaction(t transaction) string {
	serializedData := []byte(fmt.Sprintf("%s%s%s", t.sender, t.receiver, t.amount))
	return fmt.Sprintf("%s", serializedData)
}

func serializeHeader(h header) string {
	serializedHeader := []byte(fmt.Sprintf("%s%s%s", h.index, h.timestamp, h.previousHash))
	return fmt.Sprintf("%s", serializedHeader)
}

func serializeData(d []transaction) string {
	var data string
	for _, t := range d {
		data += serializeTransaction(t)
	}
	return data
}

func generateHash(h header, d []transaction) string {
	header, trans := serializeHeader(h), serializeData(d)
	serializedData := fmt.Sprintf("%s%s", header, trans)
	return fmt.Sprintf("%s", sha256.Sum256([]byte(serializedData)))
}

func createNewBlock(index uint, prevHash string, trans []transaction) block {
	header := header{
		index:        index,
		previousHash: prevHash,
		timestamp:    time.Duration(time.Now().UnixMilli()),
	}
	header.hash = generateHash(header, trans)
	return block{
		data:   trans,
		header: header,
	}
}

func (b block) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Block %d:\n", b.index))
	sb.WriteString(fmt.Sprintf("  Timestamp: %v\n", b.timestamp))
	sb.WriteString(fmt.Sprintf("  Previous Hash: %s\n", b.previousHash))
	sb.WriteString(fmt.Sprintf("  Hash: %s\n", b.hash))
	sb.WriteString("  Transactions:\n")

	for _, tx := range b.data {
		sb.WriteString(fmt.Sprintf("    %s -> %s: %.2f\n", tx.sender, tx.receiver, tx.amount))
	}

	return sb.String()
}
