package game

import "testing"

// Test a behavior where nonlocal corner capture occurs
func TestMove(t *testing.T) {
	state := initBoardState(makeSquareBoard(9, 9), 6)
	// Play a sequence of moves that gives bad behavior
	moveSequence := []int{52, 0, 53}
	for i, ID := range moveSequence {
		err := state.moveByID(ID)
		if err != nil { // log any error from this move
			t.Logf("Move %v: %v\n Error: %v", i, ID, err)
		}
	}
	// Fail if last move made a capture
	lastMove := state.moves[len(state.moves)-1]
	if lastMove.capturesMade != 0 {
		t.Logf("Captures last move: %v\n", lastMove.capturesMade)
		t.Fatal(state) // log end state
	}
}
