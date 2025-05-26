# AECH: 3D Blockchain Financial Operating System (Core Implementation)

This repository contains the initial core implementation of AECH, a conceptual 3D blockchain-based financial operating system. This early version focuses on the fundamental structures for creating blocks, storing them in a 3D spatial grid, and initializing a Genesis block.

## Current Features

*   **Block3D Structure**: Defines blocks with `X,Y,Z` coordinates, `Timestamp`, `Hash`, `PreviousHash`, and a list of `Transactions`.
*   **Plane System**: A 3D map (`Cube`) for storing and retrieving `Block3D` objects based on their spatial coordinates.
*   **Block Operations**:
    *   `block.NewBlock()`: Creates new blocks, automatically generating a timestamp and hash.
    *   `block.GenerateHash()`: Computes a SHA256 hash for a block.
    *   `block.ValidateBlock()`: A basic stub for block validation.
*   **Plane Operations**:
    *   `plane.InsertBlock()`: Inserts a block into the 3D `Cube`.
    *   `plane.GetBlock()`: Retrieves a block from specified coordinates.
    *   `plane.GetNeighbors()`: Finds adjacent blocks in the 6 cardinal directions.
*   **Genesis Block**: `main.go` initializes the blockchain with a Genesis block at coordinates (0,0,0).
*   **Unit Tests**: Basic tests for core `block` and `plane` functionalities are available in the `/test` directory.

## Directory Structure

- `/block/`: Contains the `Block3D` structure definition and related functions.
- `/plane/`: Contains the `Cube` data structure for 3D spatial storage and related functions.
- `/test/`: Contains unit tests for the project (e.g., `block_plane_test.go`).
- `main.go`: Entry point of the application; currently initializes the Genesis block.
- `go.mod`, `go.sum`: Go module files.

## Block Structure Specification (`block.Block3D`)

- `X, Y, Z (int)`: Coordinates of the block in the 3D space.
- `Timestamp (time.Time)`: Time of block creation.
- `Hash (string)`: SHA256 hash derived from X,Y,Z, Timestamp, PreviousHash, and concatenated transaction hashes.
- `PreviousHash (string)`: Hash of a conceptual preceding block. For the Genesis block, this is "0".
- `Transactions ([]Transaction)`: A list of transactions included in the block.
    - `Transaction struct { Data []byte }`: Placeholder for actual transaction data. (Note: This is outdated, will be updated in a later step if specified by another subtask, the examples below use the newer Transaction structure.)

## How to Run

1.  **Prerequisites**:
    *   Go (version 1.21 or later recommended).
2.  **Clone the repository**:
    ```bash
    # git clone <repository-url>
    # cd aech
    ```
