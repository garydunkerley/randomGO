package game

//TODO: this function is currently written to take an input node y,
//check the global game state to get groups, and augment the same state with new groups.
//It returns nothing.
//Instead, we need to remove the global state modification (so no *stoneGroup things)
//And we need to return all the stoneGroups (as a slice or a type)
//Note also: the "get neighbors" functions will be modified to return slices of node IDs.
//Formerly, they returned slices of node pointers.

//TODO: On the same note: please ensure that these functions do not change the global state.
func modifyStonegroups(y *node) {
	var friendlies []*node
	var enemies []*node

	var oldGroup *stoneGroup
	var mergedGroup *stoneGroup

	newStones := make(map[*node]bool)
	newLibs := make(map[*node]bool)
	newNodeLibs := getLiberties(y)

	friendlies = getSameColorNeighbors(y)
	enemies = getOppColorNeighbors(y)

	if len(friendlies) > 0 {
		// joins together neighboring friendly groups
		mergedGroup = friendlies[0].group

		for i := 0; i < len(friendlies); i++ {
			oldGroup = friendlies[i].group
			// set oldGroup to the group of the ith friendly node
			for key := range oldGroup.liberties {
				if key != y {
					// foreach liberty of the old stone group, set it as
					// a liberty of the new one
					newLibs[key] = true
				}
			}

			for key := range oldGroup.stones {
				newStones[key] = true
				key.inGroup = mergedGroup
				// should change each of these nodes so that they now
				// point to the group for y.
			}

		}
		for _, key := range newNodeLibs {
			newLibs[key] = true
		}
		newStones[y] = true

		mergedGroup.stones = newStones // update the stones and liberties of y's group
		mergedGroup.liberties = newLibs
		y.group = mergedGroup
	} else {
		// if y is not next to any friendly stones
		// make a group containing y
		var newGroup stoneGroup
		newStones[y] = true
		for _, z := range getLiberties(y) {
			newLibs[z] = true
		}
		newGroup.stones = newStones

		newGroup.liberties = newLibs

		y.group = &newGroup
	}
	if len(enemies) > 0 {
		for _, z := range enemies {
			delete(z.group.liberties, y)
		}
	}
}

func (x boardState) removeDead(y *node) {
	var enemies []*node
	enemies = getOppColorNeighbors(y)
	if len(enemies) == 0 {
		return nil
	} else {
		for i := 0; i < len(enemies); i++ {
			if len(enemies[i].group.liberties) == 0 {
				// if enemy group has no liberties
				for z := range enemies[i].group.stones {
					z.color = 0
					// change all stones into empty nodes
					delete(enemies[i].group.stones, z) // remove each stone from the map
				}
			}
		}
	}
	return nil
}
