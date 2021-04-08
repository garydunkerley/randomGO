package game

// makeSquareBoard returns a rectangular board topology.
func makeSquareBoard(n int, m int) BoardTop {

	var ourTopology BoardTop

	edges := make(map[int][]int)
	coords := make(map[[2]int]int)

	ourTopology.nodeCount = n * m
	ourTopology.coordBounds = []int{m, n}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			ID := i*m + j
			coords[[2]int{i, j}] = ID

			// The board is traversed top-to-bottom, left-to-right,
			// starting from the top left corner (0-0).
			// Vertical shifts correspond to ID shifts by 1.
			// Horizontal shifts correspond to ID shifts by m.

			if i != 0 { // add left node (if it exists)
				edges[ID] = append(edges[ID], ID-m)
			}
			if i != n-1 { // add right node
				edges[ID] = append(edges[ID], ID+m)
			}
			if j != 0 { // add upper node
				edges[ID] = append(edges[ID], ID-1)
			}
			if j != m-1 { // add lower node
				edges[ID] = append(edges[ID], ID+1)
			}
		}
	}

	ourTopology.edges = edges
	ourTopology.coords = coords
	return ourTopology
}
