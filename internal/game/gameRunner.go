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
	string_pieces := strings.Split(input, "-")
	if len(string_pieces) != 2 {
		return errors.New("Input is not formatted as coordinates.")
	} else if string_pieces[0] == "" { //empty string case: input is "-m"
		return errors.New("Input is formatted as a negative int.")
	} else {
		_, err1 := strconv.Atoi(string_pieces[0])
		_, err2 := strconv.Atoi(string_pieces[1])
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
	string_pieces := strings.Split(coordInput, "-")
	a, _ := strconv.Atoi(string_pieces[0])
	b, _ := strconv.Atoi(string_pieces[1])
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
	if int_err := checkIntInput(input); int_err != nil {
		if coords_err := checkCoordsInput(input); coords_err != nil {
			return errors.New("Input cannot be parsed as an integer or coordinate pair.")
		}
	}
	return nil
}

// parseIntOrCoords converts a validated input (int or coords) into a node_id
// and a validation bool.
// If the input is coordinate-form, then the second return value is true when
// the node is a key for the given coordinate map.
// If the input is integer-form, we check if it's in the range -2 <= n < node_count.
// The constant -2 is due to negative numbers coding special instructions:
// -2 is exit, -1 is pass.
func parseIntOrCoords(input string, coords map[[2]int]int, node_count int) (node_id int, valid_id bool) {
	if err := checkIntInput(input); err != nil {
		// Coordinate case: get the node id from the given coordinate map.
		if err2 := checkCoordsInput(input); err2 != nil {
			panic("Invalid input passed to parseIntOrCoords, validate first!")
		}
		a, b := parseCoords(input)
		node_id, valid_id = coords[[2]int{a, b}]
		return node_id, valid_id
	} // Integer case: check if the node id is a key for the given edge map.
	node_id = parseInt(input)
	valid_id = (-2 <= node_id) && (node_id < node_count)
	return node_id, valid_id
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
func promptNodeId(p promptui.Prompt) int {
	try_again_msg := "Try again, or enter -2 to exit."
	for attempts := 5; attempts > 0; attempts-- {
		input_string, err := prompt.Run()
		if err != nil { // malformed input string
			fmt.Println(err)
			fmt.Println(try_again_msg)
			continue
		} // good input, convert to a meaningful int
		node_id, ok := parseIntOrCoords(input_string)
		if !ok {
			fmt.Println("Invalid node ID.")
			fmt.Println(try_again_msg)
			continue
		}
		return node_id
	}
	fmt.Println("Too many invalid moves. Passing.")
	return -1
}

//TODO: divorce boardState update logic from UX
func (gameState boardState) RunGame() {
	myBoard := gameState.GoGraph
	// var valid_id bool //used to avoid code repetition later on
	// var node_id int   //store player move only

	// TODO: initialize board history
	// game_state.board_history = append(game_state.board_history, game_state.GoGraph)
	// make it so the empty board is listed as the starting state in the board history.

	// Construct the prompt here: if we make too many prompts, errors start popping up.
	// By defining the prompt inside the run_Game function, we only create one prompt.
	prompt := promptRect()

	// TODO: abstract the game-over condition into boardState or equivalent
	// TODO: use prompt functions written above to rewrite core loop interior
	for consecutive_passes < 2 && game_state.move_count < 50 {
		valid_id = false //require the move id to be verified each loop iteration
		node_id = -2

		if game_state.white_to_move == false {
			fmt.Println("It is black's move. Type a node ID between 0 and ",
				len(game_state.nodes)-1, "or -1 to pass. Exit with -2.")
		} else {
			fmt.Println("It is white's move. Type a number between 0 and ",
				len(game_state.nodes)-1, "or -1 to pass. Exit with -2.")
		}

		//This huge mess is not complex.
		//First it checks if the input is a valid coordinate pair.
		//If not, it checks if the input is a valid integer, and handles
		//special integer codes (atm just -2 for quit, -1 for pass) separately.

		//Afterward, if valid_id==true, we try to play the move node_id.
		try_again_msg := "That didn't work. Try again, or enter -2 to exit."
		input_string, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			fmt.Println(try_again_msg)
		} else {
			isCoord, _, _ := check_coords(input_string)
			if isCoord { //handle coordinate input case: set valid_id and node_id
				_, i, j := check_coords(input_string)
				coord_input := [2]int{i, j}
				id, ok := game_state.coords[coord_input]
				if !ok {
					fmt.Println("Your coordinates don't match any node.")
					fmt.Println(try_again_msg)
				} else {
					valid_id = true
					node_id = id
				}
			} else { //integer input case
				_, err := strconv.Atoi(input_string)
				if err != nil {
					fmt.Println(err)
					fmt.Println(try_again_msg)
					valid_id = false
				} else {
					node_id, _ = strconv.Atoi(input_string)
					_, ok := game_state.nodes[node_id]
					if ok {
						valid_id = true
					} else {
						switch node_id { //can add extra negative int codes here
						case -2: //exit code
							fmt.Println("Goodbye!")
							return
						case -1: //pass
							fmt.Println("Passing your turn.")
							consecutive_passes++
							//game_state.make_move(casecode)
							//don't make a move with invalid codes
							game_state.white_to_move = !game_state.white_to_move
							game_state.board_history = append(
								game_state.board_history, game_state.asdf) //TODO fix 1
							game_state.move_count += 1
						default:
							fmt.Println(try_again_msg)
						}
					}
				}
			}
		}
		if valid_id {
			if err := game_state.check_legal_move(node_id, game_state.white_to_move); err != nil {
				fmt.Println(err)
				fmt.Println(try_again_msg)
			} else {
				fmt.Println("GOOD MOVE!")
				consecutive_passes = 0
				game_state.make_move(node_id) //adjust node color, remove dead stones
				game_state.white_to_move = !game_state.white_to_move
				game_state.board_history = append(game_state.board_history, game_state.asdf) //TODO fix 2
				game_state.move_count += 1
				fmt.Println(myBoard)
			}
		}
		fmt.Println("The number of boardstates in the board history is:", len(game_state.board_history))
	}
	fmt.Println("Game over. You lose.") // a little bit rigged
}
