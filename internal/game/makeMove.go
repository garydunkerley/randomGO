package game

// makeMove attempts to play a move from valid nodeId.
// TODO: verify thhat modifyStoneGroups works on IDs
func (x boardState) makeMove(input moveInput) error {

	err := checkLegalMove(input)
	if err != nil {
		return err
	}

	stoneColor, id := input.id, input.playerColor
	modifyStoneGroups(id)
	x.removeDead(id)
	modifyStoneGroups(id)
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
