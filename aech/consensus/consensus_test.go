package consensus

import (
	"aech/block"
	"aech/plane"
	"testing"
	"time" // Added for sleep to ensure distinct blocks if NewBlock is called rapidly
)

func setupConsensusTest(t *testing.T) (*Consensus, *plane.Plane3D) {
	p := plane.NewPlane()
	cs := NewConsensus(p)
	if cs == nil {
		t.Fatal("NewConsensus returned nil")
	}
	if cs.plane == nil {
		t.Fatal("Consensus plane not initialized")
	}
	return cs, p
}

func TestNewConsensus(t *testing.T) {
	p := plane.NewPlane()
	cs := NewConsensus(p)
	if cs == nil {
		t.Fatal("NewConsensus() returned nil")
	}
	if cs.plane != p {
		t.Error("Consensus struct does not store the provided plane instance")
	}
	if len(cs.Validators) == 0 { // Based on NewConsensus adding dummy validators
		t.Error("NewConsensus should initialize some validators (dummy)")
	}
}

func TestCanAddBlock(t *testing.T) {
	cs, p := setupConsensusTest(t)
	// Ensure distinct timestamps for block creation if logic depends on unique hashes from NewBlock
	txs := []block.Transaction{*block.NewTransaction("genesis", "reward", 0)}
	genesisBlock := block.NewBlock(0, 0, 0, "0", txs) 

	// 1. Genesis block scenario: (0,0,0) is empty
	if !cs.CanAddBlock(0, 0, 0, genesisBlock) {
		t.Error("CanAddBlock should allow Genesis block at (0,0,0) when plane is empty")
	}

	// Add Genesis block to the plane
	err := p.AddBlock(genesisBlock.X, genesisBlock.Y, genesisBlock.Z, genesisBlock)
	if err != nil {
		t.Fatalf("Failed to add genesis block for testing: %v", err)
	}
	time.Sleep(1 * time.Nanosecond) // Ensure next block has different timestamp/hash

	// 2. Attempt to add another block at (0,0,0) - now it's occupied
	// CanAddBlock's current logic: if (0,0,0) is occupied, it's no longer genesis scenario.
	// It will then require neighbors. Since (0,0,0) is the only block, it has no neighbors.
	// So, CanAddBlock(0,0,0, someNewBlock) should be false.
	// The AddBlock method of plane would prevent adding to (0,0,0) anyway, this tests CanAddBlock's logic.
	if cs.CanAddBlock(0, 0, 0, block.NewBlock(0,0,0,genesisBlock.Hash,nil)) {
		t.Error("CanAddBlock should not allow adding to an already occupied (0,0,0) if it has no other neighbors")
	}
        
	// 3. Valid adjacent block
	adjBlock := block.NewBlock(0, 0, 1, genesisBlock.Hash, nil)
	if !cs.CanAddBlock(0, 0, 1, adjBlock) {
		t.Error("CanAddBlock should allow block adjacent to Genesis at (0,0,1)")
	}
    // Add it for next test
    err = p.AddBlock(adjBlock.X, adjBlock.Y, adjBlock.Z, adjBlock)
	if err != nil {
		t.Fatalf("Failed to add adjacent block (0,0,1) for testing: %v", err)
	}
	time.Sleep(1 * time.Nanosecond)


	// 4. Invalid isolated block (not adjacent to any existing block)
	isolatedBlock := block.NewBlock(5, 5, 5, "someprevhash", nil)
	if cs.CanAddBlock(5, 5, 5, isolatedBlock) {
		t.Error("CanAddBlock should not allow isolated block at (5,5,5)")
	}

	// 5. Invalid coordinates (negative)
	negCoordBlock := block.NewBlock(-1, 0, 0, "prev", nil)
	if cs.CanAddBlock(-1, 0, 0, negCoordBlock) {
		t.Error("CanAddBlock should not allow block at negative coordinates (-1,0,0)")
	}

    // 6. Valid block, adjacent to the second block (0,0,1)
    adjToSecondBlock := block.NewBlock(0,0,2, adjBlock.Hash, nil)
    if !cs.CanAddBlock(0,0,2, adjToSecondBlock){
        t.Errorf("CanAddBlock should allow block at (0,0,2) adjacent to (0,0,1). Neighbors of (0,0,2): %v", p.ListNeighbors(0,0,2))
    }
}
