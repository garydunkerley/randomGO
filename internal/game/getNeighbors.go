package game

// this script combines the get_same_color and get_opp_color functions from the original script.
// TODO: change these functions to return slices of node IDs, instead of slices of node pointers.

// getSameColorNeighbors returns the node IDs of all adjacent same-color nodes.
func getSameColorNeighbors(x *node) (friends []int) {
	if x.color == 0 {
		return friends
	}
	for _, y := range x.neighbors {
		if y.color == x.color {
			friends = append(friends, y.id)
		}
	}
	return friends
}

// getOppColorNeighbors returns the node IDs of all adjacent opp-color nodes.
func getOppColorNeighbors(x *node) (enemies []int) {
	if x.color == 0 {
		return enemies
	}
	for _, y := range x.neighbors {
		if y.color != x.color && y.color != 0 {
			enemies = append(enemies, y.id)
		}
	}
	return enemies
}

// getLiberties returns the node IDs of all adjacent empty nodes.
func getLiberties(x *node) (liberties []int) {
	for _, y := range x.neighbors {
		if y.color == 0 {
			liberties = append(liberties, y.id)
		}
	}

	return liberties
}
