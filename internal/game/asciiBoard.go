package game

import "strings"

func (G GoGraph) String() string {
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
