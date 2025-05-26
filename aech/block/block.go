package block

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Transaction represents a single transaction in the AECH system.
// It includes details about the sender, receiver, amount, and a signature.
type Transaction struct {
	// ID is the unique identifier for the transaction, typically a hash of its content.
	ID string
	// Sender is the address of the account sending the transaction.
	Sender string
	// Receiver is the address of the account receiving the transaction.
	Receiver string
	// Amount is the value being transferred in the transaction.
	Amount float64
	// Timestamp is the time at which the transaction was created.
	Timestamp time.Time
	// Signature is the cryptographic signature of the transaction, used for verification.
	// In the current stub implementation, this is a dummy value.
	Signature string
}

// NewTransaction creates and returns a new Transaction instance.
// It initializes the transaction with the sender, receiver, and amount,
// sets the current timestamp, generates a unique ID, and signs the transaction (currently a stub).
func NewTransaction(sender string, receiver string, amount float64) *Transaction {
	tx := &Transaction{
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
		Timestamp: time.Now(),
	}
	// Generate ID: SHA256(sender + receiver + amount + timestamp)
	idData := fmt.Sprintf("%s%s%s%d", sender, receiver, strconv.FormatFloat(amount, 'f', -1, 64), tx.Timestamp.UnixNano())
	hash := sha256.Sum256([]byte(idData))
	tx.ID = hex.EncodeToString(hash[:])

	_ = tx.Sign() // Call Sign to set dummy signature, explicitly ignore error for stub
	return tx
}

// Sign sets a cryptographic signature for the transaction.
// This is a stub implementation and currently sets a dummy signature.
// It returns an error if signing fails (though in the stub, it always returns nil).
func (tx *Transaction) Sign() error {
	// Stub implementation
	tx.Signature = "signed-" + tx.ID
	return nil
}

// Validate checks if the transaction is valid.
// This method currently calls stubbed validation helpers for signature and balance.
// It returns true if the transaction is considered valid, false otherwise.
func (tx *Transaction) Validate() bool {
	// Calls stubbed validation methods
	return tx.validateSignature() && tx.checkBalance()
}

// validateSignature is a stub for actual signature validation.
// It always returns true in the current implementation.
func (tx *Transaction) validateSignature() bool {
	// Stub implementation
	return true
}

// checkBalance is a stub for actual balance checking.
// It always returns true in the current implementation.
func (tx *Transaction) checkBalance() bool {
	// Stub implementation
	return true
}

// Block3D represents a block in the 3D blockchain structure of AECH.
// Each block has spatial coordinates (X, Y, Z), a timestamp, its own hash,
// the hash of the previous block (linking it in a chain), and a list of transactions it contains.
type Block3D struct {
	// X is the coordinate of the block on the X-axis in the 3D plane.
	X int
	// Y is the coordinate of the block on the Y-axis in the 3D plane.
	Y int
	// Z is the coordinate of the block on the Z-axis in the 3D plane.
	Z int
	// Timestamp is the time at which the block was created.
	Timestamp time.Time
	// Hash is the unique cryptographic hash of the block's content.
	Hash string
	// PreviousHash is the hash of the preceding block in its conceptual chain or spatial link.
	// For the Genesis block, this is typically "0" or an empty string.
	PreviousHash string
	// Transactions is a slice of Transaction objects included in this block.
	Transactions []Transaction
}

// GenerateHash calculates and sets the cryptographic hash for the Block3D instance.
// The hash is derived from the block's coordinates (X,Y,Z), timestamp, previous hash,
// and a concatenated string of all transaction IDs within the block.
// The resulting hash is stored in the block's Hash field and also returned.
func (b *Block3D) GenerateHash() string {
	// Concatenate X, Y, Z, Timestamp, PreviousHash
	dataString := fmt.Sprintf("%d%d%d%d%s", b.X, b.Y, b.Z, b.Timestamp.UnixNano(), b.PreviousHash)

	// Concatenate all Transaction IDs
	var txIDs []string
	for _, tx := range b.Transactions {
		txIDs = append(txIDs, tx.ID)
	}
	allTxIDsString := strings.Join(txIDs, "-") // Join with a separator

	// Combine block data with transaction data
	finalData := dataString + allTxIDsString

	hash := sha256.Sum256([]byte(finalData))
	b.Hash = hex.EncodeToString(hash[:]) // Set the block's hash
	return b.Hash                       // Return the hash
}

// ValidateBlock checks if a given Block3D instance is valid according to basic criteria.
// This is a standalone function. In the current stub implementation, it only checks if the block's
// Hash field is non-empty and if its X, Y, and Z coordinates are non-negative.
// TODO: Implement more comprehensive validation, including recalculating and comparing the hash.
func ValidateBlock(b Block3D) bool {
	// For now, just check if Hash is non-empty and X,Y,Z are non-negative.
	if b.Hash == "" || b.X < 0 || b.Y < 0 || b.Z < 0 {
		return false
	}
	// A more advanced validation would re-calculate the hash and compare.
	// currentHash := b.Hash
	// tempBlock := b // Create a copy to avoid modifying the original block's hash field directly
	// tempBlock.Hash = "" // Clear hash to recalculate
	// if tempBlock.GenerateHash() != currentHash { // Now calls the method
	//  return false
	// }
	return true
}

// NewBlock creates and returns a new Block3D instance.
// It initializes the block with the given spatial coordinates (x, y, z),
// the hash of the previous block, and a list of transactions.
// A new timestamp is generated, and the block's hash is computed and set.
func NewBlock(x, y, z int, previousHash string, transactions []Transaction) *Block3D {
	block := &Block3D{
		X:            x,
		Y:            y,
		Z:            z,
		Timestamp:    time.Now(),
		PreviousHash: previousHash,
		Transactions: transactions,
	}
	block.GenerateHash() // Set the hash for the new block by calling the method
	return block
}
