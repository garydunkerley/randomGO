package game

import "strings"

func (G GoGraph) String() string {
	//will edit these in the core loop
	var currentLine string
	var ij [2]int
	var color int8
	var nextChar string
	boardLines := make([]string, 1)
	boardLines[0] = "\n" //add a blank line before printing board

	//This should only try to format dim-2 graphs.
	if len(G.coordBounds) != 2 {
		return ""
	}

	iBound := G.coordBounds[0]
	jBound := G.coordBounds[1]

	for j := 0; j < jBound; j++ {
		currentLine = "" //reset the line we're working on
		ij[1] = j        //set j coordinate
		for i := 0; i < iBound; i++ {
			ij[0] = i //set i coordinate

			//At each step: redefine the current line by adding next character.
			//At the moment: use a space for missing node, B for black, W for white, . for empty.
			//Hopefully unused: will use @ for mysterious error character

			nodeID, ok := G.coords[ij]
			if !ok {
				nextChar = " " //missing node, use blank space
			} else {
				color = G.nodes[nodeID].color
				switch color {
				case 0:
					nextChar = "." //empty
				case 1:
					nextChar = "B" //black
				case 2:
					nextChar = "W" //white
				default:
					nextChar = "@" //anomalous, occurs if the color is bad.
				}
			}
			currentLine = currentLine + nextChar
		}
		boardLines = append(boardLines, currentLine)
	}
	//now merge the slice of strings into a single string
	//using the newline character to separate
	return strings.Join(boardLines, "\n") + "\n"
}
