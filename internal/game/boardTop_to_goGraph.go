package game

import (

	"strconv"

)
(a boardTop) func ConvertToGoGraph() goGraph {
	var ourBoard goGraph
	var newNode *node
	var nodeName string

	nodes := make(map[int]*node)

	for i:= 0; i < a.node_count; i++ {
		if boardTop.coords == nil {
			nodeName = strconv.Itoa(i) 
		} else {
			for j:= 0; j < a.node_count; j++ {
			nodeName = strconv()+ "-" +
		}
		newNode = initializeNode(i, nodeName)
		nodes[i]=newNode

	}

}

