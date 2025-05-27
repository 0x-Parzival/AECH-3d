package txpool

import (
	"aech/block"
	"testing"
	"fmt"
)

// TestNewTxPool tests the creation of a new transaction pool.
func TestNewTxPool(t *testing.T) {
	t.Helper()
	pool := NewTxPool()
	if pool == nil {
		t.Fatal("NewTxPool() returned nil")
	}
	if pool.pending == nil {
		t.Error("NewTxPool() did not initialize pending transactions map")
	}
	if len(pool.pending) != 0 {
		t.Errorf("NewTxPool() pending transactions map should be empty, got %d", len(pool.pending))
	}
}

// TestAddTransaction tests adding transactions to the pool, including valid and duplicate additions.
// It also implicitly tests parts of GetPendingTransactions by checking pool state.
func TestAddTransaction(t *testing.T) {
	t.Helper()
	pool := NewTxPool()
	tx1 := block.NewTransaction("sender1", "receiver1", 100)
	tx2 := block.NewTransaction("sender2", "receiver2", 200)

	t.Run("AddValidTransaction", func(t *testing.T) {
		t.Helper()
		// Test adding a valid transaction tx1
		if err := pool.AddTransaction(*tx1); err != nil {
			t.Fatalf("AddTransaction(tx1) failed: %v", err)
		}
		if len(pool.pending) != 1 {
			t.Errorf("Expected 1 transaction in pool after adding tx1, got %d", len(pool.pending))
		}
		if _, exists := pool.pending[tx1.ID]; !exists {
			t.Errorf("tx1 was not found in the pool after adding")
		}

		// Test adding another valid transaction tx2
		if err := pool.AddTransaction(*tx2); err != nil {
			t.Fatalf("AddTransaction(tx2) failed: %v", err)
		}
		if len(pool.pending) != 2 {
			t.Errorf("Expected 2 transactions in pool after adding tx2, got %d", len(pool.pending))
		}
		if _, exists := pool.pending[tx2.ID]; !exists {
			t.Errorf("tx2 was not found in the pool after adding")
		}
	})

	t.Run("AddDuplicateTransaction", func(t *testing.T) {
		t.Helper()
		// Pool should already have tx1 and tx2 from the previous subtest.
		initialCount := len(pool.pending)
		if initialCount != 2 {
			t.Fatalf("Expected pool count to be 2 before adding duplicate, got %d", initialCount)
		}

		if err := pool.AddTransaction(*tx1); err == nil {
			t.Error("Expected error when adding duplicate transaction tx1, got nil")
		}
		if len(pool.pending) != initialCount { // Should still be 'initialCount'
			t.Errorf("Expected %d transactions in pool after duplicate add attempt, got %d", initialCount, len(pool.pending))
		}
	})

	// The original TestAddAndGetPendingTransactions also tested GetPendingTransactions.
	// This functionality will be covered more directly by TestGetAllTransactions or
	// can be a separate TestGetPendingTransactions if detailed variations are needed.
	// For now, AddTransaction focuses on the add operation and basic pool state.
}

// TestRemoveTransactionByID tests removing transactions by ID, including existing and non-existent cases.
func TestRemoveTransactionByID(t *testing.T) {
	t.Helper()
	pool := NewTxPool()
	tx1 := block.NewTransaction("sender1", "receiver1", 100)
	tx2 := block.NewTransaction("sender2", "receiver2", 200)

	// Add transactions to setup the test
	_ = pool.AddTransaction(*tx1)
	_ = pool.AddTransaction(*tx2)
	if len(pool.pending) != 2 {
		t.Fatal("Setup failed: Expected 2 transactions in pool")
	}

	t.Run("RemoveExistingTransaction", func(t *testing.T) {
		t.Helper()
		err := pool.RemoveTransactionByID(tx1.ID)
		if err != nil {
			t.Fatalf("RemoveTransactionByID(tx1.ID) failed: %v", err)
		}
		if len(pool.pending) != 1 {
			t.Errorf("Expected 1 transaction after removing tx1, got %d", len(pool.pending))
		}
		if _, exists := pool.pending[tx1.ID]; exists {
			t.Error("tx1 should not exist in pool after removal")
		}
	})

	t.Run("RemoveNonExistentTransaction", func(t *testing.T) {
		t.Helper()
		nonExistentID := "fakeID"
		initialCount := len(pool.pending) // Should be 1 at this point

		err := pool.RemoveTransactionByID(nonExistentID)
		if err == nil {
			t.Error("Expected error when removing non-existent transaction, got nil")
		}
		expectedErrorMsg := fmt.Sprintf("transaction %s not found in pool", nonExistentID)
		if err != nil && err.Error() != expectedErrorMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, err.Error())
		}

		if len(pool.pending) != initialCount {
			t.Errorf("Expected %d transaction after attempting to remove non-existent ID, got %d", initialCount, len(pool.pending))
		}
	})

	t.Run("RemoveLastTransaction", func(t *testing.T) {
		t.Helper()
		// tx2 should be the only one left
		if err := pool.RemoveTransactionByID(tx2.ID); err != nil {
			t.Fatalf("RemoveTransactionByID(tx2.ID) failed: %v", err)
		}
		if len(pool.pending) != 0 {
			t.Errorf("Expected 0 transactions after removing tx2, got %d", len(pool.pending))
		}
		if _, exists := pool.pending[tx2.ID]; exists {
			t.Error("tx2 should not exist in pool after removal")
		}
	})
}

