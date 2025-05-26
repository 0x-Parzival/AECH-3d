package block

import (
	"testing"
	"time"
    "strings"
)

func TestNewTransaction(t *testing.T) {
	tx := NewTransaction("sender1", "receiver1", 100.50)
	if tx == nil {
		t.Fatal("NewTransaction returned nil")
	}
	if tx.ID == "" {
		t.Error("Transaction ID is empty")
	}
	if tx.Timestamp.IsZero() {
		t.Error("Transaction Timestamp is zero")
	}
	expectedSignaturePrefix := "signed-"
	if !strings.HasPrefix(tx.Signature, expectedSignaturePrefix) {
		t.Errorf("Transaction Signature '%s' does not have expected prefix '%s'", tx.Signature, expectedSignaturePrefix)
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
	// b3 := NewBlock(0, 0, 0, "prevHash", []Transaction{*tx1}) // Same content as b1, but different timestamp due to NewBlock
    time.Sleep(1 * time.Nanosecond)
    b4 := NewBlock(0,0,0, "prevHash", nil) // No transactions

	if b1.Hash == b2.Hash {
		t.Errorf("Blocks with different transactions have the same hash (b1, b2). b1: %s, b2: %s", b1.Hash, b2.Hash)
	}
	// b1.Hash will not be equal to b3.Hash because NewBlock assigns time.Now()
	// so their transaction lists are the same, but block timestamps differ.
    // What we can check is if the transaction part of hash matters.
    if b1.Hash == b4.Hash {
        t.Error("Block with transaction has same hash as block without (b1, b4)")
    }
    if b2.Hash == b4.Hash {
        t.Error("Block with different transaction has same hash as block without (b2, b4)")
    }
    // To properly test if identical tx lists (and other fields) give same hash,
    // we'd need to control block timestamp precisely.
    // For now, showing different tx lists -> different hashes is key.
    // And block with tx list -> different hash from block with no tx list.
}
