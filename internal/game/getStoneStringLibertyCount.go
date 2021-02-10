package game

func getLibertiesOfStoneString(S stoneString) int {
	liberties := make(map[int]bool)
	for x := range S.stones {
		for y := range getLiberties(x) {
			liberties[y] = true
		}

	}
	return len(liberties)
}
