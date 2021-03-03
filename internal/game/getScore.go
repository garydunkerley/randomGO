package game

import (
	"fmt"
)

// TODO:  

/*
NOTE
A scoring algorithm needs to be able to locate x-enclosed regions and then determine whether the blocks that comprise the boundary of these regions are healthy or not.

We don't need to detect unconditional life necessarily. I think it will be sufficient to show that all constituent blocks of a group enclose regions that are healthy for them.

If a group W of white stones sits inside of a region B of blakc stones and does not satisfy this "two healthy regions" condition, then we regard the stones as dead for scoring purposes.


*/

func getEmptyNodes(gg GoGraph) map[*node]bool{
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
	accounterFor := make(map[*node]bool)
	stillSearching := true
	counter := 0
	
	if n.color == 0 {
		accountedFor[n] = true
		for stillSearching && counter < 1000{
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
func getBoundaryColor(emptyRegion map[*node]bool) int8{
	var color int8
	var oldColor int8
	for z := range emptyRegion {
		for _, n := range z.neighbors && n.color != 0 {
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
func (gg GoGraph) inMexicanStandOff(myString stoneString) (bool, []stoneString) {
	// if my stoneString has two liberties
	if len(stoneString.liberties) == 2 {
		// then for each stone in the string, look a neighbor of the opposite color and check to see if it also has two liberties and at least one is shared
		duelDiscovered := false
		var duelists []stoneString

		for i := range myString.stones {
			for _, n := range gg.nodes[i].neighbors {
				if n.color != myString.color && n.color != 0 {
					potentialDuelist := gg.inGroup(n)
					if len(potentialDuelist.liberties) == 2{
						for z := range potentialDuelist.liberties {
							if myString.liberties[z] {
								duelists = append(dualists, myString)
								duelists = append(dualists, potentialDuelist)
								duelDiscovered = true
								break
							}
						}
					}
				}
			
				if duelDiscovered {
					break

				}

			}
			if duelDiscovered {
				break
			}

		}

		if duelDiscovered {
			// if we find that there is a mexican standoff between two stoneStrings, 
			// then we check to see if there are nearby stone strings that are also 
			// implicated in the standoff
			stillSearching := true 
			counter := 0
			accountedFor := make(map[*node]bool)
			var newPotentialDuelist stoneString

			oldDuelists := dualists
			
			// TODO: remove counter when it has been determined it won't run forever
			for stillSearching && counter < 1000 {

				newDuelistRep := make(map[*node]bool)
				newDuelists := make([]stoneString)
				
				for _, i := range oldDuelists {
					for _, l := range i.liberties {
						for _, n := range gg.nodes[l].neighbors{ 
							if n.color == i.color && !(MapKeysEqual(i.stones, gg.stringOf(n).stones) && accountedFor(n) == false {
								newPotentialDuelist = gg.stringOf[n.id]
								if len(newPotentialDuelist.liberties) == 2 {
									newDuelistRep[n] = true
									for z := range newPotentialDuelist.stones {
										accountedFor[gg.nodes[i]] = true
									}

								}
						
							}
						}
					}
				}
				if len(newDuelistRep) < 1 {
					stillSearching = false
				} else {
					for z := range newDuelistRep {
						newDuelists = append(newDuelists, gg.stringOf[z])
					}
					for z := range newDuelists {
						duelists = append(duelists, newDuelists)
					}
					oldDuelists = newDuelists
				}
			// TODO: remove counter
			counter += 1
			}
		}
		return true, duelists
	}
	return false, duelists
}

func (cs chromaticStrings) getDeadandSeki(gg GoGraph, myMap map[*node]xEnclosedRegion) {

	for _, myString := range cs.black {
		accountedfor := make(map[*node]bool)
		xEnclosedLibs := 0
		
		if myString.state == "" {
			for l := range myString.liberties {
				if accountedFor[l] == false {
					newXEnclosure := xEnclosedRegionOfNode(gg.nodes[l])
					if newXEnclosure.boundaryColor == 1 {
						xEnclosedLibs += len(newXEnclosure.region)
					}
					for n := range newXEnclosure.region {
						accountedFor[n] = true
					}


			}
			if xEnclosedLibs < 2 {
				isSeki, duelists := inMexicanStandOff(myString)
				if isSeki {
					for _,  i := range duelists {
						i.state = "seki"
					}
				} else {
					myString.state = "dead"
				}
			}
			// this accounts for stone blocks whose life or death depends on the 
			// life or death of other stone blocks whose liberties are shared.
			// This happens when two blocks are connected by a join and one of the
			// blocks has no properly contained liberties
			if xEnclosedLibs == 2 {
				for l := range myString.liberties {
					for _, n := l.neighbors {
						if n.color == myString.color && myString.stones[n.id] == false {
							if gg.stringOf(n).state == "dead" {
								myString.state = "dead"
							}
						}
					}
				}
			}
		}
	}

	for _, myString := range cs.white{
		accountedfor := make(map[*node]bool)
		xEnclosedLibs := 0
		
		if myString.state == "" {
			for l := range myString.liberties {
				if accountedFor[l] == false {
					newXEnclosure := xEnclosedRegionOfNode(gg.nodes[l])
					if newXEnclosure.boundaryColor == 2 {
						xEnclosedLibs += len(newXEnclosure.region)
					}
					for n := range newXEnclosure.region {
						accountedFor[n] = true
					}


			}
			if xEnclosedLibs < 2 {
				isSeki, duelists := inMexicanStandOff(myString)
				if isSeki {
					for _,  i := range duelists {
						i.state = "seki"
					}
				} else {
					myString.state = "dead"
				}
			}
			// this accounts for stone blocks whose life or death depends on the 
			// life or death of other stone blocks whose liberties are shared.
			// This happens when two blocks are connected by a join and one of the
			// blocks has no properly contained liberties
			if xEnclosedLibs == 2 {
				for l := range myString.liberties {
					for _, n := l.neighbors {
						if n.color == myString.color && myString.stones[n.id] == false {
							if gg.stringOf(n).state == "dead" {
								myString.state = "dead"
							}
						}
					}
				}
			}
		}
	}

	return
}



// TODO build on this so that it can actually account for shit like dead groups, seki, etc.
func getNaiveScoreSuggestion(gg GoGraph) (float64, float64) {
	
	var graphWithoutDead GoGraph

	var whiteNaive float64
	var blackNaive float64
	whiteNaive = 0
	blackNaive = 0

	myScoringMap := getAllXEnclosedRegions(gg)

	accountedFor := make(map[*node]bool)



	emptyNodes := getEmptyNodes(gg)

	//
	for z := range emptyNodes {
		if accountedFor[z] == true {
			continue
		} else {
			emptyClique := getEmptyClique(z)
			color, points := getScoreAssignment(emptyClique)
			if color == "black" {
				blackNaive += float64(points)
			} else if color == "white" {
				whiteNaive += float64(points)
			} else {
				continue
			}
			for n := range emptyClique {
				accountedFor[n] = true
			}
		}

	}
	return whiteNaive, blackNaive

}
