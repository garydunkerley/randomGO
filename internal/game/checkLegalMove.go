// this script contains functions related to checking the legality of a move in accordance with the rules:
// 1. No playing in a space currently occupied by a stone
// 2. You may not play a move that results in the suicide of your stone group.
// 3. The "super ko rule": You may not replicate a previously extant board state.

package game

import (
	"errors"
)

// checkLegalMove takes a node id and color, returns non-nil error if a stone
// of the given color cannot be played at that id
func (x boardState) checkLegalMove(n int, c int8) error {
	node, ok := x.nodes[n]
	if !ok {
		return errors.New("Illegal move: nonexistent node.")
	}
	if node.color != 0 {
		return errors.New("Illegal move: Nonempty point.")
	}

	// We temporarily change the color of the node to do error testing
	// TODO Ensure this works after removing board state hacks.
	// The suicide check and ko check should not require global state change.

	if x.suicidalMove(n, c) {
		return errors.New("Illegal move: Suicide.")

	} else if x.illegal_ko_move(n) {
		return errors.New("Illegal move: Violates ko rule.")

	} else {
		return nil
	}
}

// suicidalMove checks if given (nodeId, color) is the last liberty of a group
// of the given color.
func (y boardState) suicidalMove(n int, c int8) bool {

	for _, z := range y.nodes[n].neighbors {
		if z.color == 0 {
			// if the stone being played has a liberty,
			// it won't be regarded as a suicidal move.
			return false
		} else if z.color != c {
			// if a stone is played in the last liberty of an enemy group,
			// it is not a suicide
			if y.countLiberties(y.stringOf[z.id]) == 1 {
				return false
			}
		} else if y.countLiberties(y.stringOf[z.id]) > 1 {
			// if an adjacent friendly group has at least two liberties,
			// it is not a suicide
			return false
		}
	}
	// We return error "Suicidal Move" if all three hold
	// 1. the node has no liberties
	// 2. it does not kill an enemy group
	// 3. it occupies the last liberty of a friendly group
	return true
}

//illegal_ko_move will eventually check for ko
func (y boardState) illegal_ko_move(n int) bool {
	if n == y.history.koPoint {
		return true
	}
	return false
}
