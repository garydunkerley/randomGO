// this script contains functions related to checking the legality of a move in accordance with the rules:
// 1. No playing in a space currently occupied by a stone
// 2. You may not play a move that results in the suicide of your stone group.
// 3. The "super ko rule": You may not replicate a previously extant board state.

package game

import (
	"errors"
	// "fmt"
)

func (y boardState) checkLegalMove(n int) error {
	node, ok := y.nodes[n]
	if !ok {
		return errors.New("Illegal move: nonexistent node.")
	}
	if node.color != 0 {
		return errors.New("Illegal move: Nonempty point.")
	}

	// We temporarily change the color of the node to do error testing
	if a {
		y.current_board.nodes[n].color = 2
	} else {
		y.current_board.nodes[n].color = 1
	}

	if y.suicidal_move(n) {
		y.current_board.nodes[n].color = 0
		// we revert node back to being blank
		// before returning error

		return errors.New("Illegal move: Suicide.")

	} else if y.illegal_ko_move(n) {
		y.current_board.nodes[n].color = 0
		// we revert the node back to being
		// blank before returning an error

		return errors.New("Illegal move: Violates ko rule.")
	} else {
		y.current_board.nodes[n].color = 0
		// we revert the node back to being blank
		// so that the game can proceed to implement
		// the player's move

		return nil
	}
}

func (y boardState) suicidal_move(n int) bool {
	// suicidal_move takes an integer and looks at the corresponding node to determine whether the stone in question is surrounded

	for _, z := range y.current_board.nodes[n].neighbors {
		if z.color == 0 {
			// if the stone being played has a liberty, it won't be regarded as a suicidal move.
			return false
		} else if z.color != y.current_board.nodes[n].color {
			// if a stone is played in the last liberty of an enemy group, it is not a suicide
			if len(z.inGroup.liberties) == 1 {
				return false
			}
		} else if len(z.inGroup.liberties) > 1 { // if an adjacent friendly group has at least two liberties, it is not a suicide
			return false

		}
	}
	// We return error "Suicidal Move" if all three hold
	// 1. the node has no liberties
	// 2. it does not kill an enemy group
	// 3. it occupies the last liberty of a friendly group
	return true
}

func (y boardState) illegal_ko_move(n int) bool {
	// Need to modify y.board_history so that it stores encodings of the board states

	return false
}
