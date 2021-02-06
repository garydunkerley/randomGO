package game

import (
	"math/rand"
	"strconv"
	"time"
)

func makeRandomBoard(n int, prob float64) GoGraph {
	var ourGraph GoGraph
	var ourMap = make(map[int]*node)

	var newNode *node
	var eligibleNode *node
	var v float64

	for i := 0; i < n; i++ {
		newNode = initializeNode(i, strconv.Itoa(i))

		ourMap[i] = newNode
		// add the new node to the goGraph struct

		if i > 0 {

			for j := 0; j < i-1; j++ {
				eligibleNode = ourMap[j]

				rand.Seed(time.Now().UnixNano())
				// seed random number generator at each step of loop

				v = rand.Float64()

				if v < prob {
					newNode.neighbors = append(newNode.neighbors, eligibleNode)
					eligibleNode.neighbors = append(eligibleNode.neighbors, newNode)
					// if dice roll is passed, make j a neighbor of newNode
				}
			}
			eligibleNode = ourMap[i-1]
			newNode.neighbors = append(newNode.neighbors, eligibleNode)
			eligibleNode.neighbors = append(eligibleNode.neighbors, newNode)

			// we make it so the (i-1)th node is always a neighbor of the ith node
		}

	}
	ourGraph.nodes = ourMap
	// for i := 0; i < n; i++ {            // uncomment to get a readout of the initialized board
	//	fmt.Println("[", ourGraph.nodes[i].name, ", ", ourGraph.nodes[i].color, ", ", ourGraph.nodes[i].neighbors, ", ", ourGraph.nodes[i].inGroup, "]")
	// }
	return ourGraph
}
