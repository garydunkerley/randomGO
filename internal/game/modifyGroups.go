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
// to be subsumed into a larger stoneString given a play at nodeID of a given color
func (G GoGraph) getSubsumedStrings(nodeID int, color int8) []stoneString {
	mergeCandidates := make([]stoneString, 0)
	var friendlies []int
	for _, z := range G.nodes[nodeID].neighbors {
		if z.color == color {
			friendlies = append(friendlies, z.id)
		}
	}
	for _, z := range friendlies {
		mergeCandidates = append(mergeCandidates, G.stringOf[z])
	}
	return mergeCandidates
}

// getCapturedStrings gives a map encoding the stoneStrings
// that will be captured by a play at nodeID of a given color.
func (G GoGraph) getCapturedStrings(nodeID int, color int8) []stoneString {
	capturedStrings := make([]stoneString, 0)
	var potentialCaptives []int

	for _, z := range G.nodes[nodeID].neighbors {
		if z.color == getOppColor(color) {
			potentialCaptives = append(potentialCaptives, z.id)
		}
	}

	for _, z := range potentialCaptives {
		if G.countLiberties(G.stringOf[z]) == 1 {
			capturedStrings = append(capturedStrings, G.stringOf[z])
		}
	}
	return capturedStrings
}

// countCaptures returns the number of captured stones.
func countCaptures(capt []stoneString) int {
	sum := 0
	for _, string_ := range capt {
		sum += len(string_.stones)
	}
	return sum
}

// findCapturedLiberties takes each captured group and looks at its stones.
// If a given stone has a neighbor of the opposite color, then we
// append the nodeid of the captured stone to a list
// associated with the stoneString of the enemy neighbor.
// Wrapped by computeNewString.
// TODO finish deprecation
/*
func (G GoGraph) findCapturedLiberties(capt []stoneString) map[stoneString][]int {
	capturedLiberties := make([]stoneString)
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
*/

// mergeSubsumed uses a node id and subsumed strings
// to generate the new stoneString resulting from the given move.
// Wrapped by computeNewString.
func (G GoGraph) mergeSubsumed(id int,
	subsumed []stoneString) (new_ stoneString) {
	new_.stones = make(map[int]bool)
	new_.color = G.nodes[id].color
	new_.stones[id] = true
	if len(subsumed) == 0 { // Case: no friendly neighbors.
		return new_
	}
	for _, string_ := range subsumed {
		for stoneID := range string_.stones { // Add all stones
			new_.stones[stoneID] = true
		}
	}
	return new_
}

// computeNewString wraps mergeSubsumed.
// It computes the next string and the next map from ints to strings,
// by use of defer to avoid breaking anything.
func (G *GoGraph) computeNewString(
	subsumed []stoneString,
	m moveInput) (new_ stoneString) {
	defer func() {
		node := G.nodes[m.id]
		node.color = 0
		return
	}()
	node := G.nodes[m.id]
	node.color = m.playerColor //temporarily color a node
	new_ = G.mergeSubsumed(m.id, subsumed)
	return new_
}

//next stringOf table: update stones in new string and captured strings
/*
	//deprecate: don't need to return stringOf
	newStringOf = G.StringOf
	for stoneID := range new_.stones {
		newStringOf[stoneID] = new_
	}
	for str := range capt {
		for stoneID := range str.stones {
			delete(newStringOf, stoneID)
		}
	}
	return new_, newStringOf
*/

// computeNextChromaticStrings take the current strings and all the deltas
// and returns the strings for next turn, as a chromaticStrings object.
// It removes captured strings from the opponent's side,
// removes subsumed strings from your color.
// and adds the new string to your color.
func computeNextChromaticStrings(current chromaticStrings,
	capt []stoneString, subsumed []stoneString, new_ stoneString,
) chromaticStrings {
	for _, str := range capt {
		current.deleteStones(str)
	}
	for _, str := range subsumed {
		current.deleteStones(str)
	}
	current.addStones(new_)
	return current
}

// boardUpdate modifies the boardState by a single move, as follows.
// 0. Flip the x.whiteToMove bool.
// 1. Recolor the given node to the player's color.
// 2. Captured stones are removed.
// 3. The node groups are updated:
// 3a. x.GoGraph.stringOf[newnode] = newString
// 3b. x.GoGraph.stringOf[friendlyNodeInAdjacentGroup] = newString
// 3c. delete(x.stringOf, enemyNodeInCapturedGroup)
func (x *boardState) boardUpdate(m move, subsumed []stoneString,
	capt []stoneString, new_ stoneString) {
	x.whiteToMove = !x.whiteToMove
	newNode := x.nodes[m.id]
	newNode.color = m.playerColor // 1: add new stone
	x.stringOf[m.id] = new_       // 3a: update stone group

	for _, str := range capt {
		for id := range str.stones {
			captNode := x.nodes[id]
			captNode.color = 0     //2: remove captured stones
			delete(x.stringOf, id) //3c: update captured stone group
		}
	}

	for _, str := range subsumed {
		for id := range str.stones {
			x.stringOf[id] = new_ //3b: update subsumed stone group
		}
	}
	return
}
