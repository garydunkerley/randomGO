package game

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"errors"
)

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func sgfToString() string {
	file, err := ioutil.ReadFile("path_to_file")
	checkError(err)

	return file
}

// Really, this thing should search through the file for the string "SZ[" and then output whatever integer appears after it. We can't assume that line number reading will always work.
// Generalize later so that it just generates the appropriate board topology for a game.
func getSquareBoardFromSGF(sgfString []string) boardTop {
	boardSize, err := strconv.Atoi(strings.Trim(sgfString[12], "SZ[]"))
	if err == nil {
		myBoard := makeSquareBoard(boardSize, boardSize)
		return myBoard
	} else {
		// indicate that an error occurred I guess?
	}
}


// Instantiate a root.
// Locate the first instance of a ( symbol, call it (0
// note the location of the final move, m0,  prior to (0.
// Move from the (0 to the next instance of a ( symbol, (1,  and repeat until you encounter a ) symbol. Take all moves encountered in this line to construct the main branch.

// Now starting from the m0 encountered prior to (0,
// parse through the file adding +1 to a counter for each
// instance of ( and -1 for each instance of ).
// Once the counter is zero, check for a ( symbol as this
// will indicate that there is a new branch whose root will be
// the move m.

// How about we split the imported sgf data into two pieces: one with game metadata and one with just the moves.
// should be possible to do this by splitting at the first instance of the substring ";B["

// TODO: 	1. Split the SGF into two strings: metadata and moves
//		2. Convert each of these into rune arrays so that we can iterate over each rune to parse branching

//		3. Create either a function or a struct that outputs / stores the root and the move slice

//		4. Modify the createBranch function so that it works with the object / output of the function described in 3.

// Items to be converted into treeBranches need to be string, []string pairs where the first string will denote the root and the string slice will denote the move set located in this branch of the game tree.


// separateMetadata takes the raw string for an sgf file and 
// splits it into two strings. One containing only the metadata 
// about the game (board size, player names, etc) and the other
// containing branches in the game tree and comments
func separateMetadata(sgfData string) []string {
	myStrings := strings.SplitN(sgfData, ";B[", 2)
	myStrings[1] = ";B[" + myStrings
	
	return myStrings
}

// strSingletonArray takes a string and splits it into a slice of strings each of length 1. 
func strSingletonArray(s string) []string {
	myStrings := strings.Split(s, "")
	
}

// getBranchGenesis takes a list of strings and a particular location and determines 
// when the time at which the list begins to encode a new branch
func getBranchGenesis(symbols []string, searchStart int) (int, error) {
	parenCounter := 0

	// beginning at location searchStart, we count the number of (signed) parentheses
	// to find where the information encoding the new branch is located.
	for i := searchStart; i < len(symbols); i++ {
		if symbols[i] == "(" {
			parenCounter += 1
		} else if symbols[i] == ")" {
			parenCounter -= 1
		}
	} if parenCounter == 0 {
		// if new branch is non-empty
		// give the first symbol in the new branch

		if symbols[i+1] == "(" {
			return i+2
		// if the branch is empty, return a value 
		// that indicates such
		} else if symbols[i+1] == ")" {

			fmt.Println("Debug: New branch is empty")
			return -1 , nil

		// otherwise, there's something wrong
		} else {
			return -1, errors.New("Error: Unexpected symbol at beginning of new branch")
		}
	}
}


