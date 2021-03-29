package game

import (
	"math/rand"
	"strconv"
	"time"
)

// makeRandomBoard creates a board topology on n nodes such that
// any two nodes are connected via an edge with a fixed probability.
func makeNaiveRandomBoard(n int, prob float64) boardTop {
	var ourTopology boardTop

	edges := make(map[int][]int)

	ourTopology.nodeCount = n

	for i := 0; i < n; i++ {
		if i < n-1 {
			edges[i] = append(edges[i], i+1)

			for j := i + 2; j < n; j++ {
				rand.Seed(time.Now().UnixNano())
				v := rand.Float64()

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

// TODO change the names of nodes so that they are strings from which we can recover a coordinate plot.

func makeRandomBoard(n int) boardTop {
	var ourTopology boardTop

	edges := make(map[string][]string)
	existsAlready := make(map[string]bool)

	ourTopology.nodeCount = n
	layerCount := 0

	for i := 0; i < n; i++ {
		strconv.FormatFloat(0)
		strconv.FormatFloat()
	}

}
