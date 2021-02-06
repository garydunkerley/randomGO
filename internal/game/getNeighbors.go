package game

// this script combines the get_same_color and get_opp_color functions from the original script.

func get_same_color_neighbors(x *node) []*node {
	var friends []*node
	if x.color == 0 {
		return friends
	}
	for _, y := range x.neighbors {
		if y.color == x.color {
			friends = append(friends, y)
		}
	}
	return friends
}

func get_opp_color_neighbors(x *node) []*node {
	var enemies []*node
	if x.color == 0 {
		return enemies
	}
	for _, y := range x.neighbors {
		if y.color != x.color && y.color != 0 {
			enemies = append(enemies, y)
		}
	}
	return enemies
}

func get_liberties(x *node) []*node {
	var liberties []*node
	if x.color == 0 {
		return liberties
	} else {
		for _, y := range x.neighbors {
			if y.color == 0 {
				liberties = append(liberties, y)
			}
		}
	}
	return liberties
}
