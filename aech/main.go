package main

import (
	"aech/block"
	"aech/plane"
	"log"
)

func main() {
	// Genesis block typically has no transactions or a special coinbase transaction.
	// For simplicity, we'll use an empty transaction slice.
	genesisTransactions := []block.Transaction{}

	// Create the Genesis block.
	// The PreviousHash for a Genesis block is conventionally "0" or an empty string.
	// Our NewBlock function handles timestamping and hash generation.
	genesisBlock := block.NewBlock(0, 0, 0, "0", genesisTransactions)

	if genesisBlock == nil {
		log.Fatal("Failed to create Genesis block")
	}

	// Insert the Genesis block into the plane.
	err := plane.InsertBlock(genesisBlock.X, genesisBlock.Y, genesisBlock.Z, *genesisBlock)
	if err != nil {
		log.Fatalf("Failed to insert Genesis block into plane: %v", err)
	}

	log.Printf("Genesis block created and inserted at coordinates (%d, %d, %d)", genesisBlock.X, genesisBlock.Y, genesisBlock.Z)
	log.Printf("Genesis block Hash: %s", genesisBlock.Hash)
	log.Printf("Genesis block Timestamp: %s", genesisBlock.Timestamp.String())

	// Optional: Verify retrieval (can be commented out for normal run)
	// retrievedGenesis, found := plane.GetBlock(0, 0, 0)
	// if !found {
	// 	log.Println("Failed to retrieve Genesis block after insertion.")
	// } else {
	// 	log.Printf("Successfully retrieved Genesis block. Hash: %s", retrievedGenesis.Hash)
	// }
}