3.  **Initialize Go modules** (if you haven't run any `go` commands yet):
    ```bash
    go mod tidy
    ```
4.  **Run the main application** (initializes Genesis block):
    ```bash
    go run main.go
    ```
5.  **Run tests**:
    ```bash
    go test ./...
    ```

## Sample Output (`go run main.go`)

The output will be similar to (timestamps will vary):

```log
2025/05/26 18:47:37 Genesis block created and inserted at coordinates (0, 0, 0)
2025/05/26 18:47:37 Genesis block Hash: 05676806c90e6b384b9bcaf6aaddbe494a959fdf87390c9e72734965be91ef84
2025/05/26 18:47:37 Genesis block Timestamp: 2025-05-26 18:47:37.855062506 +0000 UTC m=+0.000084618
```

## Usage Examples

### Creating a Transaction

Illustrates how to create a new transaction using `block.NewTransaction`. The ID, Timestamp (as UnixNano), and a dummy Signature are automatically generated.

```go
package main

import (
	"aech/block"
	"fmt"
	"time" // Used for demonstration if you want to see timestamp interpretation
)

func main() {
	tx1 := block.NewTransaction("Alice", "Bob", 10.5)
	if tx1 != nil {
		fmt.Printf("Created Transaction:\n")
		fmt.Printf("  ID: %s\n", tx1.ID)
		fmt.Printf("  Sender: %s\n", tx1.Sender)
		fmt.Printf("  Receiver: %s\n", tx1.Receiver)
		fmt.Printf("  Amount: %.2f\n", tx1.Amount)
		// Timestamp is int64 (UnixNano)
		fmt.Printf("  Timestamp (UnixNano): %d\n", tx1.Timestamp)
		fmt.Printf("  Timestamp (Human-readable): %s\n", time.Unix(0, tx1.Timestamp).Format(time.RFC3339Nano))
		fmt.Printf("  Signature: %x\n", tx1.Signature) // Print signature as hex
	}
}
```

### Using the Transaction Pool (`TxPool`)

Demonstrates initializing a `TxPool`, adding transactions, retrieving pending transactions, and removing a transaction by its ID.

```go
package main

import (
	"aech/block"
	"aech/txpool"
	"fmt"
)

func main() {
	tp := txpool.NewTxPool()
	fmt.Println("Transaction Pool initialized.")

	// Create some transactions
	txA := block.NewTransaction("Charlie", "David", 50.0)
	txB := block.NewTransaction("Eve", "Frank", 75.25)

	// Add transactions to the pool
	if txA != nil {
		err := tp.AddTransaction(*txA)
		if err != nil {
			fmt.Printf("Error adding txA (%s): %v\n", txA.ID, err)
		} else {
			fmt.Printf("Added txA (%s) to pool.\n", txA.ID)
		}
	}
	if txB != nil {
		err := tp.AddTransaction(*txB)
		if err != nil {
			fmt.Printf("Error adding txB (%s): %v\n", txB.ID, err)
		} else {
			fmt.Printf("Added txB (%s) to pool.\n", txB.ID)
		}
	}
    
	// Get pending transactions
	pendingTxs := tp.GetPendingTransactions(10) // Get up to 10
	fmt.Printf("Number of pending transactions: %d\n", len(pendingTxs))
	for i, ptx := range pendingTxs {
		fmt.Printf("  Pending Tx %d ID: %s\n", i+1, ptx.ID)
	}

	// Remove a transaction by ID
	if txA != nil {
		fmt.Printf("Attempting to remove txA (%s) by ID...\n", txA.ID)
		err := tp.RemoveTransactionByID(txA.ID)
		if err != nil {
			fmt.Printf("Error removing txA by ID: %v\n", err)
		} else {
			fmt.Printf("Successfully removed txA by ID.\n")
		}
	}
    
    // Verify removal
    pendingAfterRemove := tp.GetPendingTransactions(10)
    fmt.Printf("Number of pending transactions after removal: %d\n", len(pendingAfterRemove))
    for i, ptx := range pendingAfterRemove {
		fmt.Printf("  Remaining Tx %d ID: %s\n", i+1, ptx.ID)
	}
}
```

### Checking Block Placement with Consensus Logic

Shows a conceptual example of how `Consensus.CanAddBlock` might be used with the `Plane3D` to validate if a block can be added according to spatial rules (e.g., adjacency to existing blocks).

```go
package main

import (
	"aech/block"
	"aech/plane"
	"aech/consensus"
	"fmt"
)

func main() {
	// 1. Initialize Plane and Consensus
	p := plane.NewPlane()
	cs := consensus.NewConsensus(p) 
	fmt.Println("Plane and Consensus system initialized.")

	// 2. Create a Genesis Block
	genesisBlock := block.NewBlock(0, 0, 0, "0", nil) 

	// 3. Check if Genesis block can be added and add it
	fmt.Printf("Attempting to place Genesis Block at (%d,%d,%d)...\n", genesisBlock.X, genesisBlock.Y, genesisBlock.Z)
	if cs.CanAddBlock(genesisBlock.X, genesisBlock.Y, genesisBlock.Z, genesisBlock) {
		err := p.AddBlock(genesisBlock.X, genesisBlock.Y, genesisBlock.Z, genesisBlock)
		if err != nil {
			fmt.Printf("Error adding Genesis block to plane: %v\n", err)
		} else {
			fmt.Printf("Genesis block (Hash: %s) added to the plane.\n", genesisBlock.Hash)
		}
	} else {
		fmt.Println("Consensus rules prevent adding Genesis block where specified.")
	}

	// 4. Create a subsequent block, adjacent to Genesis
	if _, err := p.GetBlock(0,0,0); err == nil { // Check if Genesis was actually added
        nextBlock := block.NewBlock(0, 0, 1, genesisBlock.Hash, nil) 
        fmt.Printf("Attempting to place Next Block at (%d,%d,%d), adjacent to Genesis...\n", nextBlock.X, nextBlock.Y, nextBlock.Z)
        if cs.CanAddBlock(nextBlock.X, nextBlock.Y, nextBlock.Z, nextBlock) {
             err := p.AddBlock(nextBlock.X, nextBlock.Y, nextBlock.Z, nextBlock)
             if err != nil {
                fmt.Printf("Error adding next block to plane: %v\n", err)
             } else {
                fmt.Printf("Next block (Hash: %s) added to the plane at (%d,%d,%d).\n", nextBlock.Hash, nextBlock.X, nextBlock.Y, nextBlock.Z)
             }
        } else {
            fmt.Printf("Consensus rules prevent adding next block at (%d,%d,%d).\n", nextBlock.X, nextBlock.Y, nextBlock.Z)
        }

        // 5. Create an isolated block (should fail consensus)
        isolatedBlock := block.NewBlock(5,5,5, genesisBlock.Hash, nil)
        fmt.Printf("Attempting to place Isolated Block at (%d,%d,%d)...\n", isolatedBlock.X, isolatedBlock.Y, isolatedBlock.Z)
        if cs.CanAddBlock(isolatedBlock.X, isolatedBlock.Y, isolatedBlock.Z, isolatedBlock) {
            fmt.Println("Consensus rules unexpectedly allowed isolated block!")
            // Try to add it to see plane's behavior (though consensus should prevent)
            _ = p.AddBlock(isolatedBlock.X, isolatedBlock.Y, isolatedBlock.Z, isolatedBlock)
        } else {
             fmt.Printf("Consensus rules correctly prevented adding isolated block at (%d,%d,%d).\n", isolatedBlock.X, isolatedBlock.Y, isolatedBlock.Z)
        }
    } else {
        fmt.Println("Genesis block not found in plane, skipping subsequent block tests.")
    }
}
```

## Architecture Diagram

*(Placeholder: A visual diagram, possibly using ASCII art or a linked image (e.g., from draw.io), will be added here in future iterations to better illustrate the 3D cube structure and block relationships.)*
