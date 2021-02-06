package game

type node struct {
	name      string
	id        int
	neighbors []*node
	color     int8        // color is 0 for empty, 1 for black, 2 for white
	inGroup   *stoneGroup // is a part of a given stone group
}

type stoneGroup struct {
	//	id      int        // in case we want to refer to them?
	stones    map[*node]bool // what stones are in here? Need to make all empty if group dies
	liberties map[*node]bool // what liberties does the group have
	// https://softwareengineering.stackexchange.com/questions/177428/sets-data-structure-in-golang
}

type GoGraph struct {
	nodes        map[int]*node
	coords       map[[2]int]int
	coord_bounds []int //for example, [9,9] for a 9 by 9 board
}

type boardState struct {
	current_board GoGraph
	white_to_move bool
	move_history  []string
	board_history []GoGraph // to determine if move violates ko rule
	move_count    int
}
