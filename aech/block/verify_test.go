package block

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	crypto_rand "crypto/rand" // Alias to avoid naming conflict if math/rand is used elsewhere
	"encoding/hex"
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"
	"time"
)

// TestHashTransactionContent ensures that the hashing function produces consistent results
// and that changes to relevant fields result in different hashes.
func TestHashTransactionContent(t *testing.T) {
	tx1 := Transaction{
		Sender:    "Alice",
		Receiver:  "Bob",
		Amount:    10.0,
		Timestamp: 1678886400000, // Fixed timestamp for predictable hash
		// ID, PubKey, Signature are not part of HashTransactionContent
	}
	hash1 := HashTransactionContent(tx1)

	// Test for consistency
	hash1_again := HashTransactionContent(tx1)
	if !reflect.DeepEqual(hash1, hash1_again) {
		t.Errorf("HashTransactionContent is not deterministic. Got %x and %x", hash1, hash1_again)
	}

	// Test that changing sender changes hash
	tx2 := tx1
	tx2.Sender = "Charlie"
	hash2 := HashTransactionContent(tx2)
	if reflect.DeepEqual(hash1, hash2) {
		t.Errorf("Changing Sender did not change hash. Hash1: %x, Hash2: %x", hash1, hash2)
	}

	// Test that changing receiver changes hash
	tx3 := tx1
	tx3.Receiver = "David"
	hash3 := HashTransactionContent(tx3)
	if reflect.DeepEqual(hash1, hash3) {
		t.Errorf("Changing Receiver did not change hash. Hash1: %x, Hash3: %x", hash1, hash3)
	}

	// Test that changing amount changes hash
	tx4 := tx1
	tx4.Amount = 20.0
	hash4 := HashTransactionContent(tx4)
	if reflect.DeepEqual(hash1, hash4) {
		t.Errorf("Changing Amount did not change hash. Hash1: %x, Hash4: %x", hash1, hash4)
	}

	// Test that changing timestamp changes hash
	tx5 := tx1
	tx5.Timestamp = 1678886500000
	hash5 := HashTransactionContent(tx5)
	if reflect.DeepEqual(hash1, hash5) {
		t.Errorf("Changing Timestamp did not change hash. Hash1: %x, Hash5: %x", hash1, hash5)
	}

	t.Logf("TestHashTransactionContent successful. Initial hash: %x", hash1)
}

