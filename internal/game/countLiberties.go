package game

func (myGraph GoGraph) countLiberties(S stoneString) int {
	liberties := make(map[int]bool)
	for x := range S.stones {
		for y := range getLiberties(myGraph.nodes[x]) {
			liberties[y] = true
		}
	}
	return len(liberties)
}
