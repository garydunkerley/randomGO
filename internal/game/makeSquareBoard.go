package game

import (
	"fmt"
	"strconv"
)

func MakeSquareBoard(n int, m int) boardTop {

	var ourTopology boardTop

	edges := make(map[int][]int)
	coords := make(map[[2]int]int)

	ourTopology.nodeCount = n * m
	ourTopology.coordBounds = []int{m, n}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {

			coords[[]int{i, j}] = i*n + j

			if j != 0 { // add left node if it exists
				edges[(i*n)+j] = append(edges[i*n+j], i*n+j-1)
			}
			if j != m-1 { // add right node if it exists
				edges[i*n+j] = append(edges[i*n+j], i*n+j+1)
			}
			if i != 0 { // add lower node if it exists
				edges[i*n+j] = append(edges[i*n+j], (i-1)*n+j)
			}
			if i != n-1 { // add upper node if it exists
				edges[i*n+j] = append(edges[i*n+j], (i+1)*n+j)
			}
		}
	}

	ourTopology.edges = edges
	ourTopology.coords = coords
	return ourTopology
}
