package block

import (
	"reflect" // Added for DeepEqual
	"testing"
	"time"
)

func TestNewTransaction(t *testing.T) {
	tx := NewTransaction("sender1", "receiver1", 100.50)
	if tx == nil {
		t.Fatal("NewTransaction returned nil")
	}
	if tx.ID == "" {
		t.Error("Transaction ID is empty")
	}
	if tx.Timestamp <= 0 { // Updated assertion for int64 timestamp
		t.Error("Timestamp should be positive nano seconds")
	}
	expectedSig := []byte("signed-" + tx.ID) // Updated assertion for []byte signature
	if !reflect.DeepEqual(tx.Signature, expectedSig) {
		t.Errorf("Expected Signature %s, got %s", string(expectedSig), string(tx.Signature))
	}
	if tx.Sender != "sender1" || tx.Receiver != "receiver1" || tx.Amount != 100.50 {
		t.Errorf("Transaction fields not set correctly: %+v", tx)
	}
}

func TestTransaction_Validate(t *testing.T) {
	tx := NewTransaction("s", "r", 1)
	if !tx.Validate() {
		t.Error("tx.Validate() returned false, expected true for stub implementation")
	}
}

func TestBlockGenerateHash_WithTransactions(t *testing.T) {
	// Test that block hash changes with different transactions
	tx1 := NewTransaction("s1", "r1", 10)
	// Ensure tx2 is different enough to produce a different ID
	time.Sleep(1 * time.Nanosecond) // Ensure timestamp difference for tx ID generation
	tx2 := NewTransaction("s2", "r2", 20)

	b1 := NewBlock(0, 0, 0, "prevHash", []Transaction{*tx1})
	time.Sleep(1 * time.Nanosecond) // Ensure block timestamp difference
	b2 := NewBlock(0, 0, 0, "prevHash", []Transaction{*tx2})
	time.Sleep(1 * time.Nanosecond)
	b4 := NewBlock(0, 0, 0, "prevHash", nil) // No transactions

	if b1.Hash == b2.Hash {
		t.Errorf("Blocks with different transactions have the same hash (b1, b2). b1: %s, b2: %s", b1.Hash, b2.Hash)
	}
	if b1.Hash == b4.Hash {
		t.Error("Block with transaction has same hash as block without (b1, b4)")
	}
	if b2.Hash == b4.Hash {
		t.Error("Block with different transaction has same hash as block without (b2, b4)")
	}
}
