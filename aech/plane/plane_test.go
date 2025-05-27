package plane

import (
	"aech/block"
	"fmt"
	"reflect"
	"sort" // For comparing slices of blocks if order doesn't matter
	"testing"
	// "time" // May be needed for block timestamps if not handled by NewBlock
)

// Helper function to create a simple block for testing
func newTestBlock(x, y, z int, previousHash string, idSuffix string) *block.Block3D {
	txs := []block.Transaction{*block.NewTransaction("testSender"+idSuffix, "testReceiver"+idSuffix, 1.0)}
	// Ensure NewBlock properly initializes the block, including generating a hash.
	// The coordinates passed to NewBlock are for its internal record,
	// Plane3D uses its own x,y,z for grid placement.
	return block.NewBlock(x, y, z, previousHash, txs)
}

func TestNewPlane(t *testing.T) {
	t.Helper()
	p := NewPlane()
	if p == nil {
		t.Fatal("NewPlane() returned nil")
	}
	if p.Grid == nil {
		t.Error("NewPlane() did not initialize Grid map")
	}
	if len(p.Grid) != 0 {
		t.Errorf("NewPlane() Grid map should be empty, got len %d", len(p.Grid))
	}
}

func TestCoordKey(t *testing.T) {
	t.Helper()
	testCases := []struct {
		name     string
		x, y, z  int
		expected string
	}{
		{"Origin", 0, 0, 0, "0_0_0"},
		{"PositiveCoords", 1, 2, 3, "1_2_3"},
		{"MixedCoords", 10, 0, 55, "10_0_55"},
		// IsValidPosition currently prevents negative coords for placement,
		// but coordKey itself might be called with them internally or in other contexts.
		// For now, focusing on valid position keys.
		// {"NegativeX", -1, 5, 10, "-1_5_10"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()
			key := coordKey(tc.x, tc.y, tc.z)
			if key != tc.expected {
				t.Errorf("coordKey(%d,%d,%d) expected %s, got %s", tc.x, tc.y, tc.z, tc.expected, key)
			}
		})
	}

	// Test that different coordinates produce different keys
	key1 := coordKey(1, 1, 1)
	key2 := coordKey(1, 1, 2)
	if key1 == key2 {
		t.Errorf("coordKey(1,1,1) and coordKey(1,1,2) produced the same key: %s", key1)
	}
}

func TestAddBlock(t *testing.T) {
	t.Helper()
	p := NewPlane()
	sampleBlock1 := newTestBlock(0, 0, 0, "0", "1")
	sampleBlock2 := newTestBlock(0, 0, 0, "0", "2") // Different content implies different hash

	t.Run("AddValidBlockToEmptyPosition", func(t *testing.T) {
		t.Helper()
		x, y, z := 0, 0, 0
		err := p.AddBlock(x, y, z, sampleBlock1)
		if err != nil {
			t.Fatalf("AddBlock(%d,%d,%d, sampleBlock1) failed: %v", x, y, z, err)
		}

		key := coordKey(x, y, z)
		if _, exists := p.Grid[key]; !exists {
			t.Errorf("Block not found in Grid at key %s after AddBlock", key)
		}
		if p.Grid[key] != sampleBlock1 {
			t.Errorf("Block in Grid at key %s is not the one that was added", key)
		}
	})

	t.Run("AddBlockToOccupiedPosition", func(t *testing.T) {
		t.Helper()
		// Position 0,0,0 is already occupied by sampleBlock1 from the previous sub-test.
		// This assumes sub-tests run sequentially and share the parent test's 'p' instance state.
		x, y, z := 0, 0, 0
		err := p.AddBlock(x, y, z, sampleBlock2)
		if err == nil {
			t.Fatal("AddBlock to occupied position expected an error, got nil")
		}

		// Check that the original block is still there
		key := coordKey(x, y, z)
		if p.Grid[key] != sampleBlock1 {
			t.Errorf("Original block at key %s was overwritten or removed. Expected sampleBlock1, got %+v", key, p.Grid[key])
		}
	})

	// Optional: Test adding a block whose internal coordinates don't match.
	// Current AddBlock doesn't check this. If it did, a test would be:
	// t.Run("AddBlockWithMismatchedCoordinates", func(t *testing.T) { ... })
}

