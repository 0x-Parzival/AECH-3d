package main

import (
	"fmt"
	"aech/block"
	"aech/plane"
	"aech/txpool"
	"time"
	"aech/api" // May need to uncomment or add this
)

func main() {
	// Initialize the transaction pool
	pool := txpool.NewTxPool()

	// Create some dummy transactions
	tx1 := block.NewTransaction("Alice", "Bob", 10.0)
	tx2 := block.NewTransaction("Bob", "Charlie", 5.0)

	// Add transactions to the pool
	_ = pool.AddTransaction(*tx1) // Errors are logged within AddTransaction
	_ = pool.AddTransaction(*tx2) // Errors are logged within AddTransaction

	// Initialize the 3D plane
	p := plane.NewPlane()

	// Retrieve some transactions to include in a block
	transactionsForBlock1 := pool.GetPendingTransactions(2)

	// Create a "genesis" block (or first block)
	// For a true genesis block, previousHash would be "0" or an empty string.
	// Let's assume there's no previous block for simplicity in this example.
	genesisBlock := block.NewBlock(0, 0, 0, "0", transactionsForBlock1)
	// fmt.Printf("Genesis Block: %+v\n", genesisBlock) // Removed

	// Add the block to the plane
	if err := p.AddBlock(genesisBlock.X, genesisBlock.Y, genesisBlock.Z, genesisBlock); err != nil {
		// fmt.Println("Error adding genesis block to plane:", err) // Error logged in AddBlock
		return // Keep return on critical error
	}
	// fmt.Printf("Added Genesis Block to plane at (%d,%d,%d)\n", genesisBlock.X, genesisBlock.Y, genesisBlock.Z) // Removed

	// Simulate adding another block
	// Assume some time has passed and new transactions have arrived (or use the same pool for simplicity)
	time.Sleep(1 * time.Second) // Simulate time passing for a new timestamp

	// Let's create new transactions for the second block
	tx3 := block.NewTransaction("Charlie", "David", 2.0)
	_ = pool.AddTransaction(*tx3) // Error logged in AddTransaction
	transactionsForBlock2 := pool.GetPendingTransactions(1) // Get one tx

	// Create a second block, linked to the genesis block
	block2 := block.NewBlock(1, 0, 0, genesisBlock.Hash, transactionsForBlock2)
	if err := p.AddBlock(block2.X, block2.Y, block2.Z, block2); err != nil {
		// fmt.Println("Error adding block2 to plane:", err) // Error logged in AddBlock
		return // Keep return on critical error
	}
	// fmt.Printf("Added Block 2 to plane at (%d,%d,%d), linked to %s\n", block2.X, block2.Y, block2.Z, block2.PreviousHash) // Removed

	// Example of retrieving a block
	// _, err := p.GetBlock(0, 0, 0) // Result not used, error handling for GetBlock is typically in API layer or specific logic
	// if err != nil {
		// fmt.Println("Error retrieving block:", err) // Removed
	// } else {
		// fmt.Printf("Retrieved Block: %+v\n", retrievedBlock) // Removed
	// }

	// TODO: Start the API server here
	// Example: api.StartServer(p, pool)
	// For now, we're just printing to console. The actual API server setup will involve more.
	// fmt.Println("AECH system initialized (simulation complete).") // This line is removed

    // Start the API server
    // The server will block, so this should be the last call in main.
    api.RunServer("8080", p, pool)

    // Keep main running if server were started in a goroutine, or let server block.
    // select {} // if server runs in goroutine // This line is removed as RunServer blocks
}
