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

// Test liberties for an edge move
func TestEdgeLiberties(t *testing.T) {
	state := initBoardState(makeSquareBoard(9, 9), 6)
	ID := 1
	err := state.moveByID(ID)
	if err != nil { // log any error from this move
		t.Logf("Move %v: %v\n Error: %v", 0, ID, err)
	}

	libCount := state.countLiberties(state.stringOf[ID])
	if libCount != 3 {
		t.Logf("Edge test: expected 3 liberties, have %v", libCount)
		t.Fatal(state)
	}
}

// Test liberties for a corner move
func TestCornerLiberties(t *testing.T) {
	state := initBoardState(makeSquareBoard(9, 9), 6)
	ID := 0
	err := state.moveByID(ID)
	if err != nil { // log any error from this move
		t.Logf("Move %v: %v\n Error: %v", 0, ID, err)
	}

	libCount := state.countLiberties(state.stringOf[ID])
	if libCount != 2 {
		t.Logf("Corner test: expected 2 liberties, have %v", libCount)
		t.Fatal(state)
	}
}

// Test liberties for a central move
func TestCenterLiberties(t *testing.T) {
	state := initBoardState(makeSquareBoard(9, 9), 6)
	ID := 40
	err := state.moveByID(ID)
	if err != nil { // log any error from this move
		t.Logf("Move %v: %v\n Error: %v", 0, ID, err)
	}

	libCount := state.countLiberties(state.stringOf[ID])
	if libCount != 4 {
		t.Logf("Corner test: expected 4 liberties, have %v", libCount)
		t.Fatal(state)
	}
}

func TestPlayer1(t *testing.T) {
	state := initBoardState(makeSquareBoard(9, 9), 6)
	if state.whiteToMove {
		t.Fatal("whiteToMove true on initial move")
	}
}

func TestSwitchOnPass(t *testing.T) {
	state := initBoardState(makeSquareBoard(9, 9), 6)
	previousPlayer := state.whiteToMove
	state.moveByID(-1)
	if state.whiteToMove == previousPlayer {
		t.Logf("Player not switching off on pass.")
	}
}
