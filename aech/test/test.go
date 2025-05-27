package main

import (
	"fmt"
	"log"

	"aech/block"
	"aech/consensus"
	"aech/plane"
	"aech/txpool"
)

func main() {
	fmt.Println("Starting AECH Simulation...")

	// 1. Setup
	p := plane.NewPlane() // Corrected: NewPlane() as per codebase
	tp := txpool.NewTxPool()
	cs := consensus.NewConsensus(p)

	fmt.Println("Setup Complete: Plane, TxPool, Consensus initialized.")
	fmt.Println("---")

	// 2. Create and Add Initial Transactions
	tx1 := block.NewTransaction("Alice", "Bob", 10)
	tx2 := block.NewTransaction("Bob", "Charlie", 5)
	tx3 := block.NewTransaction("Charlie", "Alice", 7)
	tx4 := block.NewTransaction("Dave", "Eve", 12)
	tx5 := block.NewTransaction("Eve", "Frank", 3)

	initialTxs := []*block.Transaction{tx1, tx2, tx3, tx4, tx5}
	for i, tx := range initialTxs {
		err := tp.AddTransaction(*tx) // Dereference tx
		if err != nil {
			log.Fatalf("Failed to add initial transaction %d (%s): %v", i+1, tx.ID, err)
		}
	}
	fmt.Printf("%d Initial Transactions added to TxPool.\n", len(initialTxs))
	// Using a large number to effectively get all pending transactions for display purposes
	fmt.Println("Current TxPool size:", len(tp.GetPendingTransactions(100))) 
	fmt.Println("---")

	// 3. Create Genesis Block (0,0,0)
	pendingTxsForGenesis := tp.GetPendingTransactions(2)
	if len(pendingTxsForGenesis) < 2 { // As per instruction, expecting 2 for genesis
		log.Fatalf("Not enough transactions in TxPool for Genesis block. Found: %d, Expected: 2", len(pendingTxsForGenesis))
	}
	genesisBlock := block.NewBlock(0, 0, 0, "GENESIS", pendingTxsForGenesis)

	canAddGenesis := cs.CanAddBlock(0, 0, 0, genesisBlock)
	if !canAddGenesis {
		log.Fatalf("Consensus: Cannot add Genesis Block at (0,0,0)")
	}

	err := p.AddBlock(0, 0, 0, genesisBlock)
	if err != nil {
		log.Fatalf("Failed to add Genesis Block to plane: %v", err)
	}

	for _, tx := range pendingTxsForGenesis {
		err := tp.RemoveTransactionByID(tx.ID)
		if err != nil {
			// This should ideally not happen if tx was correctly fetched from pool
			log.Printf("Warning: Failed to remove transaction %s from TxPool after Genesis block: %v", tx.ID, err)
		}
	}

	fmt.Println("Genesis Block added at (0,0,0)")
	fmt.Println("Hash:", genesisBlock.Hash)
	fmt.Println("Transactions included:", len(genesisBlock.Transactions))
	fmt.Println("Remaining Pending Transactions in Pool:", len(tp.GetPendingTransactions(100)))
	fmt.Println("---")

	// 4. Create Second Block (e.g., 1,0,0)
	tx6 := block.NewTransaction("Gina", "Harry", 8)
	tx7 := block.NewTransaction("Harry", "Ivy", 4)
	
	newTxsForBlock2 := []*block.Transaction{tx6, tx7}
	for i, tx := range newTxsForBlock2 {
		err := tp.AddTransaction(*tx) // Dereference tx
		if err != nil {
			log.Fatalf("Failed to add transaction %d (%s) for Block 2: %v", i+1, tx.ID, err)
		}
	}
	fmt.Printf("%d New Transactions added to TxPool for Block 2.\n", len(newTxsForBlock2))
	
	block2Txs := tp.GetPendingTransactions(2) // Attempt to get 2 txs
	if len(block2Txs) == 0 && len(tp.GetPendingTransactions(100)) > 0 { // Check if pool is not empty but we got 0
		log.Fatalf("TxPool is not empty, but no transactions were fetched for Block 2.")
	}
	fmt.Printf("Fetched %d transactions for Block 2.\n", len(block2Txs))


	block2 := block.NewBlock(1, 0, 0, genesisBlock.Hash, block2Txs)
	canAddBlock2 := cs.CanAddBlock(1, 0, 0, block2)
	if !canAddBlock2 {
		// For debugging consensus issues if they arise:
		// neighbors := p.ListNeighbors(1,0,0)
		// fmt.Printf("Neighbors of (1,0,0): %+v\n", neighbors)
		// existingBlock, _ := p.GetBlock(1,0,0)
		// fmt.Printf("Existing block at (1,0,0): %+v\n", existingBlock)
		log.Fatalf("Consensus: Cannot add Block 2 at (1,0,0).")
	}

	err = p.AddBlock(1, 0, 0, block2)
	if err != nil {
		log.Fatalf("Failed to add Block 2 to plane: %v", err)
	}

	for _, tx := range block2Txs {
		err := tp.RemoveTransactionByID(tx.ID)
		if err != nil {
			log.Printf("Warning: Failed to remove transaction %s from TxPool after Block 2: %v", tx.ID, err)
		}
	}

	fmt.Printf("Block added at (%d,%d,%d)\n", block2.X, block2.Y, block2.Z)
	fmt.Println("Hash:", block2.Hash)
	fmt.Println("Previous Hash:", block2.PreviousHash)
	fmt.Println("Transactions included:", len(block2.Transactions))
	fmt.Println("Remaining Pending Transactions in Pool:", len(tp.GetPendingTransactions(100)))
	fmt.Println("---")

	// 5. Create Third Block (e.g., 0,1,0)
	tx8 := block.NewTransaction("Jack", "Kim", 6)
	tx9 := block.NewTransaction("Kim", "Leo", 9)

	newTxsForBlock3 := []*block.Transaction{tx8, tx9}
	for i, tx := range newTxsForBlock3 {
		err := tp.AddTransaction(*tx) // Dereference tx
		if err != nil {
			log.Fatalf("Failed to add transaction %d (%s) for Block 3: %v", i+1, tx.ID, err)
		}
	}
	fmt.Printf("%d New Transactions added to TxPool for Block 3.\n", len(newTxsForBlock3))

	block3Txs := tp.GetPendingTransactions(2) // Attempt to get 2 txs
	if len(block3Txs) == 0 && len(tp.GetPendingTransactions(100)) > 0 {
		log.Fatalf("TxPool is not empty, but no transactions were fetched for Block 3.")
	}
	fmt.Printf("Fetched %d transactions for Block 3.\n", len(block3Txs))
	
	// Attaching to Genesis block at (0,0,0) making it adjacent at (0,1,0)
	block3 := block.NewBlock(0, 1, 0, genesisBlock.Hash, block3Txs)
	canAddBlock3 := cs.CanAddBlock(0, 1, 0, block3)
	if !canAddBlock3 {
		// For debugging consensus issues if they arise:
		// neighbors := p.ListNeighbors(0,1,0)
		// fmt.Printf("Neighbors of (0,1,0): %+v\n", neighbors)
		// existingBlock, _ := p.GetBlock(0,1,0)
		// fmt.Printf("Existing block at (0,1,0): %+v\n", existingBlock)
		log.Fatalf("Consensus: Cannot add Block 3 at (0,1,0).")
	}

	err = p.AddBlock(0, 1, 0, block3)
	if err != nil {
		log.Fatalf("Failed to add Block 3 to plane: %v", err)
	}

	for _, tx := range block3Txs {
		err := tp.RemoveTransactionByID(tx.ID)
		if err != nil {
			log.Printf("Warning: Failed to remove transaction %s from TxPool after Block 3: %v", tx.ID, err)
		}
	}
	
	fmt.Printf("Block added at (%d,%d,%d)\n", block3.X, block3.Y, block3.Z)
	fmt.Println("Hash:", block3.Hash)
	fmt.Println("Previous Hash:", block3.PreviousHash)
	fmt.Println("Transactions included:", len(block3.Transactions))
	fmt.Println("Remaining Pending Transactions in Pool:", len(tp.GetPendingTransactions(100)))
	fmt.Println("---")

	// 6. Final Output - Print Plane Structure
	fmt.Println("\nFinal Plane Structure:")
	// Accessing p.Grid directly as it's public.
	grid := p.Grid 

	if len(grid) == 0 {
		fmt.Println("Plane is empty.")
	} else {
		for key, blk := range grid {
			if blk != nil { // Should always be non-nil if added correctly
				fmt.Printf("Block at key %s (Coords: %d,%d,%d): Hash %s, PrevHash: %s, Txs: %d\n", 
					key, blk.X, blk.Y, blk.Z, blk.Hash, blk.PreviousHash, len(blk.Transactions))
			} else {
				// This case should ideally not be reached if plane.AddBlock works as expected
				fmt.Printf("Error: Block at key %s is nil in the grid\n", key)
			}
		}
	}
	
	fmt.Println("\nAECH Simulation Finished.")
}
