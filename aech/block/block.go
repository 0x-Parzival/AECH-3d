package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Transaction struct definition
type Transaction struct {
	Data []byte // Simple placeholder for transaction data
}

// Block3D struct definition
type Block3D struct {
	X            int
	Y            int
	Z            int
	Timestamp    time.Time
	Hash         string
	PreviousHash string
	Transactions []Transaction
}

// GenerateHash calculates and sets the hash for a Block3D
func GenerateHash(b *Block3D) string {
	var txHashes [][]byte
	for _, tx := range b.Transactions {
		txHash := sha256.Sum256(tx.Data)
		txHashes = append(txHashes, txHash[:])
	}

	data := fmt.Sprintf("%d%d%d%d%s%x",
		b.X, b.Y, b.Z,
		b.Timestamp.UnixNano(),
		b.PreviousHash,
		bytes.Join(txHashes, []byte{}),
	)

	hash := sha256.Sum256([]byte(data))
	b.Hash = hex.EncodeToString(hash[:])
	return b.Hash
}

// ValidateBlock checks if a Block3D is valid
func ValidateBlock(b Block3D) bool {
	// TODO: Implement more comprehensive validation
	// For now, just check if Hash is non-empty and X,Y,Z are non-negative.
	if b.Hash == "" || b.X < 0 || b.Y < 0 || b.Z < 0 {
		return false
	}
	// A more advanced validation would re-calculate the hash and compare.
	// currentHash := b.Hash
	// b.Hash = "" // Clear hash to recalculate
	// if GenerateHash(&b) != currentHash {
	//  return false
	// }
	return true
}

// NewBlock creates and returns a new Block3D
func NewBlock(x, y, z int, previousHash string, transactions []Transaction) *Block3D {
	block := &Block3D{
		X:            x,
		Y:            y,
		Z:            z,
		Timestamp:    time.Now(),
		PreviousHash: previousHash,
		Transactions: transactions,
	}
	GenerateHash(block) // Set the hash for the new block
	return block
}
