package game

import (
	"fmt"
)

func (x boardState) make_move(y int) {
	var stoneColor int8
	if x.whiteToMove {
		stoneColor = 2
	} else {
		stoneColor = 1
	}
	node, ok := x.current_board.nodes[y]
	if !ok {
		error_message := fmt.Sprintf("Node %d not found!", y)
		panic(error_message)
	}
	node.color = stoneColor
	modify_stoneGroups(x.current_board.nodes[y])
	x.remove_dead(x.current_board.nodes[y])
	modify_stoneGroups(x.current_board.nodes[y])
	// I run modify_stoneGroups again to append any liberties resulting from stone groups dying
	// TODO clean up this inefficiency
}

func contains_val(slice []*node, val *node) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
