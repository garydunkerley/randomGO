// this script contains functions related to checking the legality of a move in accordance with the rules:
// 1. No playing in a space currently occupied by a stone
// 2. You may not play a move that results in the suicide of your stone group.
// 3. You may not replicate a previously extant board state.

package game

import (
	"errors"
)

func (y boardState) check_legal_move(n int) error {
	node, ok := y.current_board.nodes[n]
	if !ok {
		return errors.New("Illegal move: nonexistent node.")
	}
	if node.color != 0 {
		return errors.New("Illegal move: Nonempty point.")
	}
	// The move is legal if:
	// 1. No ko rule. (not implemented)
	// 2. No suicide.
	// To check suicide, it suffices to:
	// A. Check if this point is the last liberty for an opponent's group.
	// B. Check if this point is the last liberty for one of your own groups.
	if y.killing_move(n) {
		return nil
	} else if y.suicidal_move(n) {
		return errors.New("Illegal move: Suicide.")
	} else if y.illegal_ko_move(n) {
		return errors.New("Illegal move: Violates ko rule.")
	} else {
		return nil
	}
}

//Syntax for these next few:
//x.killing_move(true) asks "if white plays at x, does it kill a black group?"
//x.killing_move(false) asks "if black plays ... a white group?"
//Can assume that this will only be called for an empty intersection x.

func (y boardState) killing_move(n int) bool {
	//TODO:
	// for each z neighboring y.game_state.nodes[n]
	// 	if z.color != y.game_state.nodes[n].color
	//		if len(y.roup.liberties) == 1
	//			return true
	return true
}

func (y boardState) suicidal_move(n int) bool {
	//TODO:
	// For each y neighboring x
	//	if y.color == x.color
	//		if len(y.group.liberties) == 1
	//			return true
	return false
}

func (y boardState) illegal_ko_move(n int) bool {
	//TODO:
	// for i in move_history
	// 	if i = x   // check if move has been played before
	//		if (board state corresponding to move i)= gograph
	//			return true
	return false
}
