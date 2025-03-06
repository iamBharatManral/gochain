package blockchain

import (
	"crypto/sha256"
	"fmt"
)

func GenerateHash(b Block) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(b.Serialize())))
}
