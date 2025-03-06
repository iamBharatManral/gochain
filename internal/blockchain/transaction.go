package blockchain

import "fmt"

type Transaction struct {
	Sender   string
	Receiver string
	Amount   float64
}

func (t Transaction) Serialize() string {
	serializedData := []byte(fmt.Sprintf("%s%s%s", t.Sender, t.Receiver, t.Amount))
	return fmt.Sprintf("%s", serializedData)
}
