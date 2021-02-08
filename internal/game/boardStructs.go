package game

//node represents an intersection on the board.
//Every node should have a distinct integer id,
//ranging from 0 to (# nodes - 1)
type node struct {
	name      string
	id        int
	neighbors []*node
	color     int8       // color is 0 for empty, 1 for black, 2 for white
	inGroup   stoneGroup // is a part of a given stone group
}

//stoneGroup represents a contiguous string of stones.
type stoneGroup struct {
	stones    map[int]bool // what stones are in here? Need to make all empty if group dies
	liberties map[int]bool // what liberties does the group have
	color     int8         // color is 0 for empty, 1 for black, 2 for white
}

//GoGraph contains a map of node ids to node pointers, and embeds boardTop
type GoGraph struct {
	nodes map[int]*node
	boardTop
}

//boardTop stores the information to construct a game board.
type boardTop struct {
	edges        map[int][]int
	node_count   int            // We have node ids in the range [0, node_count - 1]
	coords       map[[2]int]int // optional, for boards with nice 2d representations
	coord_bounds []int          // optional, for example, [9,9] for a 9 by 9 board
}

//boardState stores a current game state, including history
type boardState struct {
	GoGraph
	white_to_move bool
	history
}

//history stores pointer-free metadata for previous moves in game.
//TODO: Zobrist hashing.
//TODO: reconstruct a board state from a history object.
type history struct {
	moves             []move         // The move history does not store passes.
	groups            [][]stoneGroup // Stores all groups. Does not update on passes.
	moveCount         int            // The move count is incremented by passes.
	consecutivePasses int            // Incremented by passes, reset to 0 by non-pass moves.
	whitePoints       int            // Accumulated points for white (komi, if any, plus captures)
	blackPoints       int            // Accumulated points for black (captures)
}

//move stores move information, such as the color and node id of a move.
type move struct {
	playerColor  int8   // 1 black, 2 white, 0 empty (eg handicap stones)
	id           int    // node ID of the move played
	capturesMade int    // Number of stones captured with this move.
	name         string // coordinates or other string identifier
	isPass       bool   // True if the player passed. Intentional redundancy: if node ID is -1, isPass is true.
}

// saveMoveAndGroups updates a history with the given move and group slice.
// It does not check legality of the move or group slice.
func (h *history) saveMoveAndGroups(m move, g []stoneGroup) {
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
	h.groups = append(h.groups, g)
	h.consecutivePasses = 0
}

// TODO: use this function to count how many stones are captured by the given move.
// Can assume the move is legal.
func (s boardState) getCaptures(m move) int {
	captureCount = 0
	return captureCount
}

// TODO: use this function to get a (post-move) list of stonegroups.
func (s boardState) getStoneGroups(m move) []stoneGroup {
	return nil
}