func TestGetBlock(t *testing.T) {
	t.Helper()
	p := NewPlane()
	x, y, z := 1, 1, 1
	sampleBlock := newTestBlock(x, y, z, "0", "gb1") // Block's internal coords match plane coords for clarity
	_ = p.AddBlock(x, y, z, sampleBlock)

	t.Run("GetExistingBlock", func(t *testing.T) {
		t.Helper()
		retrievedBlock, err := p.GetBlock(x, y, z)
		if err != nil {
			t.Fatalf("GetBlock(%d,%d,%d) failed: %v", x, y, z, err)
		}
		if retrievedBlock == nil {
			t.Fatal("GetBlock returned nil block for existing position")
		}
		if retrievedBlock != sampleBlock {
			t.Errorf("GetBlock returned wrong block. Expected hash %s, got %s", sampleBlock.Hash, retrievedBlock.Hash)
		}
	})

	t.Run("GetBlockFromEmptyPosition", func(t *testing.T) {
		t.Helper()
		emptyX, emptyY, emptyZ := 2, 2, 2
		retrievedBlock, err := p.GetBlock(emptyX, emptyY, emptyZ)
		if err == nil {
			t.Error("GetBlock from empty position expected an error, got nil")
		}
		if retrievedBlock != nil {
			t.Errorf("GetBlock from empty position expected nil block, got %+v", retrievedBlock)
		}
		expectedErrorMsg := fmt.Sprintf("no block found at coordinates %d,%d,%d", emptyX, emptyY, emptyZ)
		if err != nil && err.Error() != expectedErrorMsg {
			t.Errorf("GetBlock error message mismatch. Expected '%s', got '%s'", expectedErrorMsg, err.Error())
		}
	})
}

func TestIsValidPosition(t *testing.T) {
	t.Helper()
	p := NewPlane() // Plane instance isn't strictly needed for current IsValidPosition, but good practice.

	validPositions := []struct{ x, y, z int }{
		{0, 0, 0}, {10, 20, 30}, {planeSize - 1, planeSize - 1, planeSize - 1},
	}
	for _, pos := range validPositions {
		t.Run(fmt.Sprintf("ValidPos_%d_%d_%d", pos.x, pos.y, pos.z), func(t *testing.T) {
			t.Helper()
			if !p.IsValidPosition(pos.x, pos.y, pos.z) {
				t.Errorf("IsValidPosition(%d,%d,%d) returned false for a valid position", pos.x, pos.y, pos.z)
			}
		})
	}

	invalidPositions := []struct{ x, y, z int }{
		{-1, 0, 0}, {0, -1, 0}, {0, 0, -1},
		{planeSize, 0, 0}, {0, planeSize, 0}, {0, 0, planeSize}, // Assuming planeSize is exclusive upper bound
	}
	for _, pos := range invalidPositions {
		t.Run(fmt.Sprintf("InvalidPos_%d_%d_%d", pos.x, pos.y, pos.z), func(t *testing.T) {
			t.Helper()
			if p.IsValidPosition(pos.x, pos.y, pos.z) {
				t.Errorf("IsValidPosition(%d,%d,%d) returned true for an invalid position", pos.x, pos.y, pos.z)
			}
		})
	}
}

