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
	//	id      int        // in case we want to refer to them?
	stones    map[int]bool // what stones are in here? Need to make all empty if group dies
	liberties map[int]bool // what liberties does the group have
	color     int8         // color is 0 for empty, 1 for black, 2 for white
	// https://softwareengineering.stackexchange.com/questions/177428/sets-data-structure-in-golang
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
	moves     []move
	groups    [][]stoneGroup
	moveCount int
}

//move stores the color, node id, and name of any move.
type move struct {
	playerColor int8   //1 black, 2 white, 0 empty (eg handicap stones)
	id          int    //node ID of the move played
	name        string //coordinates or other string identifier
}
