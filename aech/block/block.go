// Package block defines the structure and operations for blocks and transactions in the AECH system.
package block

import (
	"crypto/ecdsa"          // Added
	"crypto/elliptic"       // Added
	crypto_rand "crypto/rand" // Added with alias
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log" // Added
	"strconv"
	"strings"
	"time"
)

// Transaction represents a single transaction in the AECH system.
// It includes details about the sender, receiver, amount, and a signature.
type Transaction struct {
	ID        string  `json:"id"`        // Hash of tx content (will be set by server after verification for now)
	Timestamp int64   `json:"timestamp"` // Epoch time (client-set)
	Sender    string  `json:"sender"`    // Sender’s address or identifier (client-set)
	Receiver  string  `json:"receiver"`  // Receiver’s address (client-set)
	Amount    float64 `json:"amount"`    // Amount to transfer (client-set, kept as float64)
	Signature string  `json:"signature"` // Hex-encoded signature (client-set)
	PubKey    string  `json:"pubKey"`    // Hex-encoded sender’s public key (client-set)
}

// NewTransaction creates and returns a new Transaction instance.
// It initializes the transaction with the sender, receiver, and amount,
// sets the current timestamp (UnixNano), generates a unique ID, and signs the transaction (currently a stub).
// NewTransaction creates and returns a new cryptographically signed Transaction instance.
// This version is primarily for server-side generation (e.g., for testing or coinbase).
// The 'senderName' is a conceptual identifier; a new key pair is generated for this transaction.
func NewTransaction(senderName string, receiverName string, amount float64) *Transaction {
	// 1. Generate ECDSA key pair for the sender
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), crypto_rand.Reader)
	if err != nil {
		log.Printf("ERROR: NewTransaction - Failed to generate ECDSA key pair: %v", err)
		return nil
	}
	pubKeyCompressedBytes := elliptic.MarshalCompressed(elliptic.P256(), privKey.X, privKey.Y)
	pubKeyHex := hex.EncodeToString(pubKeyCompressedBytes)

	tx := &Transaction{
		Sender:    senderName,            // User-friendly sender name or identifier
		Receiver:  receiverName,
		Amount:    amount,
		Timestamp: time.Now().UnixNano(),
		PubKey:    pubKeyHex,             // Hex-encoded compressed public key
		// ID and Signature will be set next.
	}

	// 2. Set Transaction ID (hash of content as defined by HashTransactionContent)
	// HashTransactionContent uses tx.Sender, tx.Receiver, tx.Amount, tx.Timestamp.
	// It does NOT include PubKey or Signature itself in the hash.
	idHashBytes := HashTransactionContent(*tx) // This function is in verify.go (same package)
	tx.ID = hex.EncodeToString(idHashBytes)

	// 3. Sign the transaction ID (which is the hash of the content)
	r, s, err := ecdsa.Sign(crypto_rand.Reader, privKey, idHashBytes) // Sign the ID hash
	if err != nil {
		log.Printf("ERROR: NewTransaction - Failed to sign transaction %s: %v", tx.ID, err)
		return nil
	}

	// Convert r, s to fixed-size byte slices (32 bytes each for P256) and concatenate
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	
	// Pad r and s to 32 bytes. This is crucial for consistent signature representation.
	signatureR := make([]byte, 32)
	signatureS := make([]byte, 32)
	copy(signatureR[32-len(rBytes):], rBytes)
	copy(signatureS[32-len(sBytes):], sBytes)
	
	signatureBytes := append(signatureR, signatureS...)
	tx.Signature = hex.EncodeToString(signatureBytes)

	// The old call to tx.Sign() is no longer needed as we've implemented actual signing.
	// This line should be removed if present from the old NewTransaction:
	// _ = tx.Sign() 

	log.Printf("New signed transaction created: ID=%s, Sender=%s, PubKey=%s", tx.ID, tx.Sender, tx.PubKey)
	return tx
}

// Sign sets a cryptographic signature for the transaction.
// This is a stub implementation and currently sets a dummy signature.
// The signature is stored as []byte.
// It returns an error if signing fails (though in the stub, it always returns nil).
func (tx *Transaction) Sign() error {
	// Old dummy signing logic - to be replaced by client-side signing
	// and new server-side verification.
	// tx.Signature = []byte("signed-" + tx.ID) // Store as []byte
	return nil // Or perhaps return an error indicating it's deprecated
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
// and a comprehensive string representation of all transactions within the block.
// The resulting hash is stored in the block's Hash field and also returned.
func (b *Block3D) GenerateHash() string {
	// Concatenate X, Y, Z, Timestamp, PreviousHash
	dataString := fmt.Sprintf("%d%d%d%d%s", b.X, b.Y, b.Z, b.Timestamp.UnixNano(), b.PreviousHash)

	// Create a consistent string representation for each transaction
	var transactionRepresentations []string
	for _, tx := range b.Transactions {
		// including the new types for Timestamp and Signature.
		txString := fmt.Sprintf("%s&S:%s&R:%s&A:%f&T:%d&Sig:%s", // Using '&' and field letters as separators
			tx.ID,
			tx.Sender,
			tx.Receiver,
			tx.Amount,
			tx.Timestamp, // This is int64
			tx.Signature, // Signature is now string (hex-encoded)
		)
		transactionRepresentations = append(transactionRepresentations, txString)
	}
	// Join these transaction strings with another unique separator
	allTransactionsString := strings.Join(transactionRepresentations, "||")

	// Combine block data with transaction data
	finalDataToHash := dataString + "||TXS||" + allTransactionsString // Add a clear separator for transactions part

	hash := sha256.Sum256([]byte(finalDataToHash))
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
