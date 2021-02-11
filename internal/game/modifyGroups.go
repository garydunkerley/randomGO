package game

//TODO: Jobs for preparing a move.
//	   1. From move input: compute captured/subsumed/new strings
//DONE 2. Count captured stones. [used to track points]
//DONE 3. Using this information, compute the set of all strings after the move.
//			(This is tracked in history)
//Essentially done: playing the move.

// countCaptures will count the number of nodes in the given map.
// (used for scoring purposes)

// getOppColor(2)=1, getOppColor(1)=2
func getOppColor(color int8) int8 {
	return 3 - color
}

// getSubsumedStrings gives us a map encoding the stoneStrings that need
// to be subsumed into a larger stoneString given a play at nodeId of a given color

func (G GoGraph) getSubsumedStrings(nodeId int, color int8) map[stoneString]bool {
	stringsToBeMerged := make(map[stoneString]bool)
	var friendlies []int
	for _, z := range G.nodes[nodeId].neighbors {
		if z.color == color {
			friendlies = append(friendlies, z.id)
		}
	}
	for _, z := range friendlies {
		stringsToBeMerged[G.stringOf[z]] = true
	}
	return stringsToBeMerged
}

// getCapturedStrings gives a map encoding the stoneStrings
// that will be captured by a play at NodeId of a given color.
func (G GoGraph) getCapturedStrings(nodeId int, color int8) map[stoneString]bool {
	capturedStrings := make(map[stoneString]bool)
	var potentialCaptives []int

	for _, z := range G.nodes[nodeId].neighbors {
		if z.color == getOppColor(color) {
			potentialCaptives = append(potentialCaptives, z.id)
		}
	}

	for _, z := range potentialCaptives {
		if len(G.stringOf[z].liberties) == 1 {
			capturedStrings[G.stringOf[z]] = true
		}
	}
	return capturedStrings
}

// countCaptures returns the number of captured stones.
func countCaptures(capt map[stoneString]bool) int {
	sum := 0
	for string_ := range capt {
		sum += len(string_.stones)
	}
	return sum
}

// findCapturedLiberties takes each captured group and looks at its stones.
// If a given stone has a neighbor of the opposite color, then we
// append the nodeid of the captured stone to a list
// associated with the stoneString of the enemy neighbor.
// Wrapped by getNewStringData.
func (G GoGraph) findCapturedLiberties(capt map[stoneString]bool) map[stoneString][]int {
	capturedLiberties := make(map[stoneString][]int)
	var captors []int
	// for each captured group
	for key := range capt { // and for each stone in this captured group
		for id := range key.stones { // look at its enemy neighbors
			captive := G.nodes[id]
			captors = getOppColorNeighbors(captive)
			// if there are any, append the nodeid of the captive
			// to the list of new liberties for the stoneString of the enemy neighbor
			for _, cid := range captors {
				capturedLiberties[G.stringOf[cid]] = append(
					capturedLiberties[G.stringOf[cid]], captive.id)
			}
		}
	}
	return capturedLiberties
}

// findNewStringOf returns an updated map[int]stoneString,
// intended to populate the next GoGraph.stringOf field.
// It expects the output of findCapturedLiberties, which is in the form
// map[string which needs added liberties] == liberties to add
// Wrapped by getNewStringData.
func (G GoGraph) findNewStringOf(capturedLiberties map[stoneString][]int) (
	newStringOf map[int]stoneString) {
	newStringOf = G.stringOf
	for str, libs := range capturedLiberties {
		newLiberties := str.liberties
		for lib := range libs {
			newLiberties[lib] = true
		}
		for stoneId := range str.stones { //update group for each node in this group
			newStringOf[stoneId].liberties = newLiberties
		}
	}
	return newStringOf
}

// computeNewString uses a node id, subsumed strings, and captured strings
// to generate the new stoneString resulting from the given move.
// Note that it does not update for captures.
// Wrapped by getNewStringData.
func (G GoGraph) computeNewString(id int,
	subsumed map[stoneString]bool) (new_ stoneString) {
	new_.stones, new_.liberties = make(map[int]bool), make(map[int]bool)
	new_.color = G.nodes[id].color
	new_.stones[id] = true
	if len(subsumed) == 0 { // Case: no friendly neighbors.
		for libID := range getLiberties(G.nodes[id]) {
			new_.liberties[libID] = true
		}
		return new_
	}
	for string_, _ := range subsumed {
		for stoneId, _ := range string_.stones { // Add all stones
			new_.stones[stoneId] = true
		}
		for stoneId, _ := range string_.liberties {
			// Add all liberties, will remove the candidate move as a liberty later
			new_.liberties[stoneId] = true
			// NOTE: This does not account for liberties obtained after capturing.
		}
	}
	delete(new_.liberties, id) //remove itself as a liberty
	return new_
}

// getNewStringData wraps other string methods.
// It computes the next string and the next map from ints to strings,
// by use of defer to avoid breaking anything.
func (G GoGraph) getNewStringData(
	capt map[stoneString]bool,
	subsumed map[stoneString]bool,
	m moveInput) (newString stoneString, newStringOf map[int]stoneString) {
	defer func() {
		G.nodes[m.id] = 0
	}()
	G.nodes[m.id] = m.playerColor //temporarily color a node
	newStringOf = G.findNewStringOf(G.findCapturedLiberties(capt))
	newString = G.computeNewString(m.id, subsumed)
	return newString, newStringOf
}

// computeNextStrings take the current strings and all the deltas
// and returns the strings for next turn, as a chromaticStrings object.
// It removes captured strings from the opponent's side,
// removes subsumed strings from your color.
// and adds the new string to your color.
func computeNextStrings(current chromaticStrings,
	capt map[stoneString]bool,
	subsumed map[stoneString]bool,
	new_ stoneString) chromaticStrings {
	for str := range append(capt, subsumed) {
		current.deleteStones(str)
	}
	current.addStones(new_)
	return current
}

// boardUpdate modifies the boardState by a single move, as follows.
// 1. Recolor the given node to the player's color.
// 2. Captured stones are removed.
// 3. The node groups are updated:
// 3a. x.GoGraph.stringOf[newnode] = newString
// 3b. x.GoGraph.stringOf[friendlyNodeInAdjacentGroup] = newString
// 3c. delete(x.stringOf, enemyNodeInCapturedGroup)
func (x *boardState) boardUpdate(m move,
	subsumed map[stoneString]bool,
	capt map[stoneString]bool,
	new_ stoneString) {
	x.nodes[m.id] = m.playerColor // 1: add new stone
	x.stringOf[m.id] = new_       // 3a: update stone group

	for string_ := range capt {
		for id := range string_.stones {
			x.nodes[id].color = 0  //2: remove captured stones
			delete(x.stringOf, id) //3c: update captured stone group
		}
	}

	for string_ := range subsumed {
		for id := range string_.stones {
			x.stringOf[id] = new_ //3b: update subsumed stone group
		}
	}
	return
}
