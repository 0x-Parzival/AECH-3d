package test

import (
	"aech/block"
	"aech/plane"
	"reflect"
	"testing"
	"time"
)

// TestNewBlockAndHashGeneration tests the creation of new blocks and hash generation logic.
func TestNewBlockAndHashGeneration(t *testing.T) {
	// Use NewTransaction to create transactions
	tx1 := block.NewTransaction("sender1", "receiver1", 10.0)
	txs1 := []block.Transaction{*tx1}
	b1 := block.NewBlock(0, 0, 0, "genesis", txs1)

	if b1.Hash == "" {
		t.Errorf("Expected b1.Hash to be non-empty, got empty")
	}
	if b1.Timestamp.IsZero() {
		t.Errorf("Expected b1.Timestamp to be non-zero, got zero")
	}

	time.Sleep(1 * time.Nanosecond)
	// Recreate txs1 for b2 if its timestamp matters for distinctness from b1's txs,
	// or use the same txs1 if the block's timestamp is the primary differentiator.
	// For this test, block timestamp difference is key.
	b2 := block.NewBlock(0, 0, 0, "genesis", txs1)
	if b1.Hash == b2.Hash {
		t.Errorf("Expected b1.Hash (%s) and b2.Hash (%s) to be different due to different timestamps, but they are the same", b1.Hash, b2.Hash)
	}

	tx2 := block.NewTransaction("sender2", "receiver2", 20.0)
	txs2 := []block.Transaction{*tx2}
	b3 := block.NewBlock(1, 0, 0, "genesis", txs2)
	if b1.Hash == b3.Hash {
		t.Errorf("Expected b1.Hash (%s) and b3.Hash (%s) to be different for different inputs, but they are the same", b1.Hash, b3.Hash)
	}
}

// TestInsertAndGetBlock tests inserting and retrieving blocks from the plane using Plane3D.
func TestInsertAndGetBlock(t *testing.T) {
	p := plane.NewPlane()
	tx3 := block.NewTransaction("sender3", "receiver3", 30.0)
	txs3 := []block.Transaction{*tx3}
	testBlock := block.NewBlock(1, 2, 3, "prevhash", txs3)

	// Initial add
	err := p.AddBlock(testBlock.X, testBlock.Y, testBlock.Z, testBlock)
	if err != nil {
		t.Fatalf("AddBlock failed for initial insert: %v", err)
	}

	retrievedBlock, err := p.GetBlock(1, 2, 3)
	if err != nil {
		t.Fatalf("Expected block to be found at (1,2,3), but GetBlock returned error: %v", err)
	}
	if retrievedBlock == nil {
		t.Fatalf("Expected block to be non-nil at (1,2,3), but retrievedBlock was nil")
	}
	if retrievedBlock.Hash != testBlock.Hash {
		t.Errorf("Expected retrievedBlock.Hash (%s) to be equal to testBlock.Hash (%s)", retrievedBlock.Hash, testBlock.Hash)
	}
	if !reflect.DeepEqual(retrievedBlock.Transactions, testBlock.Transactions) {
		t.Errorf("Expected retrievedBlock.Transactions to be DeepEqual to testBlock.Transactions")
	}

	// Test retrieval of non-existent block
	nonExistentBlock, err := p.GetBlock(9, 9, 9)
	if err == nil {
		t.Errorf("Expected error when getting non-existent block at (9,9,9), but got nil error")
	}
	if nonExistentBlock != nil {
		t.Errorf("Expected nonExistentBlock to be nil, but it was not")
	}

	// Test overwrite/duplicate add
	tx4 := block.NewTransaction("sender4", "receiver4", 40.0)
	txs4 := []block.Transaction{*tx4}
	// Ensure anotherBlock has a different hash if that's part of the test logic for "original block unchanged" later
	// For AddBlock, only coordinates matter for conflict. Hash/content difference is for other tests.
	anotherBlock := block.NewBlock(testBlock.X, testBlock.Y, testBlock.Z, "anotherprevhash", txs4)
	err = p.AddBlock(testBlock.X, testBlock.Y, testBlock.Z, anotherBlock)
	if err == nil {
		t.Errorf("Expected error when adding block to occupied coordinates (%d,%d,%d), but got nil", testBlock.X, testBlock.Y, testBlock.Z)
	}

	// Ensure the original block is still there and unchanged
	originalBlockCheck, err := p.GetBlock(testBlock.X, testBlock.Y, testBlock.Z)
	if err != nil {
		t.Fatalf("GetBlock failed for original block after duplicate add attempt: %v", err)
	}
	if originalBlockCheck.Hash != testBlock.Hash {
		t.Errorf("Original block hash changed after duplicate add attempt. Expected %s, got %s", testBlock.Hash, originalBlockCheck.Hash)
	}
}

