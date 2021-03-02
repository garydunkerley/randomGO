package game

import (
	"fmt"
)

// TODO: Create program that can compute unconditional life of a stone group.

// isemptyNode takes a GoGraph and spits out a map representing
// the set of empty nodes in the GoGraph that may contribute to territory
func getEmptyNodes(gg GoGraph) map[*node]bool {
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
	stillSearching := true

	// accountedFor := make(map[*node]bool)
	oldEmpties := make(map[*node]bool)

	counter := 0

	if n.color == 0 {
		oldEmpties[n] = true
		// accountedFor[n] = true
		// newEmpties := make(map[*node]bool)
		for stillSearching && counter < 100 {
			newEmpties := make(map[*node]bool)

			for z := range oldEmpties {
				for _, i := range z.neighbors {
					if i.color == 0 {
						if oldEmpties[i] == false {
							// accountedFor[i] = true
							newEmpties[i] = true
						}
					}
				}
			}

			if len(newEmpties) < 1 {
				stillSearching = false
			}

			// clear the way for new stones

			for z := range newEmpties {
				oldEmpties[z] = true
			}
			// fmt.Println("oldEmpties are ", oldEmpties)

		}
	} else {
		fmt.Println("Error (getEmptyClique): Node is not empty.")
	}
	return oldEmpties

}

// getScoreAssignment takes an emptyClique (a set of empty nodes) and
func getScoreAssignment(emptyClique map[*node]bool) (string, int) {
	allSameColor := true
	var color string
	var oldColor string

	for z := range emptyClique {
		if allSameColor == false {
			break
		} else {
			for _, n := range z.neighbors {
				if emptyClique[n] == false {
					oldColor = color

					if n.color == 1 {
						color = "black"
					} else {
						color = "white"
					}

					if color != oldColor && oldColor != "" {
						allSameColor = false
						break
					}
				}
			}
			return color, len(emptyClique)
		}
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
				fmt.Println(n)
				accountedFor[n] = true
			}
		}

	}
	return whiteNaive, blackNaive

}
