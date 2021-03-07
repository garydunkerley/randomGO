package game

import (
	"fmt"
)

/*
NOTE
A scoring algorithm needs to be able to locate x-enclosed regions and then determine whether the blocks that comprise the boundary of these regions are healthy or not.

We don't need to detect unconditional life necessarily. I think it will be sufficient to show that all constituent blocks of a group enclose regions that are healthy for them.

If a group W of white stones sits inside of a region B of blakc stones and does not satisfy this "two healthy regions" condition, then we regard the stones as dead for scoring purposes.


*/

func getEmptyNodes(gg GoGraph) map[*node]bool {
	emptyNodes := make(map[*node]bool)
	for _, z := range gg.nodes {
		if z.color == 0 {
			emptyNodes[z] = true
		}
	}
	return emptyNodes
}

// getEmptyRegion takes a node, verifies that it's empty and then uses a bucket fill algorithm to determine its empty component.
func getEmptyRegion(n *node) map[*node]bool {
	accountedFor := make(map[*node]bool)
	stillSearching := true
	counter := 0

	if n.color == 0 {
		accountedFor[n] = true
		for stillSearching && counter < 1000 {
			newEmpties := make(map[*node]bool)
			for z := range accountedFor {
				for _, q := range z.neighbors {
					if q.color == 0 && accountedFor[q] == false {
						accountedFor[q] = true
						newEmpties[q] = true
					}
				}
			}
			if len(newEmpties) < 1 {
				stillSearching = false
			}
			counter += 1

		}
	} else {
		fmt.Println("Error: Node is not empty")
	}
	return accountedFor
}

// isXEnclosed takes an empty region and then determines whether the
// boundary is uniformly colored or mixed
func getBoundaryColor(emptyRegion map[*node]bool) int8 {

	var color int8
	var oldColor int8

	sameColor := true

	for z := range emptyRegion {
		for _, n := range z.neighbors {
			if n.color != 0 {
				oldColor = color
				if n.color == 1 {
					color = 1
				} else {
					color = 2
				}

				if oldColor != color {
					sameColor = false
					break
				}
			}
		}
		if sameColor == false {
			return 0
		}

	}
	return color

}

// getAllXEnclosedRegions takes every empty node in the GoGraph,
// locates the empty connected component of this node and then determines whether
// this component is x-enclosed for a color x (colored 0 means
// the region has a boundary of mized color). The function then outputs a map
// that will allow us to check if an empty node belongs to an x-enclowed region.
func getAllXEnclosedRegions(gg GoGraph) map[*node]xEnclosedRegion {

	myMap := make(map[*node]xEnclosedRegion)

	// this just makes sure we're not double counting nodes whose regions
	// we've already found
	accountedFor := make(map[*node]bool)

	for _, z := range gg.nodes {
		if z.color == 0 && accountedFor[z] == false {

			emptyRegion := getEmptyRegion(z)
			color := getBoundaryColor(emptyRegion)
			newXEnclosure := initXEnclosedRegion(emptyRegion, color)

			for q := range emptyRegion {

				accountedFor[q] = true
				myMap[q] = newXEnclosure
			}

		}
	}
	return myMap
}

// inMexicanStandOff determines whether a stoneString is occupying a particular type of seki
// with another stoneString where both stoneStrings have two liberties and share at least one.
// In this situtation, whoever plays in one of the shared liberties necessarily loses their group,
// and so the best choice for both players is to allow everything to stand.
// If two stonestrings are found to be in a Mexican standoff, we then check nearby stone strings to see
// if they are alo implicated
func (gg GoGraph) inMexicanStandOff(myString *stoneString) (bool, []*stoneString) {
	var duelists []*stoneString

	// if my stoneString has two liberties
	if len(myString.liberties) == 2 {

		accountedFor := make(map[*node]bool)
		// then for each stone in the string, look a neighbor of the opposite color and check to see if it also has two liberties and at least one is shared

		for i := range myString.stones {
			for _, n := range gg.nodes[i].neighbors {
				if n.color != myString.color && n.color != 0 && !accountedFor[n] {

					a := gg.stringOf[n.id]
					potentialDuelist := &a
					if len(potentialDuelist.liberties) == 2 {
						for z := range potentialDuelist.liberties {
							if myString.liberties[z] {
								duelists = append(duelists, myString)
								duelists = append(duelists, potentialDuelist)
								return true, duelists
							}
						}
					}
					accountedFor[n] = true
				}

			}

		}
	}
	return false, duelists
}

