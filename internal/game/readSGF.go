package game

import (
	"errors"
	"io/ioutil"
	"strings"
)

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func sgfToString() string {
	file, err := ioutil.ReadFile("path_to_file")
	checkError(err)

	return string(file)
}

// Really, this thing should search through the file for the string "SZ[" and then output whatever integer appears after it. We can't assume that line number reading will always work.
// Generalize later so that it just generates the appropriate board topology for a game.
/*
func getSquareBoardFromSGF(sgfString []string) boardTop, error {
	boardSize, err := strconv.Atoi(strings.Trim(sgfString[12], "SZ[]"))
	if err == nil {
		myBoard := makeSquareBoard(boardSize, boardSize)
		return myBoard, nil
	} else {
		return _, errors.New("")
		// indicate that an error occurred I guess?
	}
}
*/

// separateMetadata takes the raw string for an sgf file and
// splits it into two strings. One containing only the metadata
// about the game (board size, player names, etc) and the other
// containing branches in the game tree and comments
func separateMetadata(sgfData string) []string {

	myStrings := strings.SplitN(sgfData, ";B[", 2)

	myStrings[1] = ";B[" + myStrings[1]
	return myStrings
}

// strSingletonArray takes a string and splits it into a slice of strings each of length 1.
func strSingletonArray(s string) []string {
	myStrings := strings.Split(s, "")
	return myStrings

}

// getBranchGenesis takes a list of strings and a particular location and determines
// when the time at which the list begins to encode a new branch
// and whether that new branch will be empty
func getBranchGeneses(symbols []string, searchStart int) []int {
	parenCounter := 0

	// beginning at location searchStart, we count the number of (signed) parentheses
	// to find where the information encoding the new branch is located.
	var branchGeneses []int
	for i := searchStart; i < len(symbols); i++ {
		if symbols[i] == "(" {
			parenCounter += 1
		} else if symbols[i] == ")" {
			parenCounter -= 1
		}
		if parenCounter == 0 {
			if symbols[i+1] == "(" {
				branchGeneses = append(branchGeneses, i+2)
			} else if symbols[i+1] == ")" {
				break
			}
		}
	}
	return branchGeneses
}

// isolateMoveString grabs all information pertaining to a move beginning at a particular location in mSymb
func isolateMoveString(mSymb []string, location int) string {
	moveString := ""
	for i := location; i < len(mSymb)+1; i++ {
		if mSymb[i] != ";" && mSymb[i] != "(" && mSymb[i] != ")" {
			moveString = moveString + mSymb[i]
		} else {
			break
		}
	}
	return moveString
}

// isMoveBeginning just checks to see if the location in mSymb corresponds to the start of an encoding for a move.
func isMoveBeginning(mSymb []string, location int) bool {
	return (mSymb[location] == ";" && (mSymb[location+1]+mSymb[location+2] == "B[" || mSymb[location+1]+mSymb[location+2] == "W["))
}

// isBranchGenesis checks to see if branching is occurring at
// specified location
func isBranchGenesis(mSymb []string, location int) bool {
	return mSymb[location] == "("
}

func isBranchEnd(mSymb []string, location int) bool {
	return mSymb[location] == ")"
}

// getMove takes a string that encodes a move by
// decoding the associated node color and location
func getMove(s string) (move, error) {
	var newMoveInput moveInput
	var newMove move

	// chops off the first part of the string
	beHeaded := strings.SplitN(s, "[", 2)

	// head will contain color information
	head := beHeaded[0]

	// body will be further modified to get the NodeId
	body := beHeaded[1]

	if strings.Contains(head, "B") {
		newMoveInput.playerColor = 1
	} else if strings.Contains(head, "W") {
		newMoveInput.playerColor = 2
	} else {
		return newMove, errors.New("Error (growBranch): String does not encode player color")
	}

	// removes everything but the letters encoding the
	// NodeId
	tailRemoved := strings.SplitN(body, "]", 2)[0]

	// convert nodeData to a rune slice so as to
	// read the encoded NodeId

	nodeData := []rune(tailRemoved)

	if len(nodeData) != 2 {

		return newMove, errors.New("Error (procureNodeId): String data cannot correspond to NodeId")

	} else if len(nodeData) == 0 {

		newMoveInput.id = -1
		newMoveInput.isPass = true

	} else {

		newMoveInput.id = runeToID(nodeData[0]) + runeToID(nodeData[1])
	}

	newMove.moveInput = newMoveInput

	return newMove, nil
}

