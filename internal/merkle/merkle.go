package merkle

import (
	"crypto/sha256"
	"fmt"

	"github.com/iamBharatManral/gochain/internal/transaction"
)

type MerkleNode struct {
	Hash  string
	Left  *MerkleNode
	Right *MerkleNode
}

func NewMerkleNode(left, right *MerkleNode, hash string) *MerkleNode {
	newNode := MerkleNode{
		Right: right,
		Left:  left,
	}
	if left == nil && right == nil {
		newNode.Hash = hash
	} else {
		concateneted := left.Hash + right.Hash
		hashed := sha256.Sum256([]byte(concateneted))
		newNode.Hash = fmt.Sprintf("%x", hashed)
	}
	return &newNode
}

type MerkleTree struct {
	Root *MerkleNode
}

func NewMerkleTree(trans []transaction.Transaction) *MerkleTree {
	var hashes []string
	for _, t := range trans {
		hashes = append(hashes, t.Serialize())
	}
	return createTreeFromHashes(hashes)
}

func createTreeFromHashes(hashes []string) *MerkleTree {
	var leaves []MerkleNode
	for _, h := range hashes {
		leaves = append(leaves, *NewMerkleNode(nil, nil, h))
	}
	for len(leaves) > 1 {
		var newLevel []MerkleNode
		for i := 0; i < len(leaves); i += 2 {
			if i+1 < len(leaves) {
				parent := NewMerkleNode(&leaves[i], &leaves[i+1], "")
				newLevel = append(newLevel, *parent)
			} else {
				parent := NewMerkleNode(&leaves[i], &leaves[i], "")
				newLevel = append(newLevel, *parent)
			}
		}
		leaves = newLevel
	}
	return &MerkleTree{
		Root: &leaves[0],
	}
}
