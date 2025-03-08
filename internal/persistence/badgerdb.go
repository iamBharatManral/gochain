package persistence

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/iamBharatManral/gochain/internal/block"
)

type BadgerStorage struct {
	db *badger.DB
}

func NewBadgerStorage(dbPath string) (*BadgerStorage, error) {
	err := os.MkdirAll(dbPath, 0700)
	if err != nil {
		return nil, fmt.Errorf("failed to create database directory %v", dbPath)
	}
	db, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		return nil, fmt.Errorf("failed to open BadgerDB: %v", err)
	}
	log.Printf("BadgerDB initialized at %s", dbPath)
	return &BadgerStorage{db}, nil
}

func (bs *BadgerStorage) GetBlock(index uint) (*block.Block, error) {
	var block block.Block
	key := fmt.Sprintf("%s:%d", "block:", index)
	err := bs.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return nil
		}
		err = item.Value(func(val []byte) error {
			err = json.Unmarshal(val, block)
			if err != nil {
				return nil
			}
			return nil
		})
		return nil

	})
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (bs *BadgerStorage) SaveBlock(b block.Block) error {
	err := bs.db.Update(func(txn *badger.Txn) error {
		key := fmt.Sprintf("%s%d", "block:", b.Index)
		err := txn.Set([]byte(key), []byte(b.Serialize()))
		if err != nil {
			return err
		}
		err = txn.Set([]byte("latest_block"), []byte(b.Serialize()))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (bs *BadgerStorage) GetLatestBlock() (*block.Block, error) {
	var block block.Block
	err := bs.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("latest_block"))
		if err != nil {
			return nil
		}
		err = item.Value(func(val []byte) error {
			err = json.Unmarshal(val, block)
			if err != nil {
				return nil
			}
			return nil
		})
		return nil

	})
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (bs *BadgerStorage) Close() error {
	return bs.db.Close()
}
