package persistence

import (
	"github.com/iamBharatManral/gochain/internal/block"
)

type Storage interface {
	SaveBlock(block.Block) error
	GetBlock(uint) (*block.Block, error)
	GetLatestBlock() (*block.Block, error)
	Close() error
}