func TestListNeighbors(t *testing.T) {
	t.Helper()
	p := NewPlane()

	// Setup: Add a central block and some neighbors
	centerBlock := newTestBlock(1, 1, 1, "0", "center")
	_ = p.AddBlock(1, 1, 1, centerBlock)

	neighbor1 := newTestBlock(0, 1, 1, "0", "n1") // X-1
	_ = p.AddBlock(0, 1, 1, neighbor1)
	neighbor2 := newTestBlock(2, 1, 1, "0", "n2") // X+1
	_ = p.AddBlock(2, 1, 1, neighbor2)
	neighbor3 := newTestBlock(1, 0, 1, "0", "n3") // Y-1
	_ = p.AddBlock(1, 0, 1, neighbor3)
	// Not adding Y+1, Z-1, Z+1 neighbors for this specific test case

	expectedNeighbors := []*block.Block3D{neighbor1, neighbor2, neighbor3}

	t.Run("ForCentralBlockWithNeighbors", func(t *testing.T) {
		t.Helper()
		neighbors := p.ListNeighbors(1, 1, 1)
		if len(neighbors) != len(expectedNeighbors) {
			t.Fatalf("ListNeighbors(1,1,1) expected %d neighbors, got %d. Neighbors: %+v", len(expectedNeighbors), len(neighbors), neighbors)
		}

		// Sort slices by hash for consistent comparison
		sort.Slice(neighbors, func(i, j int) bool { return neighbors[i].Hash < neighbors[j].Hash })
		sort.Slice(expectedNeighbors, func(i, j int) bool { return expectedNeighbors[i].Hash < expectedNeighbors[j].Hash })

		for i, expected := range expectedNeighbors {
			if neighbors[i].Hash != expected.Hash { // Compare by hash as pointers will differ
				t.Errorf("Neighbor mismatch. Expected block with hash %s, got %s at index %d", expected.Hash, neighbors[i].Hash, i)
			}
		}
	})

	t.Run("ForBlockWithNoNeighbors", func(t *testing.T) {
		t.Helper()
		isolatedBlock := newTestBlock(5, 5, 5, "0", "isolated")
		_ = p.AddBlock(5, 5, 5, isolatedBlock)
		neighbors := p.ListNeighbors(5, 5, 5)
		if len(neighbors) != 0 {
			t.Errorf("ListNeighbors(5,5,5) expected 0 neighbors for isolated block, got %d", len(neighbors))
		}
	})

	t.Run("ForPositionWithoutBlockButWithNeighbors", func(t *testing.T) {
		t.Helper()
		// Position (1,1,0) has no block, but neighbor3 is at (1,0,1), centerBlock at (1,1,1)
		// ListNeighbors checks around the given coords, not *of* a block at those coords.
		// The neighbors of (1,1,0) would include blocks at (0,1,0), (2,1,0), (1,0,0), (1,2,0), (1,1,-1), (1,1,1)
		// In our current setup, only centerBlock (1,1,1) is a neighbor to (1,1,0).
		neighbors := p.ListNeighbors(1, 1, 0)
		if len(neighbors) != 1 {
			t.Fatalf("ListNeighbors(1,1,0) expected 1 neighbor, got %d. Neighbors: %+v", len(neighbors), neighbors)
		}
		if neighbors[0].Hash != centerBlock.Hash {
			t.Errorf("Expected neighbor to be centerBlock (hash %s), got hash %s", centerBlock.Hash, neighbors[0].Hash)
		}
	})
}

