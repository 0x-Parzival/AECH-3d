package plane

import (
	"aech/block"
	"fmt"
)

// Cube stores Block3D objects in a 3D spatial grid.
// Access blocks by Cube[x][y][z].
var Cube = make(map[int]map[int]map[int]block.Block3D)

// InsertBlock inserts a Block3D into the Cube at the specified coordinates.
// It overwrites if a block already exists at the coordinates.
func InsertBlock(x, y, z int, b block.Block3D) error {
	if _, ok := Cube[x]; !ok {
		Cube[x] = make(map[int]map[int]block.Block3D)
	}
	if _, ok := Cube[x][y]; !ok {
		Cube[x][y] = make(map[int]block.Block3D)
	}
	Cube[x][y][z] = b
	return nil // No specific error conditions handled for now
}

// GetBlock retrieves a Block3D from the Cube at the specified coordinates.
// It returns the block and a boolean indicating if the block was found.
func GetBlock(x, y, z int) (block.Block3D, bool) {
	if _, ok := Cube[x]; !ok {
		return block.Block3D{}, false
	}
	if _, ok := Cube[x][y]; !ok {
		return block.Block3D{}, false
	}
	b, ok := Cube[x][y][z]
	return b, ok
}

// GetNeighbors retrieves all existing adjacent blocks to the given coordinates.
func GetNeighbors(x, y, z int) []block.Block3D {
	var neighbors []block.Block3D

	offsets := [][3]int{
		{x + 1, y, z}, {x - 1, y, z}, // Right, Left
		{x, y + 1, z}, {x, y - 1, z}, // Up, Down
		{x, y, z + 1}, {x, y, z - 1}, // Front, Back
	}

	for _, off := range offsets {
		if b, found := GetBlock(off[0], off[1], off[2]); found {
			neighbors = append(neighbors, b)
		}
	}
	return neighbors
}

// PrintCoordinates is a placeholder function that might be used for debugging.
// This is added to use the fmt import and avoid "imported and not used" error during potential linting/vetting.
// This can be removed later if fmt is used by other functions or if error handling becomes more sophisticated.
func PrintCoordinates(x, y, z int) {
	fmt.Printf("Coordinates: X=%d, Y=%d, Z=%d\n", x, y, z)
}
