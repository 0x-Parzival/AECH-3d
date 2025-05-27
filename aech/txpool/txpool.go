// Package txpool manages a pool of pending transactions.
package txpool

import (
	"aech/block"
	"fmt"
	"sync"
)

// TxPool holds pending transactions that are waiting to be included in a block.
// It provides thread-safe methods to add, retrieve, and remove transactions.
type TxPool struct {
	// pending is a map of transaction IDs to Transaction objects.
	// It stores all transactions currently in the pool.
	pending map[string]block.Transaction
	// mu is a read-write mutex to ensure concurrent access to the pending map is safe.
	mu sync.RWMutex
}

// NewTxPool creates and returns a new, empty TxPool instance.
func NewTxPool() *TxPool {
	return &TxPool{
		pending: make(map[string]block.Transaction),
		// mu is initialized as a zero value RWMutex, which is ready to use.
	}
}

// AddTransaction adds a transaction to the transaction pool.
// Before adding, it validates the transaction using its Validate method.
// It returns an error if the transaction is invalid or if it already exists in the pool.
func (tp *TxPool) AddTransaction(tx block.Transaction) error {
	if !tx.Validate() { // tx.Validate() currently always returns true (stubbed)
		return fmt.Errorf("transaction %s failed validation", tx.ID)
	}

	tp.mu.Lock()
	defer tp.mu.Unlock()

	if _, exists := tp.pending[tx.ID]; exists {
		return fmt.Errorf("transaction %s already exists in the pool", tx.ID)
	}
	tp.pending[tx.ID] = tx
	return nil
}

// GetPendingTransactions retrieves a specified number of transactions from the pool.
// It takes an integer `count` specifying the maximum number of transactions to retrieve.
// It takes an integer `count` specifying the maximum number of transactions to retrieve.
// If `count == 0`, all pending transactions are retrieved.
// If `count < 0`, an empty slice is returned.
// If `count > 0` and `count > len(tp.pending)`, all pending transactions are retrieved.
// Otherwise, `count` transactions are retrieved.
// If the pool is empty, an empty slice is always returned regardless of `count` (unless `count < 0`).
// The order of transactions returned is non-deterministic due to the nature of map iteration.
func (tp *TxPool) GetPendingTransactions(count int) []block.Transaction {
	tp.mu.RLock()
	defer tp.mu.RUnlock()

	if len(tp.pending) == 0 { // If pool is empty, always return empty
		return []block.Transaction{}
	}

	if count < 0 { // If count is negative, return empty
		return []block.Transaction{}
	}

	var limit int
	if count == 0 { // count == 0 means get all
		limit = len(tp.pending)
	} else { // count > 0
		if count > len(tp.pending) {
			limit = len(tp.pending) // Cannot get more than available
		} else {
			limit = count // Get requested count
		}
	}
	
	// If limit is 0 at this point (e.g. count was 0 and pool was empty, though pool empty is handled above)
	// the make call will be make([]block.Transaction, 0, 0) and the loop won't run, returning an empty slice.
	transactions := make([]block.Transaction, 0, limit)
	i := 0
	for _, tx := range tp.pending { // Iteration order over map is not guaranteed
		if i >= limit {
			break
		}
		transactions = append(transactions, tx)
		i++
	}
	return transactions
}

// RemoveTransactions removes a slice of transactions from the pool.
// Transactions are identified by their IDs. If a transaction in the provided slice
// is not found in the pool, it is silently ignored.
// func (tp *TxPool) RemoveTransactions(txsToRemove []block.Transaction) {
// 	tp.mu.Lock()
// 	defer tp.mu.Unlock()
//
// 	for _, tx := range txsToRemove {
// 		delete(tp.pending, tx.ID)
// 	}
// }

// RemoveTransactionByID removes a single transaction from the pool by its ID.
// It returns an error if the transaction ID is not found in the pool.
func (tp *TxPool) RemoveTransactionByID(txID string) error {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	if _, exists := tp.pending[txID]; !exists {
		return fmt.Errorf("transaction %s not found in pool", txID)
	}
	delete(tp.pending, txID)
	return nil
}
