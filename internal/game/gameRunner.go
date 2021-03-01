package game

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

// checkCoordsInput checks if given input string is formatted as coordinates.
// Returns non-nil error if the input is not of the form "n-m",
// where n,m are nonnegative ints.
func checkCoordsInput(input string) error {
	stringPieces := strings.Split(input, "-")
	if len(stringPieces) != 2 {
		return errors.New("Input is not formatted as coordinates.")
	} else if stringPieces[0] == "" { //empty string case: input is "-m"
		return errors.New("Input is formatted as a negative int.")
	} else {
		_, err1 := strconv.Atoi(stringPieces[0])
		_, err2 := strconv.Atoi(stringPieces[1])
		if err1 != nil {
			return err1
		} else if err2 != nil {
			return err2
		} else {
			return nil
		}
	}
}

// parseCoords converts "a-b" to ints a, b.
// parseCoords assumes that the input is safe (checkCoordsInput returns nil error).
func parseCoords(coordInput string) (int, int) {
	stringPieces := strings.Split(coordInput, "-")
	a, _ := strconv.Atoi(stringPieces[0])
	b, _ := strconv.Atoi(stringPieces[1])
	return a, b
}

// checkIntInput checks if the given string is formatted as an integer.
func checkIntInput(input string) error {
	_, err := strconv.Atoi(input)
	return err
}

// parseInt converts "n" to int n. Assumes checkIntInput returns nil error.
func parseInt(intInput string) int {
	a, _ := strconv.Atoi(intInput)
	return a
}

// checkIntOrCoordsInput checks if the input can be parsed.
// In particular, returns nil iff checkIntInput or checkCoordsInput does.
func checkIntOrCoords(input string) error {
	if intErr := checkIntInput(input); intErr != nil {
		if coordsErr := checkCoordsInput(input); coordsErr != nil {
			return errors.New("Input cannot be parsed as an integer or coordinate pair.")
		}
	}
	return nil
}

// parseIntOrCoords converts a validated input (int or coords) into a nodeID
// and a validation bool.
// It expects a map of valid coordinate pairs and the number of nodes for the given board.
// If the input is coordinate-form, then the second return value is true when
// the node is a key for the given coordinate map.
// If the input is integer-form, we check if it's in the range -2 <= n < nodeCount.
// The constant -2 is due to negative numbers coding special instructions:
// -2 is exit, -1 is pass.
func parseIntOrCoords(input string, coo map[[2]int]int, nodeCount int) (nodeID int, validId bool) {
	if err := checkIntInput(input); err != nil {
		// Coordinate case: get the node id from the given coordinate map.
		if err2 := checkCoordsInput(input); err2 != nil {
			panic("Invalid input passed to parseIntOrCoords, validate first!")
		}
		a, b := parseCoords(input)
		nodeID, validId = coo[[2]int{a, b}]
		return nodeID, validId
	} // Integer case: check if the node id is a key for the given edge map.
	nodeID = parseInt(input)
	validId = (-2 <= nodeID) && (nodeID < nodeCount)
	return nodeID, validId
}

// promptRect creates a prompt for coordinate-based rectangular boards.
// It validates the input but does not check its legality.
func promptRect() promptui.Prompt {
	prompt := promptui.Prompt{
		Label:    "Move",
		Validate: checkIntOrCoords,
	}
	return prompt
}

// promptNodeID applies the given prompt, validates input, converts to a node id.
// If a certain number of consecutive invalid inputs are given, it passes for you.
// It only checks if the node id exists, not if the move is legal.
func promptNodeID(prompt promptui.Prompt, coo map[[2]int]int, nodeCount int) int {
	tryAgainMsg := "Try again, or enter -2 to exit."
	for attempts := 5; attempts > 0; attempts-- {
		inputString, err := prompt.Run()
		if err != nil { // malformed input string
			fmt.Println(err)
			fmt.Println(tryAgainMsg)
			continue
		} // good input, convert to a meaningful int
		nodeID, ok := parseIntOrCoords(inputString, coo, nodeCount)
		if !ok {
			fmt.Println("Invalid node ID.")
			fmt.Println(tryAgainMsg)
			continue
		}
		return nodeID
	}
	fmt.Println("Too many invalid moves. Passing.")
	return -1
}

// runGame expects an initialized boardState (GoGraph populated)
// and runs a local CLI game.
func (gameState *boardState) runGame(isRandom bool) {
	coo, nodeCount := gameState.coords, gameState.nodeCount
	prompt := promptRect()
	gameState.ongoing = true

	for gameState.ongoing {
		nodeID := promptNodeID(prompt, coo, nodeCount)
		// check for exit code
		if nodeID == -2 {
			print("Goodbye!")
			break
		}
		// instantiates an instance of the moveInput
		// struct with id equal to user input
		next := moveInput{id: nodeID}

		// Populate playerColor
		if gameState.whiteToMove == false {
			next.playerColor = 1 // black
		} else {
			next.playerColor = 2 // white
		}

		if nodeID == -1 {
			fmt.Println("Passing your turn.")
			next.isPass = true // Must be changed for passes to be counted
		}

		//The next moveInput is fully populated.
		if err := gameState.playMoveInput(next); err != nil {
			fmt.Println(err)
			fmt.Println("Please try again.")
		} else {
			fmt.Printf("%v", gameState.GoGraph)

		}

		// This should force an image of the board to come up.
		visualizeBoard(gameState.GoGraph, isRandom)
	}
	whiteNaive, blackNaive := getNaiveScore(gameState.GoGraph)
	fmt.Printf("Final score: \nWhite: %v\nBlack: %v\n", float64(gameState.whitePoints)+whiteNaive, float64(gameState.blackPoints)+blackNaive)
}

// moveByID is an alternate way to run a game. It is not interactive.
// You simply provide the node ID and the move is played.
func (gameState *boardState) moveByID(ID int) error {
	gameState.ongoing = true
	if ID == -2 {
		return errors.New("moveByID does not support this exit code")
	}
	if ID < -2 {
		return errors.New("Non-pass negative input given to moveByID")
	}
	if N := gameState.nodeCount; ID >= N {
		errorMsg := fmt.Sprintf("ID %v is too large. Max valid ID: %v", ID, N)
		return errors.New(errorMsg)
	}

	next := moveInput{id: ID} // Populate move: ID, color, pass.
	if gameState.whiteToMove == false {
		next.playerColor = 1
	} else {
		next.playerColor = 2
	}
	if ID == -1 {
		next.isPass = true // Must be changed for passes to be counted
	}
	return gameState.playMoveInput(next)
}

// StartRectangularGame initializes an n-by-m board and runs a CLI game.
func StartRectangularGame(n int, m int) {
	state := initBoardState(makeSquareBoard(n, m), 6) //6 komi
	isRandom := false

	state.runGame(isRandom)
	return
}

func StartRandomGame(n int, prob float64) {
	state := initBoardState(makeRandomBoard(n, prob), 10) // 10 komi

	isRandom := true
	state.runGame(isRandom)
}
