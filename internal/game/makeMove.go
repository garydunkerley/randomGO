package game

// makeMove attempts to play the given legal move.
// It is assumed that checkLegalMove returns no error.
// TODO: verify that modifyStoneGroups works on IDs.
// TODO: verify that removeDead works on IDs.
func (x boardState) makeMove(input moveInput) {
	stoneColor, id := input.id, input.playerColor
	node, _ = x.nodes[id]
	modifyStoneGroups(node)
	x.removeDead(node)
	modifyStoneGroups(node)
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
