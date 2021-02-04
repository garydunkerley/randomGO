package gameplay

func modify_stoneGroups(y *node) {
	var friendlies []*node
	var enemies []*node

	var old_Group *stoneGroup
	var merged_Group *stoneGroup

	newStones := make(map[*node]bool)
	newLibs := make(map[*node]bool)
	newNodeLibs := get_liberties(y)

	friendlies = get_same_color_neighbors(y)
	enemies = get_opp_color_neighbors(y)

	if len(friendlies) > 0 {
		// joins together neighboring friendly groups
		merged_Group = friendlies[0].inGroup

		for i := 0; i < len(friendlies); i++ {
			old_Group = friendlies[i].inGroup
			// set old_Group to the group of the ith friendly node
			for key := range old_Group.liberties {
				if key != y {
					// foreach liberty of the old stone group, set it as
					// a liberty of the new one
					newLibs[key] = true
				}
			}

			for key := range old_Group.stones {
				newStones[key] = true
				key.inGroup = merged_Group
				// should change each of these nodes so that they now
				// point to the group for y.
			}

		}
		for _, key := range newNodeLibs {
			newLibs[key] = true
		}
		newStones[y] = true

		merged_Group.stones = newStones // update the stones and liberties of y's group
		merged_Group.liberties = newLibs
		y.inGroup = merged_Group
	} else {
		// if y is not next to any friendly stones
		// make a group containing y
		var new_group stoneGroup
		newStones[y] = true
		for _, z := range get_liberties(y) {
			newLibs[z] = true
		}
		new_group.stones = newStones

		new_group.liberties = newLibs

		y.inGroup = &new_group
	}
	if len(enemies) > 0 {
		for _, z := range enemies {
			delete(z.inGroup.liberties, y)
		}
	}
}

func (x boardState) remove_dead(y *node) error {
	var enemies []*node
	enemies = get_opp_color_neighbors(y)
	if len(enemies) == 0 {
		return nil
	} else {
		for i := 0; i < len(enemies); i++ {
			if len(enemies[i].inGroup.liberties) == 0 {
				// if enemy group has no liberties
				for z := range enemies[i].inGroup.stones {
					z.color = 0
					// change all stones into empty nodes
					delete(enemies[i].inGroup.stones, z) // remove each stone from the map
				}
			}
		}
	}
	return nil
}
