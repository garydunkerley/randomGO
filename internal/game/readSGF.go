package game

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func sgfToStringSlice() []string {
	file, err := ioutil.ReadFile("path_to_file")
	check(err)
	myStrings := strings.Split(string(file), "\n")

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

func (h history) populateHistory(sgfString []string) {
	// 1. Should convert alphabetical data to node numbers.
	// 2. Should append these node numbers to the move struct
	// 3. should construct koHistory
	// 4. Should construct and store list of
	return
}

func getMoveHistory(sgfString []string) []move {
	var nodeData []rune
	var newMove move
	var newMoveInput moveInput

	var myMoves []move

	for z := range sgfString {
		// if a string in the sgfString slice
		// encodes a move
		if strings.Contains(";") {
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
			myMoves = append(myMoves, newMove)
		}

	}

	return myMoves

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
