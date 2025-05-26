package main

import (
	"aech/block"
	"aech/plane"
	"log"
)

func main() {
	// Initialize the new Plane3D system
	planeAPI := plane.NewPlane()
	if planeAPI == nil {
		log.Fatal("Failed to create Plane3D instance") // Should not happen with current NewPlane
	}

	log.Println("Plane3D system initialized.")

	genesisTransactions := []block.Transaction{}
	genesisBlock := block.NewBlock(0, 0, 0, "0", genesisTransactions)

	if genesisBlock == nil {
		log.Fatal("Failed to create Genesis block")
	}

	// Use IsValidPosition (optional, but good practice)
	if !planeAPI.IsValidPosition(genesisBlock.X, genesisBlock.Y, genesisBlock.Z) {
		log.Fatalf("Genesis block coordinates (%d,%d,%d) are invalid.", genesisBlock.X, genesisBlock.Y, genesisBlock.Z)
	}

	// Insert the Genesis block using the AddBlock method of the Plane3D instance
	err := planeAPI.AddBlock(genesisBlock.X, genesisBlock.Y, genesisBlock.Z, genesisBlock)
	if err != nil {
		// This might legitimately happen if main() were run twice without resetting state,
		// but for a single run, it implies an issue or unexpected state.
		log.Fatalf("Failed to add Genesis block to plane: %v", err)
	}

	log.Printf("Genesis block (Hash: %s) added to plane at coordinates (%d, %d, %d)", genesisBlock.Hash, genesisBlock.X, genesisBlock.Y, genesisBlock.Z)
	log.Printf("Genesis block Timestamp: %s", genesisBlock.Timestamp.String())

	// Optional: Verify retrieval using the new GetBlock method
	// retrievedGenesis, err := planeAPI.GetBlock(0, 0, 0)
	// if err != nil {
	// 	log.Printf("Failed to retrieve Genesis block after insertion: %v", err)
	// } else {
	// 	log.Printf("Successfully retrieved Genesis block. Hash: %s", retrievedGenesis.Hash)
	// 	if retrievedGenesis.Hash != genesisBlock.Hash {
	// 		log.Printf("Error: Retrieved Genesis block hash %s does not match original %s", retrievedGenesis.Hash, genesisBlock.Hash)
	// 	}
	// }
}
