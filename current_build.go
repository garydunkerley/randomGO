package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

type node struct {
	name      string
	id        int
	neighbors []*node
	color     int8        // color is 0 for empty, 1 for black, 2 for white
	inGroup   *stoneGroup // is a part of a given stone group
}

type stoneGroup struct {
	//	id      int        // in case we want to refer to them?
	stones    map[*node]bool // what stones are in here? Need to make all empty if group dies
	liberties map[*node]bool // what liberties does the group have
	// https://softwareengineering.stackexchange.com/questions/177428/sets-data-structure-in-golang
}

type goGraph struct {
	nodes        map[int]*node
	coords       map[[2]int]int
	coord_bounds []int //for example, [9,9] for a 9 by 9 board
}

type boardState struct {
	current_board goGraph
	white_to_move bool
	move_history  []string
	board_history []goGraph // to determine if move violates ko rule
	move_count    int
}

// Use this to establish the neighbors of a given node
func (x *node) set_neighbors(y []*node) {
	x.neighbors = y
	return
}

// Determine if a move is legal or not.`
func (y boardState) check_legal_move(n int) error {
	node, ok := y.current_board.nodes[n]
	if !ok {
		return errors.New("Illegal move: nonexistent node.")
	}
	if node.color != 0 {
		return errors.New("Illegal move: Nonempty point.")
	}
	// The move is legal if:
	// 1. No ko rule. (not implemented)
	// 2. No suicide.
	// To check suicide, it suffices to:
	// A. Check if this point is the last liberty for an opponent's group.
	// B. Check if this point is the last liberty for one of your own groups.
	if y.killing_move(n) {
		return nil
	} else if y.suicidal_move(n) {
		return errors.New("Illegal move: Suicide.")
	} else if y.illegal_ko_move(n) {
		return errors.New("Illegal move: Violates ko rule.")
	} else {
		return nil
	}
}

//Syntax for these next few:
//x.killing_move(true) asks "if white plays at x, does it kill a black group?"
//x.killing_move(false) asks "if black plays ... a white group?"
//Can assume that this will only be called for an empty intersection x.

func (y boardState) killing_move(n int) bool {
	//TODO:
	// for each z neighboring y.game_state.nodes[n]
	// 	if z.color != y.game_state.nodes[n].color
	//		if len(y.roup.liberties) == 1
	//			return true
	return true
}

func (y boardState) suicidal_move(n int) bool {
	//TODO:
	// For each y neighboring x
	//	if y.color == x.color
	//		if len(y.group.liberties) == 1
	//			return true
	return false
}

func (y boardState) illegal_ko_move(n int) bool {
	//TODO:
	// for i in move_history
	// 	if i = x   // check if move has been played before
	//		if (board state corresponding to move i)= gograph
	//			return true
	return false
}

// The next three functions could be consolidated into a single function
// that returns nodes of a given color

func get_same_color_neighbors(x *node) []*node {
	var friends []*node
	if x.color == 0 {
		return friends
	}
	for _, y := range x.neighbors {
		if y.color == x.color {
			friends = append(friends, y)
		}
	}
	return friends
}

func get_opp_color_neighbors(x *node) []*node {
	var enemies []*node
	if x.color == 0 {
		return enemies
	}
	for _, y := range x.neighbors {
		if y.color != x.color && y.color != 0 {
			enemies = append(enemies, y)
		}
	}
	return enemies
}

func get_liberties(x *node) []*node {
	var liberties []*node
	if x.color == 0 {
		return liberties
	} else {
		for _, y := range x.neighbors {
			if y.color == 0 {
				liberties = append(liberties, y)
			}
		}
	}
	return liberties
}

// added to gameplay/gamestate.go

func modify_stoneGroups(y *node) {
	var friendlies []*node
	var enemies []*node

	var old_Group *stoneGroup
	var merged_Group *stoneGroup

	newStones := make(map[*node]bool)
	newLibs := make(map[*node]bool)
	newNodeLibs := get_liberties(y)

	friendlies = get_same_color_neighbors(y)
	enemies = get_opp_color_neighbors(y)

	if len(friendlies) > 0 {
		// joins together neighboring friendly groups
		merged_Group = friendlies[0].inGroup

		for i := 0; i < len(friendlies); i++ {
			old_Group = friendlies[i].inGroup
			// set old_Group to the group of the ith friendly node
			for key := range old_Group.liberties {
				if key != y {
					// foreach liberty of the old stone group, set it as
					// a liberty of the new one
					newLibs[key] = true
				}
			}

			for key := range old_Group.stones {
				newStones[key] = true
				key.inGroup = merged_Group
				// should change each of these nodes so that they now
				// point to the group for y.
			}

		}
		for _, key := range newNodeLibs {
			newLibs[key] = true
		}
		newStones[y] = true

		merged_Group.stones = newStones // update the stones and liberties of y's group
		merged_Group.liberties = newLibs
		y.inGroup = merged_Group
	} else {
		// if y is not next to any friendly stones
		// make a group containing y
		var new_group stoneGroup
		newStones[y] = true
		for _, z := range get_liberties(y) {
			newLibs[z] = true
		}
		new_group.stones = newStones

		new_group.liberties = newLibs

		y.inGroup = &new_group
	}
	if len(enemies) > 0 {
		for _, z := range enemies {
			delete(z.inGroup.liberties, y)
		}
	}
}

