// Package plane provides types and functions for managing the 3D spatial grid of blocks.
package plane

import (
	"fmt"
	"aech/block" // Assuming aech is the module name
)

// Plane3D struct holds blocks in a 3D grid using a map with string keys.
type Plane3D struct {
	Grid map[string]*block.Block3D
}

// NewPlane is a constructor for Plane3D.
func NewPlane() *Plane3D {
	return &Plane3D{
		Grid: make(map[string]*block.Block3D),
	}
}

// coordKey generates a unique string key for a given set of X, Y, Z coordinates.
func coordKey(x, y, z int) string {
	return fmt.Sprintf("%d_%d_%d", x, y, z)
}

// AddBlock adds a block to the Plane3D at the specified coordinates.
// It returns an error if a block already exists at that position.
func (p *Plane3D) AddBlock(x, y, z int, blk *block.Block3D) error {
	key := coordKey(x, y, z)
	if _, exists := p.Grid[key]; exists {
		return fmt.Errorf("block already exists at (%d,%d,%d)", x, y, z)
	}
	// Ensure the block's internal coordinates match, if they are used.
	// For this implementation, we assume the block passed is intended for these coordinates.
	// If blk has X,Y,Z fields, they should ideally match x,y,z or be set here.
	// blk.X, blk.Y, blk.Z = x, y, z // Example if block has these fields and they should be synced
	p.Grid[key] = blk
	return nil
}

// GetBlock retrieves a block from the Plane3D at the specified coordinates.
// It returns an error if no block is found at that position.
func (p *Plane3D) GetBlock(x, y, z int) (*block.Block3D, error) {
	key := coordKey(x, y, z)
	if blk, exists := p.Grid[key]; exists {
		return blk, nil
	}
	return nil, fmt.Errorf("no block at (%d,%d,%d)", x, y, z)
}

// IsValidPosition checks if the given coordinates are valid for placing a block.
// Currently, it only checks for non-negative coordinates.
func (p *Plane3D) IsValidPosition(x, y, z int) bool {
	// Basic validation: check for non-negative coordinates.
	// More complex rules (e.g., connectivity to existing blocks, max plane dimensions) can be added later.
	return x >= 0 && y >= 0 && z >= 0
}

// ListNeighbors retrieves all existing adjacent blocks to the given coordinates.
func (p *Plane3D) ListNeighbors(x, y, z int) []*block.Block3D {
	var neighbors []*block.Block3D
	offsets := [][3]int{
		{x + 1, y, z}, {x - 1, y, z}, // Right, Left
		{x, y + 1, z}, {x, y - 1, z}, // Up, Down
		{x, y, z + 1}, {x, y, z - 1}, // Front, Back
	}

	for _, off := range offsets {
		// We use p.GetBlock which already handles key generation and existence checks.
		// If a block exists at the offset, err will be nil.
		if blk, err := p.GetBlock(off[0], off[1], off[2]); err == nil {
			neighbors = append(neighbors, blk)
		}
	}
	return neighbors
}
