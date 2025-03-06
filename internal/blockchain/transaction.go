package blockchain

import (
	"errors"
	"fmt"
	"strings"
)

type Transaction struct {
	Sender   string
	Receiver string
	Amount   float64
}

func (t Transaction) Serialize() string {
	serializedData := []byte(fmt.Sprintf("%s%s%s", t.Sender, t.Receiver, t.Amount))
	return fmt.Sprintf("%s", serializedData)
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
