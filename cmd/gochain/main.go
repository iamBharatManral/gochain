package main

import (
	"fmt"

	"github.com/iamBharatManral/gochain/internal/blockchain"
)

func main() {
	bc := blockchain.New()
	trans := []blockchain.Transaction{
		{
			Sender:   "Bharat",
			Receiver: "Raul",
			Amount:   10,
		},
	}

	bc.AddBlock(trans)
	fmt.Println(bc)
}
