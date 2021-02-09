package game

//graph() initializes a goGraph from the given topology.
//TODO add node connections
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
	for i := 0; i < a.nodeCount; i++ {
		for z := range a.edges[i] {
			nodes[i].neighbors = append(nodes[i].neighbors, nodes[z])
		}
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
func initBoardState(a boardTop) boardState {
	return boardState{GoGraph: a.initGraph()}
}
