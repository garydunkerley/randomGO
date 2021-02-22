package game

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func sgfToStringSlice() []string {
	file, err := ioutil.ReadFile("path_to_file")
	check(err)
	myStrings := strings.Split(string(file), "(;")

	return myStrings
}

func getSquareBoardFromSGF(sgfString []string) boardTop {
	boardSize, err := strconv.Atoi(strings.Trim(myStrng[12], "SZ[]"))
	if err == nil {
		myBoard := makeSquareBoard(boardSize, boardSize)
		return myBoard
	} else {
		// indicate that an error occurred I guess?
	}
}

// initTreeNode specifies a node in the game tree by a move, treeNode pointer pair.
func initTreeNode(value move, parent *treeNode) *treeNode {
	var myTreeNode treeNode
	myTreeNode.value = value
	myTreeNode.parent = parent

	return *myTreeNode

}

// newChild method on a parent accepts a move and modifies the parent so that it has a child corresponding to the specified move
func (parent treeNode) newChild(m move) *treeNode {

	child := initTreeNode(m, parent)
	parent.children = append(parent.children, newChild)

	return child

}

func continueAlongCurrentBranch(crawler *treeNode) *treeNode {
	return crawler.children[0]
}

func choosePath(crawler *treeNode, i int) *treeNode {
	return crawler.children[i]
}

func moveBack(crawler *treeNode) *treeNode {
	return crawler.parent
}

// TODO: Find a way to spit the SGF moveset into parts, one object of type []string for each branch.

// createBranch generates treeNodes that are attached to a given root
func createBranch(branch []move, root *treeNode) {
	var child *treeNode

	var crawler *treeNode

	// we start by making the root have a child and set crawler to this location
	crawler = &root.newChild(branch[0])

	// for each item in the branch, we initialize a treeNode to be a child of crawler and
	// then move crawler to the newly created child
	for i := 1; i < len(branch); i++ {
		&crawler.newChild(branch[i])
		crawler = child
	}

}

// convertBranch takes a collection of strings and outputs their corresponding sequence of moves.
func convertBranch(branchString []string) []move {
	var nodeData []rune
	var newMove move

	var branch []move
	// a designated location on the tree

	// crawler will be used to move up and down the tree as we construct it starting from the root

	for z := range branchString {
		// if a string in the sgfString slice
		// encodes a move
		if strings.Contains(";") {
			// if we are beginning to look at variations, break the loop
			if strings.Contains(")") {
				break
			} else {
				// save the player color for the move
				if strings.Contains("B") {
					newMoveInput.playerColor = 1
				} else {
					newMoveInput.playerColor = 2
				}
				// isolates the individual runes that encode the nodeid
				nodeData = []rune(strings.Trim(z, "BW;()[]"))
				if nodeData == "" {
					newMoveInput.id = -1
					newMoveInput.isPass = true
				} else {
					newMoveInput.id = runeToID(nodeData[0]) + runeToID(nodeData[1])
				}
				newMove.moveInput = newMoveInput

				branch = append(branch, newMove)

			}
		}

	}

	return branch

}

// latinAlphaToNumber assigns a-z to 0-25 and A-Z to 26-51
func runeToID(letter rune) int {
	// if the letter is upper case
	if int(letter) > 64 && int(letter) < 91 {
		return (int(letter) % 65) + 26
		// the letter is lowercase
	} else if int(letter) > 96 && int(letter) < 123 {
		return int(letter) % 97
	} else {
		fmt.Println("Error: not a valid rune")
		return
	}
}