// getSimpleDeadandSeki takes a GoGraph and a map assigning emptynodes to their xEnclosedRegions and determines which
// stoneStrings can be immediately ruled as dead and those occupying a type of seki we have termed
// a "Mexican standoff". Determining more sophistacted dead and seki (those strings whose status depends on another string)
// is dealt with in the getDependents function
func (cs chromaticStrings) getSimpleDeadandSeki(gg GoGraph, myMap map[*node]xEnclosedRegion) {

	for _, a := range cs.black {
		myString := &a
		accountedFor := make(map[*node]bool)
		xEnclosedRegions := 0

		if myString.state == "" {
			for l := range myString.liberties {
				if accountedFor[gg.nodes[l]] == false {
					newXEnclosure := myMap[gg.nodes[l]]
					if newXEnclosure.boundaryColor == 1 {
						xEnclosedRegions += 1
					}
					for n := range newXEnclosure.region {
						accountedFor[n] = true
					}

				}
			}
			if xEnclosedRegions < 2 {
				isSeki, duelists := gg.inMexicanStandOff(myString)
				if isSeki {
					for _, i := range duelists {
						i.state = "seki"
					}
				} else {
					myString.state = "dead"
				}
			}
		}
	}

	for _, a := range cs.white {
		myString := &a
		accountedFor := make(map[*node]bool)
		xEnclosedRegions := 0

		if myString.state == "" {
			for l := range myString.liberties {
				if accountedFor[gg.nodes[l]] == false {
					newXEnclosure := myMap[gg.nodes[l]]
					if newXEnclosure.boundaryColor == 2 {
						xEnclosedRegions += 1
					}
					for n := range newXEnclosure.region {
						accountedFor[n] = true
					}

				}
			}
			if xEnclosedRegions < 2 {
				isSeki, duelists := gg.inMexicanStandOff(myString)
				if isSeki {
					for _, i := range duelists {
						i.state = "seki"
					}
				} else {
					myString.state = "dead"
				}
			}
		}
	}

	return
}

