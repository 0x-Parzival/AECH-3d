package txpool

import (
	"aech/block"
	"testing"
)

func TestNewTxPool(t *testing.T) {
	tp := NewTxPool()
	if tp == nil {
		t.Fatal("NewTxPool returned nil")
	}
	if tp.pending == nil {
		t.Error("TxPool pending map not initialized")
	}
}

func TestTxPool_AddTransaction(t *testing.T) {
	tp := NewTxPool()
	tx1 := block.NewTransaction("s1", "r1", 10)

	err := tp.AddTransaction(*tx1)
	if err != nil {
		t.Fatalf("AddTransaction failed for tx1: %v", err)
	}
	if len(tp.pending) != 1 {
		t.Errorf("Expected 1 pending transaction, got %d", len(tp.pending))
	}
	if _, exists := tp.pending[tx1.ID]; !exists {
		t.Errorf("tx1 not found in pending map by ID")
	}

	err = tp.AddTransaction(*tx1) // Add duplicate
	if err == nil {
		t.Error("Expected error when adding duplicate transaction, got nil")
	}
}

func TestTxPool_GetPendingTransactions(t *testing.T) {
	tp := NewTxPool()
	txs := []*block.Transaction{
		block.NewTransaction("s1", "r1", 10),
		block.NewTransaction("s2", "r2", 20),
		block.NewTransaction("s3", "r3", 30),
	}

	// Test with empty pool
	if len(tp.GetPendingTransactions(0)) != 0 {
		t.Errorf("Expected GetPendingTransactions(0) to return 0 transactions for an empty pool, got %d", len(tp.GetPendingTransactions(0)))
	}
	if len(tp.GetPendingTransactions(-1)) != 0 {
		t.Error("Expected GetPendingTransactions(-1) to return 0 transactions for an empty pool")
	}

	// Add transactions to the pool
	for _, tx := range txs {
		if err := tp.AddTransaction(*tx); err != nil {
			t.Fatalf("Error adding tx %s during setup: %v", tx.ID, err)
		}
	}

	// Test GetPendingTransactions(0) for all transactions
	allTxs := tp.GetPendingTransactions(0)
	if len(allTxs) != len(txs) {
		t.Errorf("Expected GetPendingTransactions(0) to return all %d transactions, but got %d", len(txs), len(allTxs))
	}

	// Test GetPendingTransactions(-1) still returns empty slice
    if len(tp.GetPendingTransactions(-1)) != 0 {t.Error("GetPendingTransactions(-1) should return 0 transactions even when pool is not empty")}
	
    got2 := tp.GetPendingTransactions(2)
	if len(got2) != 2 {
		t.Errorf("Expected 2 transactions from GetPendingTransactions(2), got %d", len(got2))
	}

	got5 := tp.GetPendingTransactions(5) // Attempt to get more than available
	if len(got5) != 3 {
		t.Errorf("Expected 3 transactions from GetPendingTransactions(5) when 3 available, got %d", len(got5))
	}
    // Check if all original tx IDs are present in got5
    retrievedIDs := make(map[string]bool)
    for _, rTx := range got5 {
        retrievedIDs[rTx.ID] = true
    }
    for _, oTx := range txs {
        if !retrievedIDs[oTx.ID] {
            t.Errorf("Original tx %s not found in retrieved transactions", oTx.ID)
        }
    }
}

func TestTxPool_RemoveTransactionByID(t *testing.T) { // Renamed from TestTxPool_RemoveTransactions
	tp := NewTxPool()
	tx1 := block.NewTransaction("s1", "r1", 10)
	tx2 := block.NewTransaction("s2", "r2", 20)
	tx3 := block.NewTransaction("s3", "r3", 30)

	if err := tp.AddTransaction(*tx1); err != nil {
		t.Fatalf("Error adding tx1 during setup: %v", err)
	}
	if err := tp.AddTransaction(*tx2); err != nil {
		t.Fatalf("Error adding tx2 during setup: %v", err)
	}
	if err := tp.AddTransaction(*tx3); err != nil {
		t.Fatalf("Error adding tx3 during setup: %v", err)
	}

	// Remove tx1
	err := tp.RemoveTransactionByID(tx1.ID)
	if err != nil {
		t.Errorf("Failed to remove tx1 by ID: %v", err)
	}
	if _, exists := tp.pending[tx1.ID]; exists {
		t.Error("tx1 should have been removed from the pool")
	}

	// Remove tx3
	err = tp.RemoveTransactionByID(tx3.ID)
	if err != nil {
		t.Errorf("Failed to remove tx3 by ID: %v", err)
	}
	if _, exists := tp.pending[tx3.ID]; exists {
		t.Error("tx3 should have been removed from the pool")
	}
    
    // Attempt to remove a non-existent ID
    err = tp.RemoveTransactionByID("nonexistentID")
    if err == nil {
        t.Error("Expected error when trying to remove non-existent transaction ID, got nil")
    }

	// Check remaining transactions
	if len(tp.pending) != 1 {
		t.Errorf("Expected 1 pending transaction after removal, got %d", len(tp.pending))
	}
	if _, exists := tp.pending[tx2.ID]; !exists {
		t.Error("tx2 should still be in the pool")
	}
}