// TestGetAllTransactions tests retrieving all transactions from the pool.
func TestGetAllTransactions(t *testing.T) {
	t.Helper()
	pool := NewTxPool()

	t.Run("EmptyPool", func(t *testing.T) {
		t.Helper()
		allTxs := pool.GetAllTransactions()
		if allTxs == nil {
			t.Error("GetAllTransactions on empty pool returned nil, expected empty slice")
		}
		if len(allTxs) != 0 {
			t.Errorf("Expected 0 transactions from GetAllTransactions on empty pool, got %d", len(allTxs))
		}
	})

	tx1 := block.NewTransaction("sender1", "receiver1", 100)
	tx2 := block.NewTransaction("sender2", "receiver2", 200)
	tx3 := block.NewTransaction("sender3", "receiver3", 300)

	_ = pool.AddTransaction(*tx1)
	_ = pool.AddTransaction(*tx2)
	_ = pool.AddTransaction(*tx3)

	t.Run("PopulatedPool", func(t *testing.T) {
		t.Helper()
		allTxs := pool.GetAllTransactions()
		if len(allTxs) != 3 {
			t.Errorf("Expected 3 transactions from GetAllTransactions, got %d", len(allTxs))
		}

		// Verify transactions (order is not guaranteed from map iteration)
		foundTx1, foundTx2, foundTx3 := false, false, false
		for _, tx := range allTxs {
			switch tx.ID {
			case tx1.ID:
				foundTx1 = true
			case tx2.ID:
				foundTx2 = true
			case tx3.ID:
				foundTx3 = true
			}
		}
		if !foundTx1 || !foundTx2 || !foundTx3 {
			t.Errorf("Not all added transactions were found in GetAllTransactions result. Found: tx1=%t, tx2=%t, tx3=%t", foundTx1, foundTx2, foundTx3)
		}
	})

	t.Run("PoolAfterRemoval", func(t *testing.T) {
		t.Helper()
		// Remove tx2
		_ = pool.RemoveTransactionByID(tx2.ID)
		allTxs := pool.GetAllTransactions()

		if len(allTxs) != 2 {
			t.Errorf("Expected 2 transactions after removing one, got %d", len(allTxs))
		}

		foundTx1, foundTx2, foundTx3 := false, false, false
		for _, tx := range allTxs {
			switch tx.ID {
			case tx1.ID:
				foundTx1 = true
			case tx2.ID: // Should not be found
				foundTx2 = true
			case tx3.ID:
				foundTx3 = true
			}
		}
		if !foundTx1 || foundTx2 || !foundTx3 {
			t.Errorf("GetAllTransactions result after removal is incorrect. Found: tx1=%t, tx2=%t (should be false), tx3=%t", foundTx1, foundTx2, foundTx3)
		}
	})
}

// TestGetPendingTransactions is also important as it's used by other parts of the system (e.g., block creation).
// The original TestAddAndGetPendingTransactions covered this. Let's ensure it's still covered adequately.
func TestGetPendingTransactions(t *testing.T) {
	t.Helper()
	pool := NewTxPool()
	tx1 := block.NewTransaction("sender1", "receiver1", 100)
	tx2 := block.NewTransaction("sender2", "receiver2", 200)

	_ = pool.AddTransaction(*tx1)
	_ = pool.AddTransaction(*tx2)


	t.Run("CountLessThanTotal", func(t *testing.T) {
		t.Helper()
		retrievedTxs := pool.GetPendingTransactions(1)
		if len(retrievedTxs) != 1 {
			t.Errorf("Expected 1 transaction from GetPendingTransactions(1), got %d", len(retrievedTxs))
		}
	})

	t.Run("CountEqualToTotal", func(t *testing.T) {
		t.Helper()
		retrievedTxs := pool.GetPendingTransactions(2)
		if len(retrievedTxs) != 2 {
			t.Errorf("Expected 2 transactions from GetPendingTransactions(2), got %d", len(retrievedTxs))
		}
	})
	
	t.Run("CountMoreThanTotal", func(t *testing.T) {
		t.Helper()
		retrievedTxs := pool.GetPendingTransactions(3)
		if len(retrievedTxs) != 2 { // Should return all available
			t.Errorf("Expected 2 transactions from GetPendingTransactions(3), got %d", len(retrievedTxs))
		}
	})

	t.Run("CountIsZero", func(t *testing.T) {
		t.Helper()
		retrievedTxs := pool.GetPendingTransactions(0)
		if len(retrievedTxs) != 0 {
			t.Errorf("Expected 0 transactions from GetPendingTransactions(0), got %d", len(retrievedTxs))
		}
	})

	t.Run("PoolIsEmpty", func(t *testing.T) {
		t.Helper()
		emptyPool := NewTxPool()
		retrievedTxs := emptyPool.GetPendingTransactions(1)
		if len(retrievedTxs) != 0 {
			t.Errorf("Expected 0 transactions from GetPendingTransactions(1) on empty pool, got %d", len(retrievedTxs))
		}
	})
}

// Note on "invalid transaction": block.Transaction.Validate() currently always returns true.
// If Validate() logic becomes meaningful, TestAddTransaction should include cases for invalid transactions
// where AddTransaction would be expected to return an error.
// e.g., if tx.Validate() returns false, pool.AddTransaction(tx) should return an error.
// This would require modifying AddTransaction to call Validate().
// func (tp *TxPool) AddTransaction(tx block.Transaction) error {
// 	if !tx.Validate() { // Assuming Validate exists and is meaningful
// 		return fmt.Errorf("transaction %s is invalid", tx.ID)
// 	}
// ... rest of the logic
// }
```
