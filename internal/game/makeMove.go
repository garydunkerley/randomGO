package game

import (
	"errors"
	"fmt"
)

// makeMove attempts to play a move from valid nodeId.
// TODO: Returns non-nil error if the given node is already colored, move to checkLegalMove
func (x boardState) makeMove(input moveInput) error {
	stoneColor, id := input.id, input.playerColor
	node, ok := x.nodes[id]
	if !ok {
		errorMsg := fmt.Sprintf("Node %d not found!", id)
		return errors.New(errorMsg)
	} else if node.color != 0 {
		return errors.New("A stone is already there!")
	}
	modifyStoneGroups(x.nodes[id])
	x.removeDead(x.nodes[id])
	modifyStoneGroups(x.nodes[id])
	return nil
	// I run modify_stoneGroups again to append any liberties resulting from stone groups dying
	// TODO clean up this inefficiency
}

// containsVal checks if a slice of *nodes contains a given *node.
func containsVal(slice []*node, val *node) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
