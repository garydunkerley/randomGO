package game

import (
	"fmt"
)

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
	var isNew bool

	var mergeCandidates []stoneString

	mergeStringReps := make(map[int]bool)

	for _, z := range G.nodes[nodeID].neighbors {
		if z.color == color {
			isNew = true
			for i := range mergeStringReps {
				if mapKeysEqual(G.stringOf[z.id].stones, G.stringOf[i].stones) {
					isNew = false
				}

			}
			if isNew {
				mergeStringReps[z.id] = true
			}
		}
	}
	fmt.Println("Debug: Our merge representatives are: ")
	for z := range mergeStringReps {
		fmt.Println(z)
		mergeCandidates = append(mergeCandidates, G.stringOf[z])
	}

	return mergeCandidates
}

// getCapturedStrings gives a map encoding the stoneStrings
// that will be captured by a play at nodeID of a given color.
// Expects to be called before the move in question is played.
func (G GoGraph) getCapturedStrings(nodeID int, color int8) []stoneString {

	var isNew bool
	var capturedStrings []stoneString
	captiveStringReps := make(map[int]bool)

	for _, z := range G.nodes[nodeID].neighbors {
		if z.color == getOppColor(color) {
			isNew = true
			for i := range captiveStringReps {
				if mapKeysEqual(G.stringOf[z.id].stones, G.stringOf[i].stones) {
					isNew = false
				}

			}
			if isNew {
				captiveStringReps[z.id] = true
			}

		}
	}

	for z := range captiveStringReps {
		myStoneString := G.stringOf[z]
		if len(myStoneString.liberties) == 1 {
			for s := range myStoneString.liberties {
				if s == nodeID {
					capturedStrings = append(capturedStrings, myStoneString)
				}
			}
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

// assignAllLiberties looks at all stoneStrings, determines their liberties, and then assigns
// them to the stoneString
func (cs *chromaticStrings) assignAllLiberties(gg GoGraph) {

	blackStrings := cs.black
	whiteStrings := cs.white

	for _, b := range blackStrings {
		accountedFor := make(map[int]bool)
		newLibs := make(map[int]bool)

		for s := range b.stones {
			for n := range gg.nodes[s].neighbors {
				if accountedFor[n] == false && gg.nodes[n].color == 1 {
					newLibs[n] = true
					accountedFor[n] = true
				} else {
					accountedFor[n] = true
				}
			}
		}
		b.liberties = newLibs
	}

	for _, w := range whiteStrings {
		accountedFor := make(map[int]bool)
		newLibs := make(map[int]bool)

		for s := range w.stones {
			for n := range gg.nodes[s].neighbors {
				if accountedFor[n] == false && gg.nodes[n].color == 1 {
					newLibs[n] = true
					accountedFor[n] = true
				} else {
					accountedFor[n] = true
				}
			}
		}
		w.liberties = newLibs
	}

	return

}

// getKoPoints takes the outputs of computeNewString and getCapturedStrings to determine whether a given play produces a ko point, which will prevent the player from playing in that position next turn.
func (s *boardState) setKoPoint(inputID int, newString stoneString, capt []stoneString) {
	var setKo simpleKo
	if countCaptures(capt) == 1 {
		if len(newString.stones) == 1 {
			if s.GoGraph.countLiberties(newString) == 1 {
				setKo.hasKo = true
				setKo.koPoint = getLiberties(s.GoGraph.nodes[inputID])[0]
			} else {
				setKo.hasKo = false
				setKo.koPoint = -1
			}
		} else {
			setKo.hasKo = false
			setKo.koPoint = -1
		}
	} else {
		setKo.hasKo = false
		setKo.koPoint = -1
	}
	s.history.koHistory = append(s.history.koHistory, setKo)
	return
}

// mergeSubsumed uses a node id and subsumed strings
// to generate the new stoneString resulting from the given move.
// Wrapped by computeNewString.
func (G GoGraph) mergeSubsumed(id int,
	subsumed []stoneString) (new_ stoneString) {

	new_.stones = make(map[int]bool)
	new_.liberties = make(map[int]bool)
	new_.color = G.nodes[id].color
	new_.stones[id] = true
	if len(subsumed) == 0 { // Case: no friendly neighbors.
		return new_
	}
	for _, string_ := range subsumed {
		for stoneID := range string_.stones { // Add all stones
			new_.stones[stoneID] = true
		}
		for n := range string_.liberties {
			if n != id {
				new_.liberties[n] = true
			}
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
	newLibs := make(map[int]bool)

	node := G.nodes[m.id]
	node.color = m.playerColor //temporarily color a node
	new_ = G.mergeSubsumed(m.id, subsumed)
	newLibs = new_.liberties
	for _, n := range G.nodes[m.id].neighbors {
		if G.nodes[n.id].color == 0 {
			newLibs[n.id] = true

			// if you find an enemy node, remove the
			// player's move input from its list of liberties
		} else if G.nodes[n.id].color != m.playerColor && G.nodes[n.id].color != 0 {
			removeLib := make(map[int]bool)
			oldEnemyString := G.stringOf[n.id]
			removeLib = oldEnemyString.liberties
			delete(removeLib, m.id)

			oldEnemyString.liberties = removeLib
		}

	}
	new_.liberties = newLibs
	return new_

}

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
