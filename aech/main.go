package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"aech/block"
	"aech/plane"
	// "aech/txpool" // Not strictly needed for stateless CLI as per plan
	// "aech/consensus" // Not strictly needed for stateless CLI
)

func main() {
	// Define flags for add-tx command
	addTxCmd := flag.Bool("add-tx", false, "Add a new transaction. Requires --sender, --receiver, --amount.")
	senderFlag := flag.String("sender", "", "Sender's address for the transaction.")
	receiverFlag := flag.String("receiver", "", "Receiver's address for the transaction.")
	amountFlag := flag.Float64("amount", 0.0, "Amount to transact.")

	// Define flags for mine-block command
	mineBlockCmd := flag.Bool("mine-block", false, "Mine a new block. Requires --x, --y, --z, --prevhash.")
	xFlag := flag.Int("x", 0, "X coordinate for the new block.")
	yFlag := flag.Int("y", 0, "Y coordinate for the new block.")
	zFlag := flag.Int("z", 0, "Z coordinate for the new block.")
	prevHashFlag := flag.String("prevhash", "", "Previous block's hash.")

	// Define flags for print-plane command
	printPlaneCmd := flag.Bool("print-plane", false, "Print a simulated blockchain plane with a few blocks.")

	flag.Parse()

	// Check which command is active
	if *addTxCmd {
		if *senderFlag == "" || *receiverFlag == "" || *amountFlag <= 0 {
			fmt.Println("Error: --sender, --receiver, and a positive --amount are required for --add-tx.")
			flag.Usage()
			os.Exit(1)
		}
		fmt.Println("Executing 'add-tx' command...")
		tx := block.NewTransaction(*senderFlag, *receiverFlag, *amountFlag)
		if tx == nil {
			log.Fatal("Failed to create transaction.") // Should not happen with current NewTransaction
		}
		fmt.Printf("Transaction Created:\n")
		fmt.Printf("  ID: %s\n", tx.ID)
		fmt.Printf("  Sender: %s\n", tx.Sender)
		fmt.Printf("  Receiver: %s\n", tx.Receiver)
		fmt.Printf("  Amount: %.2f\n", tx.Amount)
		fmt.Printf("  Timestamp: %d\n", tx.Timestamp)
		fmt.Printf("  Signature: %x\n", tx.Signature)

	} else if *mineBlockCmd {
		if *prevHashFlag == "" { // Basic check, coordinates default to 0,0,0 if not provided
			fmt.Println("Error: --prevhash is required for --mine-block.")
			flag.Usage()
			os.Exit(1)
		}
		fmt.Println("Executing 'mine-block' command...")
		// Create dummy transactions for the block
		dummyTx1 := block.NewTransaction("miner_reward_sender", "miner_address", 10.0)
		dummyTx2 := block.NewTransaction("some_user_A", "some_user_B", 5.0)
		if dummyTx1 == nil || dummyTx2 == nil {
			log.Fatal("Failed to create dummy transactions for block.")
		}
		transactions := []block.Transaction{*dummyTx1, *dummyTx2}

		newBlk := block.NewBlock(*xFlag, *yFlag, *zFlag, *prevHashFlag, transactions)
		if newBlk == nil {
			log.Fatal("Failed to create new block.")
		}
		fmt.Printf("Block Mined:\n")
		fmt.Printf("  Hash: %s\n", newBlk.Hash)
		fmt.Printf("  Coordinates: (%d, %d, %d)\n", newBlk.X, newBlk.Y, newBlk.Z)
		fmt.Printf("  Previous Hash: %s\n", newBlk.PreviousHash)
		fmt.Printf("  Timestamp: %s\n", newBlk.Timestamp.String())
		fmt.Printf("  Transaction Count: %d\n", len(newBlk.Transactions))
		for i, tx := range newBlk.Transactions {
			fmt.Printf("    Tx %d ID: %s (Sender: %s, Receiver: %s, Amount: %.2f)\n", i+1, tx.ID, tx.Sender, tx.Receiver, tx.Amount)
		}

	} else if *printPlaneCmd {
		fmt.Println("Executing 'print-plane' command...")
		p := plane.NewPlane()

		// Add Genesis block
		genesisTxs := []block.Transaction{} // Genesis might have no app-level txs or specific coinbase
		genesisBlock := block.NewBlock(0, 0, 0, "0000000000000000000000000000000000000000000000000000000000000000", genesisTxs) // Using full length zero hash for clarity
		if genesisBlock == nil {
			log.Fatal("Failed to create Genesis block for print-plane.")
		}
		err := p.AddBlock(genesisBlock.X, genesisBlock.Y, genesisBlock.Z, genesisBlock)
		if err != nil {
			log.Fatalf("Failed to add Genesis block to plane for simulation: %v", err)
		}

		// Add a second block
		block2Txs := []block.Transaction{*block.NewTransaction("user1", "user2", 25.0)}
		block2 := block.NewBlock(1, 0, 0, genesisBlock.Hash, block2Txs)
		if block2 == nil {
			log.Fatal("Failed to create second block for print-plane.")
		}
		err = p.AddBlock(block2.X, block2.Y, block2.Z, block2)
		if err != nil {
			log.Fatalf("Failed to add second block to plane for simulation: %v", err)
		}
		
		// Add a third block, perhaps adjacent to genesis in another direction
		block3Txs := []block.Transaction{*block.NewTransaction("user3", "user4", 15.0)}
		block3 := block.NewBlock(0, 1, 0, genesisBlock.Hash, block3Txs)
		if block3 == nil {
			log.Fatal("Failed to create third block for print-plane.")
		}
		err = p.AddBlock(block3.X, block3.Y, block3.Z, block3)
		if err != nil {
			log.Fatalf("Failed to add third block to plane for simulation: %v", err)
		}


		fmt.Println("\nSimulated Plane Structure:")
		grid := p.Grid // Assuming Grid is public, as per previous use in test.go
		if len(grid) == 0 {
			fmt.Println("Plane is empty.")
		} else {
			// For consistent output, it might be better to iterate over known keys or sort keys if possible
			// For now, direct iteration is fine for this stateless demo.
			for key, blk := range grid {
				if blk != nil {
					fmt.Printf("Block at key '%s' (Coords: %d,%d,%d): Hash %s, PrevHash: %s, Txs: %d\n",
						key, blk.X, blk.Y, blk.Z, blk.Hash, blk.PreviousHash, len(blk.Transactions))
				}
			}
		}
		fmt.Println("\nNote: This is a stateless, simulated plane for demonstration.")

	} else {
		fmt.Println("No command specified or invalid command combination.")
		fmt.Println("Usage: go run aech/main.go <command_flag> [options]")
		fmt.Println("\nAvailable commands and their specific flags:")
		fmt.Println("  --add-tx: Add a new transaction")
		fmt.Println("    --sender <string>: Sender's address")
		fmt.Println("    --receiver <string>: Receiver's address")
		fmt.Println("    --amount <float64>: Amount to transact (must be > 0)")
		fmt.Println("\n  --mine-block: Mine a new block")
		fmt.Println("    --x <int>: X coordinate for the block (defaults to 0)")
		fmt.Println("    --y <int>: Y coordinate for the block (defaults to 0)")
		fmt.Println("    --z <int>: Z coordinate for the block (defaults to 0)")
		fmt.Println("    --prevhash <string>: Previous block's hash (required)")
		fmt.Println("\n  --print-plane: Print a simulated blockchain plane")
		fmt.Println("\nExample for add-tx: go run aech/main.go --add-tx --sender Alice --receiver Bob --amount 10.5")
		fmt.Println("Example for mine-block: go run aech/main.go --mine-block --x 1 --y 0 --z 0 --prevhash <hash_of_block_at_0,0,0>")
		fmt.Println("Example for print-plane: go run aech/main.go --print-plane")
		os.Exit(1)
	}
}