// TestValidateBlock tests the block validation logic.
func TestValidateBlock(t *testing.T) {
	validBlock := block.NewBlock(0, 0, 0, "prev", nil)
	if !block.ValidateBlock(*validBlock) {
		t.Errorf("Expected validBlock to be valid, but it was not. Hash: %s", validBlock.Hash)
	}

	invalidBlock1 := block.NewBlock(0, 0, 0, "prev", nil)
	invalidBlock1.Hash = "" // Manually invalidate
	if block.ValidateBlock(*invalidBlock1) {
		t.Errorf("Expected invalidBlock1 (empty hash) to be invalid, but it was valid")
	}

	invalidBlock2 := block.NewBlock(-1, 0, 0, "prev", nil)
	// ValidateBlock itself doesn't use IsValidPosition from plane, it checks block's own fields based on its current simple logic.
	// The block.ValidateBlock currently only checks for Hash == "" or X < 0 || Y < 0 || Z < 0.
	// NewBlock sets X,Y,Z so block.ValidateBlock will use these.
	if block.ValidateBlock(*invalidBlock2) {
		t.Errorf("Expected invalidBlock2 (X=-1) to be invalid due to negative coordinate, but it was valid. Hash: %s, X: %d", invalidBlock2.Hash, invalidBlock2.X)
	}
}

// TestListNeighbors tests retrieving neighbors of a block using Plane3D.
func TestListNeighbors(t *testing.T) {
	p := plane.NewPlane()

	b000 := block.NewBlock(0, 0, 0, "genesis", nil)
	time.Sleep(1 * time.Nanosecond) // Ensure distinct timestamps for distinct hashes
	b100 := block.NewBlock(1, 0, 0, b000.Hash, nil)
	time.Sleep(1 * time.Nanosecond)
	b010 := block.NewBlock(0, 1, 0, b000.Hash, nil)

	if err := p.AddBlock(b000.X, b000.Y, b000.Z, b000); err != nil {
		t.Fatalf("Failed to add b000: %v", err)
	}
	if err := p.AddBlock(b100.X, b100.Y, b100.Z, b100); err != nil {
		t.Fatalf("Failed to add b100: %v", err)
	}
	if err := p.AddBlock(b010.X, b010.Y, b010.Z, b010); err != nil {
		t.Fatalf("Failed to add b010: %v", err)
	}

	neighbors := p.ListNeighbors(0, 0, 0)
	if len(neighbors) != 2 {
		t.Fatalf("Expected 2 neighbors for block (0,0,0), got %d", len(neighbors))
	}

	foundB100 := false
	foundB010 := false
	for _, nb := range neighbors {
		if nb.Hash == b100.Hash {
			foundB100 = true
		}
		if nb.Hash == b010.Hash {
			foundB010 = true
		}
	}

	if !foundB100 {
		t.Errorf("Neighbor b100 (hash %s) not found in neighbors of (0,0,0)", b100.Hash)
	}
	if !foundB010 {
		t.Errorf("Neighbor b010 (hash %s) not found in neighbors of (0,0,0)", b010.Hash)
	}

	// Test neighbors for a block with no neighbors
	isolatedNeighbors := p.ListNeighbors(5, 5, 5)
	if len(isolatedNeighbors) != 0 {
		t.Errorf("Expected 0 neighbors for isolated coordinates (5,5,5), got %d", len(isolatedNeighbors))
	}
}

func TestIsValidPosition(t *testing.T) {
	p := plane.NewPlane() // Instantiate Plane3D

	// Test cases for valid positions (non-negative coordinates)
	validCoordinates := [][3]int{
		{0, 0, 0},
		{1, 2, 3},
		{100, 100, 100},
	}
	for _, coord := range validCoordinates {
		if !p.IsValidPosition(coord[0], coord[1], coord[2]) {
			t.Errorf("Expected IsValidPosition(%d, %d, %d) to be true, but got false", coord[0], coord[1], coord[2])
		}
	}

	// Test cases for invalid positions (negative coordinates)
	invalidCoordinates := [][3]int{
		{-1, 0, 0},
		{0, -1, 0},
		{0, 0, -1},
		{-1, -1, -1},
		{1, -5, 3},
	}
	for _, coord := range invalidCoordinates {
		if p.IsValidPosition(coord[0], coord[1], coord[2]) {
			t.Errorf("Expected IsValidPosition(%d, %d, %d) to be false, but got true", coord[0], coord[1], coord[2])
		}
	}
}
