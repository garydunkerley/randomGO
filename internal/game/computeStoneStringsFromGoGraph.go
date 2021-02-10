package game

func initStoneString(someStones map[int]bool) stoneString {
	var myStoneString stoneString
	myStoneString.stones = someStones

	for z := range someStones {
		G.stringOf[z] = myStoneString
	}
}

func (G GoGraph) computeStoneStrings() {

	anEmptyMap := make(map[int]bool)

	for node := range G.nodes {
		// once you find a colored node
		if node.color != 0 {
			// if it does not already belong to a stoneString
			if G.stringOf[node.id] == nil {

				nodesInStoneString := make(map[int]bool)

				nextLevelDown := make(map[int]bool)
				nextnextLevelDown := make(map[int]bool)

				// establish that the first node will be
				// in the stone string
				nodesInStoneString[node.id] = true

				// collect all friendly neighbors of node
				// to be designated as belonging to the
				// next level of the tree
				for friendlyID := range sameColorNeighbors(node) {

					nextLevelDown[friendlyId] = true
				}

				// while the next level of nodes in
				// the spanning tree for the stoneString
				// is non-empty
				for len(nextLevelDown) > 0 {
					for x := range nextLevelDown {

						nodesInStoneString[x] = true

						// append the neighbors of x to a list
						// corresponding to all stones the
						// next layer after the current one
						for y := range sameColorNeighbors(G.nodes[x]) {

							// if a given same color neighbor of x
							// is not already in a group
							// indicate it belongs to the next layer to be checked
							if nodesInStoneString[y] == false {
								nextnextLevelDown[y] = true
							}
						}

					}
					// once you've documented all NodeIds in the current level,
					// move on to the next one
					nextLevelDown = nextnextLevelDown
				}

				// creates a stoneString whose constituents are determined
				// by the map and update these nodes to indicate that they belong to
				// this stoneString
				initStoneString(nodesInStoneString)
			}
		}
	}
}