func TestGetBlockByID(t *testing.T) {
	t.Helper()
	p := NewPlane()

	b1 := newTestBlock(0, 0, 0, "0", "id1")
	b2 := newTestBlock(0, 0, 1, b1.Hash, "id2")
	b3 := newTestBlock(0, 1, 0, b1.Hash, "id3")

	_ = p.AddBlock(b1.X, b1.Y, b1.Z, b1)
	_ = p.AddBlock(b2.X, b2.Y, b2.Z, b2)
	// Not adding b3 to test a non-existent case for its ID later, if needed.

	t.Run("GetExistingBlockByID", func(t *testing.T) {
		t.Helper()
		retrieved, err := p.GetBlockByID(b1.Hash)
		if err != nil {
			t.Fatalf("GetBlockByID for b1.Hash (%s) failed: %v", b1.Hash, err)
		}
		if retrieved == nil {
			t.Fatal("GetBlockByID for b1.Hash returned nil block")
		}
		if retrieved.Hash != b1.Hash {
			t.Errorf("GetBlockByID expected block with hash %s, got %s", b1.Hash, retrieved.Hash)
		}
	})

	t.Run("GetAnotherExistingBlockByID", func(t *testing.T) {
		t.Helper()
		retrieved, err := p.GetBlockByID(b2.Hash)
		if err != nil {
			t.Fatalf("GetBlockByID for b2.Hash (%s) failed: %v", b2.Hash, err)
		}
		if retrieved.Hash != b2.Hash {
			t.Errorf("GetBlockByID expected block with hash %s, got %s", b2.Hash, retrieved.Hash)
		}
	})

	t.Run("GetNonExistentBlockByID", func(t *testing.T) {
		t.Helper()
		nonExistentHash := b3.Hash // b3 was not added
		retrieved, err := p.GetBlockByID(nonExistentHash)
		if err == nil {
			t.Errorf("GetBlockByID for nonExistentHash (%s) expected an error, got nil", nonExistentHash)
		}
		if retrieved != nil {
			t.Errorf("GetBlockByID for nonExistentHash (%s) expected nil block, got %+v", nonExistentHash, retrieved)
		}
		expectedErrorMsg := fmt.Sprintf("block with ID %s not found", nonExistentHash)
		if err != nil && err.Error() != expectedErrorMsg {
			t.Errorf("GetBlockByID error message mismatch. Expected '%s', got '%s'", expectedErrorMsg, err.Error())
		}
	})

	t.Run("GetBlockByIDFromEmptyPlane", func(t *testing.T) {
		t.Helper()
		emptyPlane := NewPlane()
		retrieved, err := emptyPlane.GetBlockByID("anyHash")
		if err == nil {
			t.Error("GetBlockByID from empty plane expected an error, got nil")
		}
		if retrieved != nil {
			t.Errorf("GetBlockByID from empty plane expected nil block, got %+v", retrieved)
		}
	})
}

func TestGetAllBlocks(t *testing.T) {
	t.Helper()
	p := NewPlane()

	t.Run("EmptyPlane", func(t *testing.T) {
		t.Helper()
		allBlocks := p.GetAllBlocks()
		if allBlocks == nil { // Should be an empty slice, not nil
			t.Error("GetAllBlocks on empty plane returned nil, expected empty slice")
		}
		if len(allBlocks) != 0 {
			t.Errorf("GetAllBlocks on empty plane expected 0 blocks, got %d", len(allBlocks))
		}
	})

	b1 := newTestBlock(0, 0, 0, "0", "ga1")
	b2 := newTestBlock(1, 0, 0, b1.Hash, "ga2")
	b3 := newTestBlock(0, 1, 0, b1.Hash, "ga3")

	_ = p.AddBlock(b1.X, b1.Y, b1.Z, b1)
	_ = p.AddBlock(b2.X, b2.Y, b2.Z, b2)
	_ = p.AddBlock(b3.X, b3.Y, b3.Z, b3)

	t.Run("PopulatedPlane", func(t *testing.T) {
		t.Helper()
		allBlocks := p.GetAllBlocks()
		if len(allBlocks) != 3 {
			t.Fatalf("GetAllBlocks expected 3 blocks, got %d", len(allBlocks))
		}

		expectedHashes := map[string]bool{b1.Hash: true, b2.Hash: true, b3.Hash: true}
		retrievedHashes := make(map[string]bool)
		for _, b := range allBlocks {
			if _, found := expectedHashes[b.Hash]; !found {
				t.Errorf("GetAllBlocks returned unexpected block with hash %s", b.Hash)
			}
			retrievedHashes[b.Hash] = true
		}
		if len(retrievedHashes) != len(expectedHashes) {
			t.Errorf("GetAllBlocks did not return all expected blocks. Expected %d unique hashes, got %d", len(expectedHashes), len(retrievedHashes))
		}
	})
}
```
