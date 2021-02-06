package game

func initializeNode(i int, a string) *node {
	// initializes empty nodes whose ids are integers

	var x node
	x.name = a
	x.id = i
	x.color = 0
	return &x
}
