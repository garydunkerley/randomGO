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

// parseIntOrCoords converts a validated input (int or coords) into a nodeId
// and a validation bool.
// It expects a map of valid coordinate pairs and the number of nodes for the given board.
// If the input is coordinate-form, then the second return value is true when
// the node is a key for the given coordinate map.
// If the input is integer-form, we check if it's in the range -2 <= n < nodeCount.
// The constant -2 is due to negative numbers coding special instructions:
// -2 is exit, -1 is pass.
func parseIntOrCoords(input string, coo map[[2]int]int, nodeCount int) (nodeId int, validId bool) {
	if err := checkIntInput(input); err != nil {
		// Coordinate case: get the node id from the given coordinate map.
		if err2 := checkCoordsInput(input); err2 != nil {
			panic("Invalid input passed to parseIntOrCoords, validate first!")
		}
		a, b := parseCoords(input)
		nodeId, validId = coo[[2]int{a, b}]
		return nodeId, validId
	} // Integer case: check if the node id is a key for the given edge map.
	nodeId = parseInt(input)
	validId = (-2 <= nodeId) && (nodeId < nodeCount)
	return nodeId, validId
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

// promptNodeId applies the given prompt, validates input, converts to a node id.
// If a certain number of consecutive invalid inputs are given, it passes for you.
// It only checks if the node id exists, not if the move is legal.
func promptNodeId(prompt promptui.Prompt, coo map[[2]int]int, nodeCount int) int {
	tryAgainMsg := "Try again, or enter -2 to exit."
	for attempts := 5; attempts > 0; attempts-- {
		inputString, err := prompt.Run()
		if err != nil { // malformed input string
			fmt.Println(err)
			fmt.Println(tryAgainMsg)
			continue
		} // good input, convert to a meaningful int
		nodeId, ok := parseIntOrCoords(inputString, coo, nodeCount)
		if !ok {
			fmt.Println("Invalid node ID.")
			fmt.Println(tryAgainMsg)
			continue
		}
		return nodeId
	}
	fmt.Println("Too many invalid moves. Passing.")
	return -1
}

//TODO: initialize game function, boardtop > boardstate. Initialize history in this function.

//RunGame expects an initialized boardState and runs a local CLI game.
func (gameState boardState) RunGame() {
	coo, nodeCount := gameState.coords, gameState.nodeCount
	prompt := promptRect()

	for nodeId := promptNodeId(prompt, coo, nodeCount); nodeId != -2 && gameState.ongoing; {
		next := moveInput{id: nodeId}

		// Populate playerColor
		if gameState.whiteToMove == false {
			next.playerColor = 2
		} else {
			next.playerColor = 1
		}

		if nodeId == -1 {
			next.isPass = true
			fmt.Println("Passing your turn.")
		}
		//The next moveInput is fully populated.

		if ok := gameState.playMoveInput(next); !ok { //print error if any
			fmt.Println(ok)
			fmt.Println("Please try again.")
		}
		fmt.Println("The number of boardstates in the board history is:", len(game_state.board_history))
	}
	fmt.Println("Game over. You lose.") // a little bit rigged
}