// have this output the distinct string representing the 
// root 
func getRootNameForBranch(symbols []string, location int) (string, error) {
	if symbols[location] != ( {
		return "Error", errors.New("Debug: This location is not a branching point!")
	} else {
		for j:= 0; j < location + 1; j++ {
			// work backwards until you encounter a move
			if symbols[location-j] == ";" {
				// save the location of the move
				rootName := isolateMoveString(location - j + 1)
				rootName = rootName + " @(*THIS_IS_A_ROOT*)@"
			}
		}
	}
}

// isolateMoveString grabs all information pertaining to a move beginning at a particular location in mSymb 
func isolateMoveString(mSymb []string, location int) string {
	moveString := ""
	for i:= location; i < len(mSymb)+1; i++ {
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
	return (mSymb[i] == ";" && (mSymb[i+1]+mSymb[i+2] == "B[" || mSymb[i+1]+mSymb[i+2] == "W[" ))   
}


// getMoveInput takes a string that encodes a single move, 
// decodes the associated node color and location, 
// and outputs the moveInput data
getMoveInput(s string) (moveInput, error) {
	var newMoveInput moveInput 
	
	// chops off the first part of the string
	beHeaded := strings.SplitN(s, "[", 2)
	
	// head will contain sonte color information
	head = beHeaded[0]

	// body will be further modified to get the NodeId
	body = beHeaded[1]
	 
	if strings.Contains(head, "B") {
			newMoveInput.playerColor = 1
	} else if strings.Contains(head, "W") {
		newMoveInput.playerColor = 2
	} else {
		return _, errors.New("Error (growBranch): 
		String does not encode player color")
	}


	// removes everything but the letters encoding the 
	// NodeId
	tailRemoved := strings.SplitN(body, "]", 2)[0]
	
	// convert nodeData to a rune slice so as to 
	// read the encoded NodeId

	nodeData := []rune(tailRemoved)

	if len(nodeData) != 2 && nodeData != "" {
		
		return _, errors.New("Error (procureNodeId): 
		String data cannot correspond to NodeId")

	} else if nodeData == "" {

		newMoveInput.id = -1
		newMoveInput.isPass = true

	} else {

		newMoveInput.id = runeToID(nodeData[0]) + runeToID(nodeData[1])
	}
	
	return newMoveInput , nil
}

// growBranch takes a collection of strings and a treeNode and outputs the collection of treeNodes that compose the associated branch.
func (root *treeNode) growBranch(movesInBranch []string, branchNum int) error {
	
	// a place to store the names of the treeNodes
	branchChart := make(map[string]treeNode)

	var nodeData []rune
	var newMove move

	for z := range movesInBranch {
		// iterate over the moves and do some checks
		
		newMove.moveInput, _ = getMoveInput(z)
		
		// create a new child of the root treeNode
		newNode := root.newChild(newMove)

		// TODO: find a way to create a naming convention so that we 
		// can identify 


		// now set things up so that the next treeNode will have
		// the newly created one as its parent
		root = newNode
			
		
	}

}


// buildTree takes a slice of strings (each string having length 1), records the moves comprising a given branch as well as the roots and beginnings of subsequent branches, constructing the entire tree inductively.
func buildTree(mSymb []string, root *treeNode) error {
	// n := getBranchCount(moveSymbols)

	for z := range mSymb {
		if len(z) != 1 {
			return errors.New("Error: Non-singleton detected.")
		}
	}
 
	
	var branchGeneses []int{0}
	var newBranchGeneses []int

	var rootLocations []int
	var newRootLocations []int

	var roots []*treeNode{root}
	var newRoots []*treeNode

	// stores the strings that encode moves in a given branch, the main branch will begin
	// with an empty "root" node
	var movesInBranch := []string{"root"}
	
	var rootMove move
	// root should take the empty move and an empty pointer for its parent
	root := initTreeNode(rootMove, )


	gameTree := make([string] treeNode)
	gameTree("root") = &root 

	// iterate while there are new branches to make
	for branchGeneses != [] {
		newBranchGeneses = [] // set up a new save location for roots
		for z := range branchGeneses {
			
			// begin constructing a branch at each root
			for i := z; i < len(mSymb)+1; i++ {

				// if the location corresponds to the beginning of a move,
				// construct and save the string encoding that move.
				if isMoveBeginning(mSymb, i) {
			
					// find the string encoding the located move and 
					// append it to the list of moves in the branch

					movesInBranch = append(movesInBranch, isolateMoveString(mSymb, i+1))

				} else if mSymb[i] == "(" {
				
					// store the location of the root
					newRoot := getRootForBranch(mSymb, i) 
					newRootLocations = append(newRootLocations, newRoot) 

					newBranchGenesis := getBranchGenesis(mSymb, i)
					newBranchGeneses = append(newBranchGeneses, newBranchGenesis)

					movesinBranch[len(movesinBranch)-1] = movesinBranch[len(movesinBranch-1)] + " @(*THIS_IS_A_ROOT*)@"
				// encountering a ")" means that the branch we are on has ended 
				} else if mSymb[i] == ")" {

					constructBranch(movesInBranch)
					// begin construction of branch 


					// don't forget to empty the old branch!
					// movesInBranch = []

					// reset the root genesis locations to iterate over
					// branches emerging from your newly constructed 
					// branch

					// rootGenesisLoc = newRootGenesisLoc
				}
			}
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
func (parent *treeNode) newChild(m move) *treeNode {

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





// latinAlphaToNumber assigns a-z to 0-25 and A-Z to 26-51
func runeToID(letter rune) int {
	// if the letter is upper case
	if int(letter) > 64 && int(letter) < 91 {
		return (int(letter) % 65) + 26
		// the letter is lowercase
	} else if int(letter) > 96 && int(letter) < 123 {
		return int(letter) % 97
	} else {
		return errors.New("Error: not a valid rune")
	}
}
