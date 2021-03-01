package game

import (
	"fmt"
)

// TODO: Create program that can compute unconditional life of a stone group.

// isemptyNode takes a GoGraph and spits out a map representing
// the set of empty nodes in the GoGraph that may contribute to territory
func getEmptyNodes(gg GoGraph) map[*node]bool {
	fmt.Println("Debug: begun getEmptyNodes")
	emptyNodes := make(map[*node]bool)
	for _, z := range gg.nodes {
		if z.color == 0 {
			emptyNodes[z] = true
		}
	}
	return emptyNodes
}

// getEmptyClique gets the connected component that a stone belongs to
func getEmptyClique(n *node) (emptyClique map[*node]bool) {
	fmt.Println("Debug: begun getEmptyClique")
	stillSearching := true

	accountedFor := make(map[*node]bool)

	if n.color == 0 {
		if stillSearching {
			newEmpties := make(map[*node]bool)
			for _, z := range n.neighbors {
				if z.color == 0 {
					if accountedFor[z] {
						continue
					} else {
						accountedFor[z] = true
						newEmpties[z] = true
					}
				}

				if len(newEmpties) == 0 {
					stillSearching = false
					break
				}
			}
		}
	} else {
		fmt.Println("Error (getEmptyClique): Node is not empty.")
	}

	return accountedFor

}

// getScoreAssignment takes an emptyClique (a set of empty nodes) and
func getScoreAssignment(emptyClique map[*node]bool) (string, int) {
	fmt.Println("Debug: Begun get score assignment")
	allSameColor := true
	var color string
	var oldColor string

	for z := range emptyClique {
		if allSameColor {
			continue
		} else {
			break
		}
		for _, n := range z.neighbors {
			if emptyClique[n] {
				continue
			} else {

				if n.color == 1 {
					color = "black"
				} else {
					color = "white"
				}

				if color == oldColor || oldColor == "" {
					continue
				} else {
					allSameColor = false
					break
				}
			}
		}
		return color, len(emptyClique)

	}

	return "mixed", 0
}

// TODO build on this so that it can actually account for shit like dead groups, seki, etc.
func getNaiveScore(gg GoGraph) (float64, float64) {
	var whiteNaive float64
	var blackNaive float64
	whiteNaive = 0
	blackNaive = 0

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
