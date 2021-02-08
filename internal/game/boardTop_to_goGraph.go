package game

import (
	"strconv"
)

func (a boardTop) ConvertToGoGraph() goGraph {
	var ourBoard goGraph
	var newNode *node
	var nodeName string

	ourNodes := make(map[int]*node)

	if a.coords == nil {
		// if the coordinate map is blank, initialize nodes
		// whose node name is the conversion of the node id into a string
		for i := 0; i < a.node_count; i++ {
			newNode = initializeNode(i, strconv.Itoa(i))
			ourNodes[i] = newNode
		}
	} else {
		// if the coordinate map is not blank, initialize
		// nodes whose names obey a grid convention
		for key := range boardTop.coords {
			nodeName = strconv.Itoa(key[0]) + "-" + strconv.Itoa(key[1])
			newNode = initializeNode(a.coords[key], nodeName)
			ourNodes[a.coords[key]] = newNode
		}

	}
	for i := 0; i < a.node_count; i++ {
		for z := range a.edges[i] {
			ourNodes[i].neighbors = append(ourNodes[i].neighbors, ourNodes[z])
		}
	}

	ourBoard.nodes = ourNodes
	ourBoard.boardTop = a

	return ourBoard

}
