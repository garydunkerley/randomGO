package game

//graph() initializes a goGraph from the given topology.
func (a boardTop) initGraph() GoGraph {
	nodes := make(map[int]*node)
	for i := 0; i < a.nodeCount; i++ {
		newNode := initNode(i)
		nodes[i] = newNode
	}

	G := GoGraph{
		nodes:    nodes,
		boardTop: a,
	}
	return G
}

// initNode initializes an empty node with given integer id.
func initNode(i int) *node {
	var x node
	x.id = i
	x.color = 0
	return &x
}

// initBoardState initializes a board state with empty history,
// black to move, game not ongoing.
func initBoardState(boardTop) boardState {
	return boardState{GoGraph: boardTop.initGraph()}
}
