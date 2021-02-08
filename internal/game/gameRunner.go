package game

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

func (myBoard GoGraph) RunGame() {
	// RunGame takes the current state of the board and, if it determines the game is on-going, requests and processes player input in order to act on the board state.

	var game_state boardState
	// game_State contains the colored graph representing the board, tracks whose turn it is, and contains the history of the board

	var consecutive_passes int
	// consecutive_passes is a value that increases each time a player passes and resets to zero if a player makes a non-passing move

	var valid_id bool //used to avoid code repetition later on
	var node_id int   //store player move only
	game_state.current_board = myBoard
	game_state.board_history = append(game_state.board_history, game_state.current_board)
	// make it so the empty board is listed as the starting state in the board history.

	//Construct the prompt here: if we make too many prompts, errors start popping up.
	//By defining the prompt inside the run_Game function, we only create one prompt.
	//move is checked for validity but not legality
	validator := func(input string) error {
		string_pieces := strings.Split(input, "-")
		if len(string_pieces) > 2 { //Case: a-b-c input given.
			return errors.New("Too many coordinates given.")
		} else if len(string_pieces) == 1 { //Case: integer input given.
			_, err1 := strconv.Atoi(input)
			if err1 != nil { //parse error with our 1 int
				return errors.New("Input parsing error: cannot parse int.")
			}
			return nil
		} else if len(string_pieces) == 2 { //Case: either 2 coordinates given, or a negative number.
			_, err1 := strconv.Atoi(string_pieces[0])
			if err1 != nil {
				if string_pieces[0] != "" {
					//This isn't a negative number, just a bad first coordinate.
					return errors.New("Input parsing error: cannot parse your first coordinate.")
				}
			}
			_, err2 := strconv.Atoi(string_pieces[1])
			if err2 != nil {
				return errors.New("Input parsing error: cannot parse your second coordinate.")
			}
			return nil
		}
		return errors.New("Unimplemented coordinate case reached.")
	}

	prompt := promptui.Prompt{
		Label:    "Move (int)",
		Validate: validator,
	}

	//check_coords("i-j") = true, i, j
	//check_coords("23") = false, 0, 0
	//bool indicates if the string represents dash-separated coordinates
	//ints are the coordinates, if applicable
	check_coords := func(input string) (bool, int, int) {
		string_pieces := strings.Split(input, "-")
		if len(string_pieces) != 2 {
			return false, 0, 0
		} else if string_pieces[0] == "" { //empty string case
			return false, 0, 0
		} else {
			i, err1 := strconv.Atoi(string_pieces[0])
			j, err2 := strconv.Atoi(string_pieces[1])
			if err1 != nil {
				panic(err1)
			} else if err2 != nil {
				panic(err2)
			} else {
				return true, i, j
			}
		}
	}

	for consecutive_passes < 2 && game_state.move_count < 50 {
		valid_id = false //require the move id to be verified each loop iteration
		node_id = -2

		if game_state.white_to_move == false {
			fmt.Println("It is black's move. Type a node ID between 0 and ",
				len(game_state.current_board.nodes)-1, "or -1 to pass. Exit with -2.")
		} else {
			fmt.Println("It is white's move. Type a number between 0 and ",
				len(game_state.current_board.nodes)-1, "or -1 to pass. Exit with -2.")
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
				id, ok := game_state.current_board.coords[coord_input]
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
					_, ok := game_state.current_board.nodes[node_id]
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
								game_state.board_history, game_state.current_board)
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
				game_state.board_history = append(game_state.board_history, game_state.current_board)
				game_state.move_count += 1
				fmt.Println(myBoard)
			}
		}
		fmt.Println("The number of boardstates in the board history is:", len(game_state.board_history))
	}
	fmt.Println("Game over. You lose.") // a little bit rigged
}