// getDependents takes a GoGrah, a map connecting nodes and xEnclosedRegions, and a state ("dead" or "seki" or "alive")
func (cs chromaticStrings) getDependents(gg GoGraph, myMap map[*node]xEnclosedRegion, state string) {

	var stateBlackStrings []*stoneString
	var stateWhiteStrings []*stoneString

	for _, myString := range cs.black {
		stringPointer := &myString
		if stringPointer.state == state {
			stateBlackStrings = append(stateBlackStrings, stringPointer)
		}
	}

	for _, myString := range cs.white {
		stringPointer := &myString
		if myString.state == state {
			stateWhiteStrings = append(stateWhiteStrings, stringPointer)
		}
	}

	// We first get all dependent dead / seki for black

	stillSearching := true
	counter := 0
	accountedFor := make(map[*node]bool)

	// We record all the stones that already belong to
	// stoneStrings that are dead / seki

	for _, strings := range stateBlackStrings {
		for i := range strings.stones {
			accountedFor[gg.nodes[i]] = true
		}
	}

	for stillSearching && counter < 1000 {
		var newState []stoneString
		for _, myString := range stateBlackStrings {

			// for each liberty of the dead string
			for L := range myString.liberties {

				// get the xenclosed region of this liberty
				xEnclosed := myMap[gg.nodes[L]]
				color := xEnclosed.boundaryColor
				// region := xEnclosed.region
				// if it is in fact xenclosed
				if color == 1 {
					// iterate over neighbors until you find a black node

					for _, n := range gg.nodes[L].neighbors {
						if !accountedFor[n] && n.color == 1 && gg.stringOf[n.id].state == "" {
							if len(gg.stringOf[n.id].liberties) == 2 {
								newState = append(newState, gg.stringOf[n.id])

								for q := range gg.stringOf[n.id].stones {
									accountedFor[gg.nodes[q]] = true

								}

							}
						}
					}

				}
			}
		}

		if len(newState) < 1 {
			stillSearching = false
		}
		// add the newly found stoneStrings to the pile
		for _, i := range newState {
			newPointer := &i
			stateBlackStrings = append(stateBlackStrings, newPointer)
		}

		counter += 1

	}

	// We now repeat the process for white

	stillSearching = true
	counter = 0

	// We record all the stones that already belong to
	// stoneStrings having a particular state
	for _, strings := range stateWhiteStrings {
		for i := range strings.stones {
			accountedFor[gg.nodes[i]] = true
		}
	}

	for stillSearching && counter < 1000 {
		var newState []stoneString
		for _, myString := range stateWhiteStrings {

			// for each liberty of the dead string
			for L := range myString.liberties {

				// get the xenclosed region of this liberty
				xEnclosed := myMap[gg.nodes[L]]
				color := xEnclosed.boundaryColor
				// region := xEnclosed.region

				// if it is in fact xenclosed
				if color == 2 {
					// iterate over neighbors until you find a white node
					// whose associated string
					// 1. has not been visited yet
					// 2. has not been assigned a state

					for _, n := range gg.nodes[L].neighbors {
						if accountedFor[n] == false && n.color == 2 && gg.stringOf[n.id].state == "" {
							if len(gg.stringOf[n.id].liberties) == 2 {
								newState = append(newState, gg.stringOf[n.id])

								for q := range gg.stringOf[n.id].stones {
									accountedFor[gg.nodes[q]] = true

								}

							}
						}
					}

				}
			}
		}
		if len(newState) < 1 {
			stillSearching = false
		}
		// add the newly found stoneStrings to the pile
		for _, i := range newState {
			newString := &i
			stateWhiteStrings = append(stateWhiteStrings, newString)
		}

		counter += 1

	}

	// hopefully, because we are using slices of
	// pointers, it will update the copies that are
	// found in chromatic strings.
	for _, z := range stateBlackStrings {
		z.state = state
	}
	for _, z := range stateWhiteStrings {
		z.state = state
	}

	return
}

func cleanerBoard(gg GoGraph, cs chromaticStrings) (GoGraph, map[*node]bool, map[*node]bool) {

	cleanBoard := gg

	deadBlackStones := make(map[*node]bool)

	deadWhiteStones := make(map[*node]bool)

	for _, strings := range cs.black {
		if strings.state == "dead" {
			for z := range strings.stones {
				deadBlackStones[gg.nodes[z]] = true
			}
		}
	}
	for _, strings := range cs.white {
		if strings.state == "dead" {
			for z := range strings.stones {
				deadWhiteStones[gg.nodes[z]] = true
			}
		}
	}

	for z := range deadBlackStones {
		cleanBoard.nodes[z.id].color = 0
	}

	for z := range deadWhiteStones {
		cleanBoard.nodes[z.id].color = 0
	}
	return cleanBoard, deadBlackStones, deadWhiteStones

}

// getNewGoGraph takes
func getNaiveScoreSuggestion(gg GoGraph, cs chromaticStrings) (float64, float64, map[*node]bool) {

	var blackScore float64
	var whiteScore float64

	dead := make(map[*node]bool)

	myMap := getAllXEnclosedRegions(gg)
	cs.getSimpleDeadandSeki(gg, myMap)
	cs.getDependents(gg, myMap, "dead")
	cs.getDependents(gg, myMap, "seki")

	cleanBoard, deadBlackStones, deadWhiteStones := cleanerBoard(gg, cs)

	newEmpties := getEmptyNodes(cleanBoard)
	newMap := getAllXEnclosedRegions(cleanBoard)

	for z := range newEmpties {
		if newMap[z].boundaryColor == 1 {
			blackScore += 1
		} else if newMap[z].boundaryColor == 2 {
			whiteScore += 1
		}
	}

	for z := range deadWhiteStones {
		dead[z] = true
	}
	for z := range deadBlackStones {
		dead[z] = true
	}

	blackScore += float64(len(deadBlackStones))
	whiteScore += float64(len(deadWhiteStones))

	return blackScore, whiteScore, dead
}
