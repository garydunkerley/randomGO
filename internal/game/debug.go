package game

import "fmt"

func DebugPrint() {

	fmt.Println("\ncheckIntInput: %v\n", checkIntInput("-2"))
	fmt.Println("\ncheckIntInput: %v\n", checkIntInput("2"))
	fmt.Println("\ncheckIntInput: %v\n", checkIntInput("2-2"))

	fmt.Println("\n\ncheckCoordsInput: %v\n", checkCoordsInput("-2"))
	fmt.Println("\ncheckCoordsInput: %v\n", checkCoordsInput("2"))
	fmt.Println("\ncheckCoordsInput: %v\n", checkCoordsInput("2-2"))

	fmt.Println("\nparseCoords\n")
	// causes panic
	// x, z := parseCoords("-2")
	// fmt.Printf("Parse -2:\n nodeID: %v, valid: %v\n", x, z)

	// x, z := parseCoords("2")
	// fmt.Printf("Parse 2:\n nodeID: %v, valid: %v\n", x, z)

	x, z := parseCoords("2-2")
	fmt.Printf("Parse 2-2:\n nodeID: %v, valid: %v\n", x, z)

	fmt.Println("\nparseIntOrCoords\n")
	w, y := parseIntOrCoords("-2", nil, 100)
	fmt.Printf("Parse -2:\n nodeID: %v, valid: %v\n", w, y)

	w, y = parseIntOrCoords("2", nil, 100)
	fmt.Printf("Parse 2:\n nodeID: %v, valid: %v\n", w, y)

	J := make(map[[2]int]int)
	J[[2]int{2, 2}] = 20
	w, y = parseIntOrCoords("2-2", J, 25)
	fmt.Printf("Parse 2-2:\n nodeID: %v, valid: %v\n", w, y)
}
