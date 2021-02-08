package game

import (
	"math/rand"
	"time"
)

// makeRandomBoard creates a board topology on n nodes such that
// any two nodes are connected via an edge with a fixed probability.
func makeRandomBoard(n int, prob float64) boardTop {
	var ourTopology boardTop

	edges := make(map[int][]int)

	ourTopology.node_count = n

	for i = 0; i < n; i++ {
		if i < n-1 {
			edges[i] = append(edges[i], i+1)

			for j := i + 2; j < n; j++ {
				rand.Seed(time.Now().UnixNano())
				v = rand.Float64

				if v < prob {
					edges[i] = append(edges[i], j)
					edges[j] = append(edges[j], i)
				}
			}
		}
		if i > 0 {
			edges[i] = append(edges[i], i-1)
		}

	}

	ourTopology.edges = edges

	return ourTopology
}
