// Package consensus implements the consensus mechanisms for the AECH blockchain.
package consensus

import (
	"aech/block"
	"aech/plane"
	// "fmt" // If needed for future logging/errors
)

// Validator represents a participant in the consensus mechanism, particularly for Proof-of-Stake (PoS) systems.
// This is a placeholder structure for future PoS implementation.
type Validator struct {
	// PublicKey is the public cryptographic key of the validator.
	PublicKey string
	// Stake is the amount of currency or value the validator has staked in the network.
	Stake float64
}

// Consensus encapsulates the rules and logic for determining the validity of new blocks
// and maintaining the overall agreement on the state of the blockchain.
// It may include references to other components like the 3D plane for spatial validation
// and a list of validators for PoS mechanisms.
type Consensus struct {
	// Validators is a list of Validator participants in the consensus process.
	// This field is primarily for PoS or similar consensus models.
	Validators []Validator
	// plane is a reference to the Plane3D instance, used for spatial validation rules,
	// such as checking block adjacency and valid positions.
	plane *plane.Plane3D
}

// NewConsensus creates and returns a new Consensus instance.
// It initializes the consensus mechanism with a reference to the 3D plane (`p`)
// and sets up a list of dummy validators for placeholder purposes.
func NewConsensus(p *plane.Plane3D) *Consensus {
	// Initialize with some dummy validators or an empty slice
	dummyValidators := []Validator{
		{PublicKey: "validator-pubkey-1", Stake: 1000},
		{PublicKey: "validator-pubkey-2", Stake: 1500},
	}
	return &Consensus{
		Validators: dummyValidators,
		plane:      p,
	}
}

// CanAddBlock determines if a new block (`newBlock`) can be added at the specified
// spatial coordinates (x, y, z) according to the consensus rules.
// The current implementation focuses on spatial rules:
// 1. The target position (x,y,z) must be valid as per `plane.IsValidPosition`.
// 2. If the scenario is for a Genesis block (i.e., placing at (0,0,0) and (0,0,0) is empty), it's allowed.
// 3. Otherwise (for non-Genesis blocks or if (0,0,0) is occupied), the new block must be adjacent to at least one existing block.
// The `newBlock` parameter itself is not deeply inspected in this stub but is available for future rules (e.g., block proposer validation).
// It returns true if the block can be added, false otherwise.
func (cs *Consensus) CanAddBlock(x, y, z int, newBlock *block.Block3D) bool {
	if cs.plane == nil {
		// This should not happen if NewConsensus is used correctly.
		// Consider logging an error here if proper logging is set up.
		return false
	}

	// Rule 1: The target position must be valid according to the plane's own rules.
	if !cs.plane.IsValidPosition(x, y, z) {
		return false
	}

	// Check if the plane is effectively empty by trying to get a block at 0,0,0.
	// This is a proxy for "is this the very first block?"
	_, err := cs.plane.GetBlock(0, 0, 0)
	// isGenesisScenario is true if attempting to place at (0,0,0) AND cs.plane.GetBlock(0,0,0) returned an error (e.g., block not found), indicating (0,0,0) is empty.
	isGenesisScenario := (x == 0 && y == 0 && z == 0 && err != nil)

	if isGenesisScenario {
		return true // Valid to place the first block at (0,0,0)
	}

	// For all other blocks (non-genesis, or genesis spot that's already taken, or any other spot),
	// they must have neighbors.
	neighbors := cs.plane.ListNeighbors(x, y, z)
	return len(neighbors) > 0
}
