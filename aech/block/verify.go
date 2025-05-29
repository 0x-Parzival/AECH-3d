package block

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
)

// HashTransactionContent generates the SHA-256 hash of the transaction's core content
// that is signed and verified.
// Note: The user's plan was: tx.Sender + tx.Receiver + string(tx.Amount) + string(tx.Timestamp)
// We will use strconv for numeric types for proper string representation.
// The ID field is NOT part of this hash, as the ID is often derived from this hash.
func HashTransactionContent(tx Transaction) []byte {
	// Consistent string representation of the transaction data to be signed/verified.
	// Order matters here.
	amountStr := strconv.FormatFloat(tx.Amount, 'f', -1, 64)
	timestampStr := strconv.FormatInt(tx.Timestamp, 10)

	input := tx.Sender + tx.Receiver + amountStr + timestampStr
	// Future consideration: Include other fields like Nonce if replay protection is added.
	// Also, if tx.ID were generated *before* signing and meant to be part of signed content,
	// it would be included here. But current plan has ID derived from this hash.

	hash := sha256.Sum256([]byte(input))
	return hash[:]
}

// VerifySignature checks if the given transaction's signature is valid
// for its content, using the provided public key.
func VerifySignature(tx Transaction) error {
	if tx.PubKey == "" {
		return errors.New("missing public key for signature verification")
	}
	if tx.Signature == "" {
		return errors.New("missing signature for verification")
	}

	// Decode public key (expecting compressed hex format from client as per user plan)
	pubKeyBytes, err := hex.DecodeString(tx.PubKey)
	if err != nil {
		return fmt.Errorf("invalid public key hex encoding: %w", err)
	}

	// Unmarshal compressed public key
	// elliptic.P256() is a common choice.
	curve := elliptic.P256()
	x, y := elliptic.UnmarshalCompressed(curve, pubKeyBytes)
	if x == nil { // UnmarshalCompressed returns x=nil if the point is not on the curve or data is invalid
		return errors.New("failed to unmarshal compressed public key, point may not be on curve or data invalid")
	}
	pubKey := ecdsa.PublicKey{Curve: curve, X: x, Y: y}

	// Decode signature (expecting hex format, 64 bytes for P256: 32 for R, 32 for S)
	sigBytes, err := hex.DecodeString(tx.Signature)
	if err != nil {
		return fmt.Errorf("invalid signature hex encoding: %w", err)
	}
	if len(sigBytes) != 64 { // For P256, R and S are 32 bytes each.
		return fmt.Errorf("invalid signature length, expected 64 bytes, got %d", len(sigBytes))
	}

	r := big.NewInt(0).SetBytes(sigBytes[:32])
	s := big.NewInt(0).SetBytes(sigBytes[32:])

	// Hash transaction content that was supposedly signed
	hashedContent := HashTransactionContent(tx)

	// Verify the signature
	if !ecdsa.Verify(&pubKey, hashedContent, r, s) {
		return errors.New("ECDSA signature verification failed")
	}

	return nil // Signature is valid
}
