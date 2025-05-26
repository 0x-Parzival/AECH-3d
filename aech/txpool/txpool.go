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
// If `count` is zero or negative, or if the pool is empty, an empty slice is returned.
// If `count` is greater than the number of transactions in the pool, all pending transactions are returned.
// The order of transactions returned is non-deterministic due to the nature of map iteration.
func (tp *TxPool) GetPendingTransactions(count int) []block.Transaction {
	tp.mu.RLock()
	defer tp.mu.RUnlock()

	if count <= 0 || len(tp.pending) == 0 {
		return []block.Transaction{}
	}

	// Ensure we don't try to retrieve more transactions than available
	if count > len(tp.pending) {
		count = len(tp.pending)
	}

	transactions := make([]block.Transaction, 0, count)
	i := 0
	for _, tx := range tp.pending { // Iteration order over map is not guaranteed
		if i >= count {
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
