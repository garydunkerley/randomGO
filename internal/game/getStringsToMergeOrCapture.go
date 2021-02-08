package game

// getOppColor(2)=1, getOppColor(1)=2
func getOppColor(color int8) int {
	return 3 - color
}

// getSubsumedStrings gives us a map encoding the stoneStrings that need
// to be subsumed into a larger stoneString given a play at nodeId of a given color

func (G GoGraph) getSubsumedStrings(nodeId int, color int8) map[stoneString]bool {
	stringsToBeMerged := make(map[stoneString]bool)
	var friendlies []int
	for z := range G.nodes[nodeId].neighbors {
		if z.color == color {
			friendlies = append(friendlies, z.id)
		}
	}
	for z := range friendlies {
		stringsToBeMerged[G.stringOf(z)] = true
	}
	return stringsToBeMerged
}

// getCapturedStrings gives a map encoding the stoneStrings
// that will be captured by a play at NodeId of a given color.
func (G GoGraph) getCapturedStrings(nodeId int, color int8) map[stoneString]bool {
	capturedStrings := make(map[stoneString]bool)
	var potentialCaptives []int

	for z := range G.nodes[nodeId].neighbors {
		if z.color == getOppColor(color) {
			potentialCaptives = append(potentialCaptives, z.id)
		}
	}

	for z := range potentialCaptives {
		if len(G.stringOf(z).liberties) == 1 {
			capturedStrings[G.stringOf(z)] = true
		}
	}
	return capturedStrings
}