// TestSignAndVerifyTransaction tests the full cycle of signing a transaction (using a modified NewTransaction for testing)
// and then verifying its signature.
func TestSignAndVerifyTransaction(t *testing.T) {
	// 1. Generate a new key pair for this test
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), crypto_rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}
	pubKeyCompressedBytes := elliptic.MarshalCompressed(elliptic.P256(), privKey.X, privKey.Y)
	pubKeyHex := hex.EncodeToString(pubKeyCompressedBytes)

	// 2. Create a transaction instance (manually, as NewTransaction now does too much for this direct test)
	tx := Transaction{
		Sender:    "TestSender", // Conceptual sender name
		Receiver:  "TestReceiver",
		Amount:    100.0,
		Timestamp: time.Now().UnixNano(),
		PubKey:    pubKeyHex, // Assign the generated public key
		// ID and Signature will be populated by our signing process
	}

	// 3. Populate ID using HashTransactionContent (as NewTransaction would)
	idHashBytes := HashTransactionContent(tx)
	tx.ID = hex.EncodeToString(idHashBytes)

	// 4. Sign the transaction ID (content hash) using the private key
	r, s, err := ecdsa.Sign(crypto_rand.Reader, privKey, idHashBytes)
	if err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	signatureR := make([]byte, 32)
	signatureS := make([]byte, 32)
	copy(signatureR[32-len(rBytes):], rBytes)
	copy(signatureS[32-len(sBytes):], sBytes)
	signatureBytes := append(signatureR, signatureS...)
	tx.Signature = hex.EncodeToString(signatureBytes)

	// 5. Verify the signature
	err = VerifySignature(tx)
	if err != nil {
		t.Errorf("Valid signature failed verification: %v", err)
		t.Logf("Transaction details: %+v", tx)
	} else {
		t.Logf("Successfully signed and verified transaction: ID=%s", tx.ID)
	}

	// --- Test Case: Tampered Data ---
	t.Run("TamperedData", func(t *testing.T) {
		tamperedTx := tx
		tamperedTx.Amount = 200.0 // Modify amount after signing
		err := VerifySignature(tamperedTx)
		if err == nil {
			t.Error("Signature verification passed for tampered data, but it should fail.")
		} else if !strings.Contains(err.Error(), "ECDSA signature verification failed") {
			t.Errorf("Expected ECDSA verification error for tampered data, got: %v", err)
		} else {
			t.Logf("Correctly failed to verify tampered data: %v", err)
		}
	})

	// --- Test Case: Invalid Signature ---
	t.Run("InvalidSignature", func(t *testing.T) {
		invalidSigTx := tx
		invalidSigTx.Signature = "abcdef123456" // Clearly invalid signature
		err := VerifySignature(invalidSigTx)
		if err == nil {
			t.Error("Signature verification passed for invalid signature format, but it should fail.")
		} else if !strings.Contains(err.Error(), "invalid signature hex encoding") && !strings.Contains(err.Error(), "invalid signature length") {
			// Check for either hex decoding error or length error
			t.Errorf("Expected signature encoding/length error, got: %v", err)
		} else {
			t.Logf("Correctly failed to verify invalid signature: %v", err)
		}
	})
	
	// --- Test Case: Invalid Signature (Correct Length, Wrong Content) ---
	t.Run("WrongSignatureContent", func(t *testing.T) {
		wrongSigTx := tx
		// Create a signature that is hex and correct length but not for this data
		dummySigBytes := make([]byte, 64)
		_, _ = crypto_rand.Read(dummySigBytes) // Fill with random bytes
		wrongSigTx.Signature = hex.EncodeToString(dummySigBytes)

		err := VerifySignature(wrongSigTx)
		if err == nil {
			t.Error("Signature verification passed for wrong signature content, but it should fail.")
		} else if !strings.Contains(err.Error(), "ECDSA signature verification failed") {
			t.Errorf("Expected ECDSA verification error for wrong signature, got: %v", err)
		} else {
			t.Logf("Correctly failed to verify wrong signature content: %v", err)
		}
	})


	// --- Test Case: Malformed Public Key (Invalid Hex) ---
	t.Run("MalformedPubKeyHex", func(t *testing.T) {
		malformedKeyTx := tx
		malformedKeyTx.PubKey = "not-a-hex-string"
		err := VerifySignature(malformedKeyTx)
		if err == nil {
			t.Error("Signature verification passed for malformed public key (hex), but it should fail.")
		} else if !strings.Contains(err.Error(), "invalid public key hex encoding") {
			t.Errorf("Expected public key hex encoding error, got: %v", err)
		} else {
			t.Logf("Correctly failed to verify with malformed public key (hex): %v", err)
		}
	})
	
	// --- Test Case: Malformed Public Key (Not on Curve) ---
	t.Run("MalformedPubKeyNotOnCurve", func(t *testing.T) {
		malformedKeyTx_NotOnCurve := tx
		// Create some random bytes that are unlikely to be a valid compressed public key on P256
		randomBytes := make([]byte, 33) // Compressed keys for P256 are 33 bytes (0x02/0x03 + X coord)
		_, _ = crypto_rand.Read(randomBytes)
		randomBytes[0] = 0x02 // Set a valid prefix
		malformedKeyTx_NotOnCurve.PubKey = hex.EncodeToString(randomBytes)
		
		// It's hard to guarantee this random data will *always* fail unmarshal in a specific way,
		// but it's highly probable it won't be a valid point.
		err := VerifySignature(malformedKeyTx_NotOnCurve)
		if err == nil {
			t.Errorf("Signature verification passed for malformed public key (not on curve), but it should fail. PubKey: %s", malformedKeyTx_NotOnCurve.PubKey)
		} else if !strings.Contains(err.Error(), "failed to unmarshal compressed public key") {
			t.Errorf("Expected 'failed to unmarshal compressed public key' error, got: %v", err)
		} else {
			t.Logf("Correctly failed to verify with malformed public key (not on curve): %v", err)
		}
	})


	// --- Test Case: Missing Signature ---
	t.Run("MissingSignature", func(t *testing.T) {
		missingSigTx := tx
		missingSigTx.Signature = ""
		err := VerifySignature(missingSigTx)
		if err == nil {
			t.Error("Signature verification passed with missing signature, but it should fail.")
		} else if !strings.Contains(err.Error(), "missing signature for verification") {
			t.Errorf("Expected 'missing signature' error, got: %v", err)
		} else {
			t.Logf("Correctly failed to verify with missing signature: %v", err)
		}
	})

	// --- Test Case: Missing Public Key ---
	t.Run("MissingPubKey", func(t *testing.T) {
		missingPubKeyTx := tx
		missingPubKeyTx.PubKey = ""
		err := VerifySignature(missingPubKeyTx)
		if err == nil {
			t.Error("Signature verification passed with missing public key, but it should fail.")
		} else if !strings.Contains(err.Error(), "missing public key for signature verification") {
			t.Errorf("Expected 'missing public key' error, got: %v", err)
		} else {
			t.Logf("Correctly failed to verify with missing public key: %v", err)
		}
	})
}

// Example of using NewTransaction for a full integration test, though it generates its own keys.
// This test ensures NewTransaction itself produces verifiable transactions.
func TestNewTransaction_SignAndVerify(t *testing.T) {
	// NewTransaction now generates its own keys and signs.
	// The senderName param is conceptual.
	tx := NewTransaction("GeneratedSender", "GeneratedReceiver", 50.0)
	if tx == nil {
		t.Fatal("NewTransaction returned nil")
	}

	// Verify the signature of the transaction created by NewTransaction
	err := VerifySignature(*tx)
	if err != nil {
		t.Errorf("Signature verification failed for transaction created by NewTransaction: %v", err)
		t.Logf("Transaction details: %+v", tx)
	} else {
		t.Logf("Successfully verified transaction created by NewTransaction: ID=%s", tx.ID)
	}
}

```
