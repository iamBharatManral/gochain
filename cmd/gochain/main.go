package main

import (
	"fmt"
	"log"

	"github.com/iamBharatManral/gochain/internal/blockchain"
	"github.com/iamBharatManral/gochain/internal/persistence"
	"github.com/iamBharatManral/gochain/internal/transaction"
)

func main() {
	storage, err := persistence.NewBadgerStorage("./data")
	if err != nil {
		log.Fatal("stoage is not initialized")
	}
	bc := blockchain.New(storage)
	trans := []transaction.Transaction{
		{
			Sender:   "Bharat",
			Receiver: "Raul",
			Amount:   10,
		},
	}

	err = bc.AddBlock(trans)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(bc)
	if err = bc.Validate(); err != nil {
		log.Fatal(err.Error())
	}
	defer storage.Close()

}
