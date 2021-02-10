package game

// node represents an intersection on the board.
// Every node should have a distinct integer id,
// ranging from 0 to (# nodes - 1)
type node struct {
	name      string // deprecating unless there is a reason to care
	id        int
	neighbors []*node
	color     int8        // color is 0 for empty, 1 for black, 2 for white
	group     stoneString // deprecating: is a part of a given stone group
}

// stoneString represents a contiguous string of stones.
type stoneString struct {
	stones    map[int]bool // what stones are in here? Need to make all empty if group dies
	liberties map[int]bool // deprecating: we have decided just to compute liberties as needed.
	color     int8         // color is 0 for empty, 1 for black, 2 for white
}

// chromaticStrings represents two sets of stoneStrings, one per color
// Usage: ensure that the colors are copacetic.
//   Each key G of chromaticStrings.black must satisfy G.color == 1
//   (resp. chromaticStrings.white, G.color == 2)
type chromaticStrings struct {
	black map[stoneString]bool
	white map[stoneString]bool
}

// boardTop stores the information to construct a game board.
type boardTop struct {
	edges       map[int][]int
	nodeCount   int            // We have node ids in the range [0, nodeCount - 1]
	coords      map[[2]int]int // optional, for boards with nice 2d representations
	coordBounds []int          // optional, for example, [9,9] for a 9 by 9 board
}

// GoGraph holds board topology and maps node IDs to their strings and *nodes.
type GoGraph struct {
	nodes    map[int]*node
	stringOf map[int]stoneString // Maps node ID to its containing string
	boardTop
}

// boardState stores a current game state, including history
type boardState struct {
	GoGraph
	status
	history
}

// status tracks current turn and if the game is currently active
type status struct {
	whiteToMove bool
	ongoing     bool // true if game has not finished
}

// history stores pointer-free metadata for previous moves in game.
// TODO: Zobrist hashing.
// TODO: reconstruct a board state from a history object.
type history struct {
	moves             []move             // The move history does not store passes.
	allStoneStrings   []chromaticStrings // Stores all groups. Does not update on passes.
	moveCount         int                // The move count is incremented by passes.
	consecutivePasses int                // Incremented by passes, reset to 0 by non-pass moves.
	whitePoints       int                // Accumulated points for white (komi, if any, plus captures)
	blackPoints       int                // Accumulated points for black (captures)
}

// move is a moveInput and the associated number of captures made with the move
type move struct {
	moveInput
	capturesMade int
}

// moveInput stores the color and id of a valid move, and a bool for passing (-1)
type moveInput struct {
	playerColor int8 // 1 black, 2 white, 0 empty (eg handicap stones)
	id          int  // node ID of the move played
	isPass      bool // True if the player passed. Intentional redundancy: if node ID is -1, isPass is true.
}

// addMoveAndStrings updates a history with the given move and computed group slice.
// It is assumed that the move is a legal move, and it is not a pass.
// It does not check legality of the move. Called by playMoveInput.
func (h *history) addMoveAndStrings(m move, C chromaticStrings) {
	h.moveCount++
	if m.isPass {
		h.consecutivePasses++
		return
	}
	if m.capturesMade != 0 {
		if m.playerColor == 1 {
			h.blackPoints += m.capturesMade
		} else if m.playerColor == 2 {
			h.whitePoints += m.capturesMade
		} else {
			panic("Attempted to save an invalid move to history (colorless capture).")
		}
	}
	h.moves = append(h.moves, m) // don't update moves or groups on pass
	h.allStoneStrings = append(h.allStoneStrings, C)
	h.consecutivePasses = 0
}

// playMoveInput will play a moveInput, changing the history and then mutating
// board state.
// NOTE: If non-nil error is returned, no changes should be made to global state.
func (s *boardState) playMoveInput(input moveInput) error {
	// Check legality
	err := s.checkLegalMove(input.id, input.playerColor)
	if err != nil {
		return err
	}

	// Check pass
	if input.isPass {
		s.history.consecutivePasses++
		s.history.moveCount++
		if s.consecutivePasses >= 2 {
			s.ongoing = false
		}
		return nil
	}

	nodeID, color := input.id, input.playerColor
	subsumed := s.getSubsumedStrings(nodeID, color)
	capt := s.getCapturedStrings(nodeID, color)
	//Note that the liberties of these strings do not account for captures yet.
	//It's just the old strings, before the move is played.
	newString, newStringOf := s.getNewStringData(capt, subsumed, input)

	m := move{
		moveInput:    input,
		capturesMade: countCaptures(capt),
	}

	// Store the info as a chromaticStrings object
	var last chromaticStrings
	if L := len(s.allStoneStrings); L != 0 {
		lastChromaticStrings := s.allStoneStrings[L-1]
	}

	// For computeNextStrings to be accurate, newString should have accurate liberties.
	// capt and subsumed need only have correct stones.
	// FIXME: There is probably going to be some bug with liberties here.
	// Be sure to test cases where a single stone kills a string, or when a group kills a string.
	next := computeNextStrings(last, capt, subsumed, newString, m.playerColor)

	s.history.addMoveAndStrings(m, next)

	s.boardUpdate(m, subsumed, capt, newString)
	return nil
}
