package gameplay

// this script combines the get_same_color and get_opp_color functions from the original script.

func getNeighbors(x *node, neighb_color int) {
	var friends []*node
	var enemies []*node
	var liberties []*node

	if neighb_color == x.color {
		for _, y := range x.neighbors {
			if y.color == x.color {
				friends = append(friends, y)
			}
		}

		return friends

	} else {

		if neighb_color != 0 {

			for _, y := range x.neighbors {
				if y.color != x.color && y.color != 0 {
					enemies = append(enemies, y)
				}
			}

			return enemies

		} else {
			for _, y := range x.neighbors {
				if y.color == 0 {
					liberties = append(liberties, y)
				}
			}
		}

	}
}
