package game

//TODO: Jobs for preparing a move.
//	   1. From move input: compute captured/subsumed/new strings
//DONE 2. Count captured stones. [used to track points]
//DONE 3. Using this information, compute the set of all strings after the move.
//			(This is tracked in history)

//TODO: Jobs for playing a move.
// Data needed:
// * The sets of all strings for each color, before the move was played.
//		This is accessible via the last history entry.
// * Captured and subsumed strings, as map[stoneString]bool
// * The new string, as stoneString
// 0. Copy the starting string information (chromaticStrings, from history)
// 1. Remove all captured strings from the opponent's string set.
// 2. Remove all subsumed strings from your own string set.
// 3. Iterate over all nodes in the captured strings, set their color to 0.

// countCaptures will count the number of nodes in the given map.
// (used for scoring purposes)
func countCaptures(capt map[stoneString]bool) int {
	sum := 0
	for string_, _ := range capt {
		sum += len(string_.stones)
	}
	return sum
}

// computeNewString takes the subsumed strings and candidate move id,
// then generates a new stoneString with the union of all of them
func computeNewString(
	id int,
	subsumed map[stoneString]bool,
) (new_ stoneString) {
	new_.stones, new_.liberties = make(map[int]bool), make(map[int]bool)
	new_.color = id.color
	new_.stones[id] = true
	if len(subsumed) == 0 { // Case: no friendly neighbors.
		//TODO: new_.liberties should be the "get liberties" function for a single node
		return new_
	}
	for string_, _ := range subsumed {
		for stoneId, _ := range string_.stones { // Add all stones
			new_.stones[stoneId] = true
		}
		for stoneId, _ := range string_.liberties {
			// Add all liberties, will remove the candidate move as a liberty later
			new_.liberties[stoneId] = true
			// TODO: This does not account for liberties obtained after capturing.
			// Perhaps add the captured strings as an argument? Not sure
		}
	}
	delete(new_.liberties, id) //remove itself as a liberty
	return new_
}

// computeNextStrings take the current strings and all the deltas
// and returns the strings for next turn, as a chromaticStrings object.
// It removes captured strings from the opponent's side,
// removes subsumed strings from your color.
// and adds the new string to your color.
func computeNextStrings(current chromaticStrings,
	capt map[stoneString]bool,
	subsumed map[stoneString]bool,
	new_ stoneString,
	moveColor int8,
) (next chromaticStrings) {
	var ownStrings map[stoneString]bool
	var oppStrings map[stoneString]bool
	if moveColor == 1 {
		ownStrings = current.black
		oppStrings = current.white
	} else {
		ownStrings = current.white
		oppStrings = current.black
	}
	for string_, _ := range capt {
		delete(oppStrings, string_)
	}
	for string_, _ := range subsumed {
		delete(ownStrings, string_)
	}
	ownStrings[new_] = true
	// now give them the white/black labels
	if moveColor == 1 {
		next.black = ownStrings
		next.white = oppStrings
	} else {
		next.black = oppStrings
		next.white = ownStrings
	}
	return next
}

// boardUpdate modifies the boardState by a single move, as follows.
// 1. Recolor the given node to the player's color.
// 2. Captured stones are removed.
// 3. The node groups are updated:
// 3a. x.GoGraph.stringOf[newnode] = newString
// 3b. x.GoGraph.stringOf[friendlyNodeInAdjacentGroup] = newString
// 3c. delete(x.stringOf, enemyNodeInCapturedGroup)
func (x *boardState) boardUpdate(
	m move,
	subsumed map[stoneString]bool,
	capt map[stoneString]bool,
	new_ stoneString,
) {
	x.nodes[m.id] = m.color // 1: add new stone
	x.stringOf[m.id] = new_ // 3a: update stone group

	for string_, _ := range capt {
		for id, _ := range string_.stones {
			x.nodes[id].color = 0  //2: remove captured stones
			delete(x.stringOf, id) //3c: update captured stone group
		}
	}
	for string_, _ := range subsumed {
		for id, _ := range string_.stones {
			x.stringOf[id] = new_ //3b: update subsumed stone group
		}
	}
	return
}
