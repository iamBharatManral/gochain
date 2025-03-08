package utils

import (
	"crypto/sha256"
	"fmt"

	"github.com/iamBharatManral/gochain/internal/block"
)

func GenerateHash(b block.Block) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(b.Serialize())))
}
