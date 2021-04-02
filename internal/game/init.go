package game

//graph() initializes a goGraph from the given topology.
//TODO add node connections
func (a *boardTop) initGraph() GoGraph {
	nodes := make(map[int]*node)
	for i := 0; i < a.nodeCount; i++ {
		newNode := initNode(i)
		nodes[i] = newNode
	}

	G := GoGraph{
		nodes:    nodes,
		stringOf: make(map[int]stoneString),
		boardTop: *a,
	}

	for i := 0; i < a.nodeCount; i++ {
		for _, z := range a.edges[i] {
			nodes[i].neighbors = append(nodes[i].neighbors, nodes[z])
		}
	}

	return G
}

func (gg GoGraph) getColorAssignments() map[int]int8 {

	colorMap := make(map[int]int8)
	for i := range gg.nodes {
		colorMap[i] = gg.nodes[i].color
	}

	return colorMap
}

/*
func (gameState boardState) initEbitenBoardInfo() EbitenBoardInfo {
	var boardInfo EbitenBoardInfo

	boardInfo.colorAssignments = gameState.GoGraph.getColorAssignments()
	boardInfo.boardTop = gameState.GoGraph.boardTop
	boardInfo.status = gameState.status
	boardInfo.internalPointer = &gameState

	return boardInfo
}
*/

// initNode initializes an empty node with given integer id.
func initNode(i int) *node {
	var x node
	x.id = i
	x.color = 0
	return &x
}

// initBoardState initializes a board state with empty history,
// black to move, game not ongoing.
// world's lowest priority TODO: fractional komi
func initBoardState(a boardTop, komi int) (X boardState) {
	X.GoGraph = a.initGraph()
	//X.status needs no initialization
	X.history.moves = make([]move, 0)
	X.history.allStoneStrings = make([]chromaticStrings, 0)
	X.history.whitePoints = komi
	// appends the simpleKo struct that indicates there are no ko points in the board.
	X.history.koHistory = append(X.history.koHistory, simpleKo{hasKo: false, koPoint: -1})
	return X
}

func initXEnclosedRegion(region map[*node]bool, color int8) xEnclosedRegion {
	var newXEnclosure xEnclosedRegion

	newXEnclosure.region = region
	newXEnclosure.boundaryColor = color

	return newXEnclosure
}