// removes stones from groups without liberties
// MOVED TO gamestate.go
func (x boardState) remove_dead(y *node) error {
	var enemies []*node
	enemies = get_opp_color_neighbors(y)
	if len(enemies) == 0 {
		return nil
	} else {
		for i := 0; i < len(enemies); i++ {
			if len(enemies[i].inGroup.liberties) == 0 {
				// if enemy group has no liberties
				for z := range enemies[i].inGroup.stones {
					z.color = 0
					// change all stones into empty nodes
					delete(enemies[i].inGroup.stones, z) // remove each stone from the map
				}
			}
		}
	}
	return nil
}

// function that allows us play a move
func (x boardState) make_move(y int) {
	var stoneColor int8
	if x.white_to_move {
		stoneColor = 2
	} else {
		stoneColor = 1
	}
	node, ok := x.current_board.nodes[y]
	if !ok {
		error_message := fmt.Sprintf("Node %d not found!", y)
		panic(error_message)
	}
	node.color = stoneColor
	modify_stoneGroups(x.current_board.nodes[y])
	x.remove_dead(x.current_board.nodes[y])
	modify_stoneGroups(x.current_board.nodes[y])
	// I run modify_stoneGroups again to append any liberties resulting from stone groups dying
	// TODO clean up this inefficiency
}

func contains_val(slice []*node, val *node) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func initializeNode(i int, a string) *node {
	// initializes empty nodes whose ids are integers

	var x node
	x.name = a
	x.id = i
	x.color = 0
	return &x
}

func makeGoban(n int, m int) goGraph {
	var name string
	goban_map := make(map[int]*node)
	coord_map := make(map[[2]int]int)
	var ourGraph goGraph
	var newNode *node
	//initialize all nodes

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			//TODO: initialize a node named "i-j"
			name = strconv.Itoa(i) + "-" + strconv.Itoa(j)
			fmt.Println(name)

			ourSlice := [2]int{i, j}
			newNode = initializeNode(i*n+j, name)
			goban_map[i*n+j] = newNode
			coord_map[ourSlice] = i*n + j

		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j != 0 { // add left node if it exists
				goban_map[i*n+j].neighbors = append(goban_map[i*n+j].neighbors, goban_map[i*n+j-1])
			}
			if j != m-1 { // add right node if it exists
				goban_map[i*n+j].neighbors = append(goban_map[i*n+j].neighbors, goban_map[i*n+j+1])
			}
			if i != 0 { // add lower node if it exists
				goban_map[i*n+j].neighbors = append(goban_map[i*n+j].neighbors, goban_map[(i-1)*n+j])
			}
			if i != n-1 { // add upper node if it exists
				goban_map[i*n+j].neighbors = append(goban_map[i*n+j].neighbors, goban_map[(i+1)*n+j])
			}

		}
	}
	ourGraph.nodes = goban_map
	ourGraph.coords = coord_map
	ourGraph.coord_bounds = []int{n, m}
	return ourGraph
}

func (myBoard goGraph) runGame() {
	var game_state boardState
	var consecutive_passes int
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
			if err := game_state.check_legal_move(node_id); err != nil {
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
	}
	fmt.Println("Game over. You lose.") // a little bit rigged
}

//Define a way to print our graph in a nice way.
func (G goGraph) String() string {
	//will edit these in the core loop
	var current_line string
	var ij [2]int
	var color int8
	var next_char string
	board_lines := make([]string, 1)
	board_lines[0] = "\n" //add a blank line before printing board

	//This should only try to format dim-2 graphs.
	if len(G.coord_bounds) != 2 {
		return ""
	}

	i_bound := G.coord_bounds[0]
	j_bound := G.coord_bounds[1]

	for j := 0; j < j_bound; j++ {
		current_line = "" //reset the line we're working on
		ij[1] = j         //set j coordinate
		for i := 0; i < i_bound; i++ {
			ij[0] = i //set i coordinate

			//At each step: redefine the current line by adding next character.
			//At the moment: use a space for missing node, B for black, W for white, . for empty.
			//Hopefully unused: will use @ for mysterious error character

			node_id, ok := G.coords[ij]
			if !ok {
				next_char = " " //missing node, use blank space
			} else {
				color = G.nodes[node_id].color
				switch color {
				case 0:
					next_char = "." //empty
				case 1:
					next_char = "B" //black
				case 2:
					next_char = "W" //white
				default:
					next_char = "@" //anomalous, occurs if the color is bad.
				}
			}
			current_line = current_line + next_char
		}
		board_lines = append(board_lines, current_line)
	}
	//now merge the slice of strings into a single string
	//using the newline character to separate
	return strings.Join(board_lines, "\n") + "\n"
}

func main() {
	theBoard := makeGoban(5, 5)
	theBoard.runGame()
}
