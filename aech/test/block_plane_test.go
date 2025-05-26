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
	txs := []block.Transaction{{Data: []byte("tx1")}}
	b1 := block.NewBlock(0, 0, 0, "genesis", txs)

	if b1.Hash == "" {
		t.Errorf("Expected b1.Hash to be non-empty, got empty")
	}
	if b1.Timestamp.IsZero() {
		t.Errorf("Expected b1.Timestamp to be non-zero, got zero")
	}

	// NewBlock uses time.Now(), so two blocks created sequentially will likely have different timestamps and thus hashes.
	// We'll ensure this by sleeping for a tiny duration.
	time.Sleep(1 * time.Nanosecond)
	b2 := block.NewBlock(0, 0, 0, "genesis", txs)
	if b1.Hash == b2.Hash {
		t.Errorf("Expected b1.Hash (%s) and b2.Hash (%s) to be different due to different timestamps, but they are the same", b1.Hash, b2.Hash)
	}

	b3 := block.NewBlock(1, 0, 0, "genesis", []block.Transaction{{Data: []byte("tx2")}})
	if b1.Hash == b3.Hash {
		t.Errorf("Expected b1.Hash (%s) and b3.Hash (%s) to be different for different inputs, but they are the same", b1.Hash, b3.Hash)
	}

	// Test that re-hashing the same block data (manually constructed) yields the same hash
	// This is a bit more involved as NewBlock sets the timestamp internally.
	// Instead, let's focus on GenerateHash's determinism if we can control all inputs.
	// For now, the check that different inputs (b1 vs b3) yield different hashes is a good start.
}

// TestInsertAndGetBlock tests inserting and retrieving blocks from the plane.
func TestInsertAndGetBlock(t *testing.T) {
	// Clear the Cube for a clean test environment
	plane.Cube = make(map[int]map[int]map[int]block.Block3D)

	testBlock := block.NewBlock(1, 2, 3, "prevhash", []block.Transaction{{Data: []byte("data")}})
	err := plane.InsertBlock(testBlock.X, testBlock.Y, testBlock.Z, *testBlock)
	if err != nil {
		t.Fatalf("InsertBlock failed: %v", err)
	}

	retrievedBlock, found := plane.GetBlock(1, 2, 3)
	if !found {
		t.Fatalf("Expected block to be found at (1,2,3), but it was not")
	}
	if retrievedBlock.Hash != testBlock.Hash {
		t.Errorf("Expected retrievedBlock.Hash (%s) to be equal to testBlock.Hash (%s)", retrievedBlock.Hash, testBlock.Hash)
	}
	if !reflect.DeepEqual(retrievedBlock.Transactions, testBlock.Transactions) {
		t.Errorf("Expected retrievedBlock.Transactions to be DeepEqual to testBlock.Transactions")
	}

	_, found = plane.GetBlock(9, 9, 9)
	if found {
		t.Errorf("Expected no block to be found at (9,9,9), but it was")
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

	// Based on current stub, negative coordinates make it invalid.
	// NewBlock will generate a hash, so we don't need to worry about that part.
	invalidBlock2 := block.NewBlock(-1, 0, 0, "prev", nil)
	if block.ValidateBlock(*invalidBlock2) {
		t.Errorf("Expected invalidBlock2 (X=-1) to be invalid, but it was valid. Hash: %s, X: %d", invalidBlock2.Hash, invalidBlock2.X)
	}
}

// TestGetNeighbors tests retrieving neighbors of a block.
func TestGetNeighbors(t *testing.T) {
	// Clear the Cube for a clean test environment
	plane.Cube = make(map[int]map[int]map[int]block.Block3D)

	b000 := block.NewBlock(0, 0, 0, "genesis", nil)
	b100 := block.NewBlock(1, 0, 0, b000.Hash, nil)
	b010 := block.NewBlock(0, 1, 0, b000.Hash, nil)
	// Ensure distinct hashes for b100 and b010 even if timestamp resolution is low
	time.Sleep(1 * time.Nanosecond) 
	b010 = block.NewBlock(0, 1, 0, b000.Hash, nil)


	plane.InsertBlock(b000.X, b000.Y, b000.Z, *b000)
	plane.InsertBlock(b100.X, b100.Y, b100.Z, *b100)
	plane.InsertBlock(b010.X, b010.Y, b010.Z, *b010)

	neighbors := plane.GetNeighbors(0, 0, 0)
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

	isolatedNeighbors := plane.GetNeighbors(5, 5, 5)
	if len(isolatedNeighbors) != 0 {
		t.Errorf("Expected 0 neighbors for isolated block (5,5,5), got %d", len(isolatedNeighbors))
	}
}