// growBranch is a method on a treeNode that takes a slice of strings,
// a location in the slice, and a counter. It iteratively adds children
// according to the instructions starting at the specified location in
// mSymb. As it encounters the roots of new branches, it records them
// as well as the locations in mSymb corresponding to their instructions.
// It also outputs a map that will be used to construct a global map for the tree.
func (root *treeNode) growBranch(mSymb []string, location int, nodeCounter int) ([]*treeNode, []int, int, map[int]treeNode) {

	// a place to store the names of the treeNodes
	branchChart := make(map[int]treeNode)

	var newRoots []*treeNode

	var newBranchGeneses []int
	for z := range mSymb {
		if isMoveBeginning(mSymb, z) {
			// make a move corresponding to the move string encoding found at location z
			newMove, _ := getMove(isolateMoveString(mSymb, z))

			// create a child of the root
			child := root.newChild(newMove)
			nodeCounter += 1
			branchChart[nodeCounter] = *child

			// designate the child as the new root
			root = child

		} else if isBranchGenesis(mSymb, z) {
			for s := range getBranchGeneses(mSymb, z) {
				newRoots = append(newRoots, root)
				newBranchGeneses = append(newBranchGeneses, s)
			}
		} else if isBranchEnd(mSymb, z) {
			break
		}

	}
	return newRoots, newBranchGeneses, nodeCounter, branchChart

}

// buildTree takes a slice of strings (each string having length 1) and recursively constructs
// the gameTree by constructing a branch, for each root of a new branch, it records the root
// as well as the location in mSymb where the new branch is encoded. At each node initialization,
// a map associates an integer with the new node.
// The code iterates until there are no more branches to grow.
func buildTree(mSymb []string) (map[int]treeNode, error) {

	var rootMove move
	root := initTreeNode(rootMove)
	roots := []*treeNode{root}

	branchGeneses := []int{0}
	branchesToGrow := true

	nodeCounter := 0
	gameTree := make(map[int]treeNode)

	// while there are still branches to grow, construct branches
	for branchesToGrow {
		if len(roots) != len(branchGeneses) {
			return gameTree, errors.New("Error (buildTree): len(roots) =/= len(branchGeneses) ")
		}
		newRoots, newBranchGeneses, newNodeCounter, chart := roots[0].growBranch(mSymb, branchGeneses[0], nodeCounter)

		for z := range chart {
			gameTree[z] = chart[z]
		}

		roots = append(roots[1:], newRoots...)
		branchGeneses = append(branchGeneses[1:], newBranchGeneses...)
		nodeCounter = newNodeCounter
		// if there are no new branches to grow, end the loop
		if len(roots) == 0 {
			branchesToGrow = false
		}
	}

	return gameTree, nil

}

// initTreeNode specifies a node in the game tree by a move, treeNode pointer pair.
func initTreeNode(value move) *treeNode {
	var myTreeNode treeNode

	myTreeNode.value = value

	return &myTreeNode

}

// newChild method on a parent accepts a move and modifies the parent so that it has a child corresponding to the specified move
func (parent *treeNode) newChild(m move) *treeNode {

	child := initTreeNode(m)
	child.parent = parent

	parent.children = append(parent.children, child)

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

// latinAlphaToNumber assigns a-z to 0-25 and A-Z to 26-51
func runeToID(letter rune) int {
	// if the letter is upper case

	var integer int

	if int(letter) > 64 && int(letter) < 91 {
		integer = (int(letter) % 65) + 26
		// the letter is lowercase
	} else if int(letter) > 96 && int(letter) < 123 {
		integer = int(letter) % 97
	}
	return integer
}
