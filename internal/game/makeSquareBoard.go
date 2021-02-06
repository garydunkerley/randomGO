package game

import (
	"fmt"
	"strconv"
)

func MakeSquareBoard(n int, m int) GoGraph {
	var name string
	goban_map := make(map[int]*node)
	coord_map := make(map[[2]int]int)
	var ourGraph GoGraph
	var newNode *node
	//initialize all nodes

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			//TODO: initialize a node named "i-j"
			name = strconv.Itoa(i) + "-" + strconv.Itoa(j)
			fmt.Println(name)

			ourSlice := [2]int{i, j}
			newNode = initializeNode(i*n+j, name)
			goban_map[i*n+j] = newNode
			coord_map[ourSlice] = i*n + j

		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j != 0 { // add left node if it exists
				goban_map[i*n+j].neighbors = append(goban_map[i*n+j].neighbors, goban_map[i*n+j-1])
			}
			if j != m-1 { // add right node if it exists
				goban_map[i*n+j].neighbors = append(goban_map[i*n+j].neighbors, goban_map[i*n+j+1])
			}
			if i != 0 { // add lower node if it exists
				goban_map[i*n+j].neighbors = append(goban_map[i*n+j].neighbors, goban_map[(i-1)*n+j])
			}
			if i != n-1 { // add upper node if it exists
				goban_map[i*n+j].neighbors = append(goban_map[i*n+j].neighbors, goban_map[(i+1)*n+j])
			}

		}
	}
	ourGraph.nodes = goban_map
	ourGraph.coords = coord_map
	ourGraph.coord_bounds = []int{n, m}
	return ourGraph
}
