package game

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// makeNaiveRandomBoard creates a board topology on n nodes such that
// any two nodes are connected via an edge with a fixed probability.
// This is decrepitated by the newer function that makes baords that
// are planar and more likely to be fun to play on.
func makeNaiveRandomBoard(n int, prob float64) boardTop {
	var ourTopology boardTop

	edges := make(map[int][]int)

	ourTopology.nodeCount = n

	for i := 0; i < n; i++ {
		if i < n-1 {
			edges[i] = append(edges[i], i+1)

			for j := i + 2; j < n; j++ {
				rand.Seed(time.Now().UnixNano())
				v := rand.Float64()

				if v < prob {
					edges[i] = append(edges[i], j)
					edges[j] = append(edges[j], i)
				}
			}
		}
		if i > 0 {
			edges[i] = append(edges[i], i-1)
		}

	}

	ourTopology.edges = edges

	return ourTopology
}

// sum just sums the elements of an integer array
func sum(array []int) int {
	result := 0

	for _, v := range array {
		result += v
	}

	return result
}

// euclideanNorm gets the magnitude of a point as a vector in two-dimensional Euclidean space
func euclideanNorm(x []float64) float64 {
	var sum float64
	for i := 0; i < len(x); i++ {
		sum += math.Pow(x[i], 2)
	}
	return math.Sqrt(sum)
}

func vecSum(x []float64, y []float64) []float64 {
	var z []float64
	for i := 0; i < len(x); i++ {
		z = append(z, x[i]+y[i])
	}
	return z
}

// getHexagonalNumbers computes a slice of centered hexagonal numbers less than or equal to n
func getCenteredHexagonalNumbers(n int) []int {
	var hexNums []int
	// var newHexNums []int

	var newHex int

	count := 0
	searching := true

	for searching {

		newHex = 3*count*(count+1) + 1

		if newHex <= n {
			hexNums = append(hexNums, newHex)
			count += 1
		} else {
			searching = false
		}
	}

	return hexNums
}

// getHexagonalLattice takes an integer and outputs the edges and nodes of a
// normal hexagonal lattice in two-dimensional Euclidean space
func getHexagonalLattice(n int) (map[int][]int, map[int][]float64) {

	theta := math.Pi / 6
	multiplier := float64(0)

	edges := make(map[int][]int)
	coordMap := make(map[int][]float64)

	// roomsMade will track how far we have iterated along
	// direction to create a layer in Pascal's triangle.
	cellsMade := 1

	// counter tracks the number of nodes
	// that have been generated so far
	counter := 1

	// the first node will be centered at the coordinate 0,0

	currentVec := []float64{0, 0}
	coordMap[0] = currentVec

	// we will also initialize the direction of generation
	// as being trivial
	genDirection := []float64{0, 0}

	hexNums := getCenteredHexagonalNumbers(n)
	hexIterate := 0

	corner := 1
	cornerCount := 1

	fmt.Println("Our hexNumbers: ", hexNums)
	for counter < n {
		levelingNumber := int(math.Min(float64(hexIterate+1), float64(len(hexNums)-1)))
		if counter == hexNums[levelingNumber] {
			// move along the current direction to the
			// node that was seen previously
			currentVec = vecSum(currentVec, genDirection)
			multiplier = 0
			genDirection = []float64{1, 0}

			if hexIterate < len(hexNums)-1 {
				hexIterate += 1
			}

		}

		newVec := vecSum(currentVec, genDirection)
		coordMap[counter] = newVec

		// After establishing the existence of a new node,
		// we attach edges between it and nodes we have already created.

		// if counter happens to be a hexagonal number and greater than 1,
		// we ensure there is an edge between it and the previous hexagonal number.

		// we distinguish three types of numbers, corners (numbers whose distance)
		// from a hexagonal number is a multiple of the hexIterate+1), nodes that are
		// next to hexNumbers (and so are subject to special index management), and the rest.

		// if a number is a corner
		if counter == corner {

			if counter != hexNums[hexIterate] {
				edges[counter] = append(edges[counter], counter-1)
				edges[counter-1] = append(edges[counter-1], counter)
			}

			edges[counter] = append(edges[counter], counter-cornerCount)
			edges[counter-cornerCount] = append(edges[counter-cornerCount], counter)

			corner += (hexIterate + 1)

			if counter+1 != hexNums[levelingNumber] {
				cornerCount += 1
			}
			if counter+1 == 7 {
				edges[counter] = append(edges[counter], 1)
				edges[1] = append(edges[hexNums[hexIterate]], counter)

				edges[counter] = append(edges[counter], counter-1)
				edges[counter-1] = append(edges[counter-1], counter)

			}

		} else if counter+1 == hexNums[levelingNumber] {
			edges[counter] = append(edges[counter], hexNums[hexIterate])
			edges[hexNums[hexIterate]] = append(edges[hexNums[hexIterate]], counter)

			edges[counter] = append(edges[counter], hexNums[hexIterate-1])
			edges[hexNums[hexIterate-1]] = append(edges[hexNums[hexIterate-1]], counter)

			edges[counter] = append(edges[counter], counter-1)
			edges[counter-1] = append(edges[counter-1], counter)

			second := counter - (cornerCount)

			edges[counter] = append(edges[counter], second)
			edges[second] = append(edges[second], counter)

		} else {

			first := counter - (cornerCount)
			second := counter - (cornerCount - 1)

			edges[counter] = append(edges[counter], first, second)
			edges[first] = append(edges[first], counter)
			edges[second] = append(edges[second], counter)

			if counter > 1 {
				edges[counter] = append(edges[counter], counter-1)
				edges[counter-1] = append(edges[counter-1], counter)
			}
		}

		cellsMade += 1

		if cellsMade == hexIterate {

			multiplier = float64((int(multiplier) + 1) % 6)
			genDirection = []float64{math.Sin(multiplier * theta), math.Cos(multiplier * theta)}
			cellsMade = 0

		}

		counter += 1
	}

	for i := 0; i < n; i++ {
		slice := edges[i]
		ascendingSlice := makeAscending(slice)
		edges[i] = ascendingSlice
		fmt.Println(i, " is sent to ", ascendingSlice)
	}

	return edges, coordMap
}

// makeAscending rearranges an integer slice so that its elements are in ascending order
func makeAscending(mySlice []int) []int {

	targetLength := len(mySlice)

	var ascendingSlice []int

	for len(ascendingSlice) < targetLength {
		currentMinIndex := 0
		// iterate over the elements of your slice
		// when you find an element smaller than the
		// current minimum, record it and move on.
		for i := 0; i < len(mySlice); i++ {
			if mySlice[i] <= mySlice[currentMinIndex] {
				currentMinIndex = i
			}
		}
		// add the smallest element to the ascending slice
		// and then remove it from the original
		ascendingSlice = append(ascendingSlice, mySlice[currentMinIndex])
		mySlice = append(mySlice[:currentMinIndex], mySlice[currentMinIndex+1:]...)
	}

	return ascendingSlice
}

// getCircuit takes our collection of edges and outputs an
// encoding of a spanning tree that is also a path
func getCircuit(n int, edges map[int][]int) map[int][]int {

	var circuit []int

	circuitMap := make(map[int][]int)
	visited := make(map[int]bool)

	currentNode := 0

	// look at the list of nodes adjacent to the current node arranged in
	// increasing order and then append the smallest, unvisited node
	// to the circuit

	// TODO: remove counter for bug testing
	counter := 0
	for len(visited) < n && counter < 1 {

		ascendingEdges := makeAscending(edges[currentNode])

		fmt.Println("Debug (makeRandomBoard l.235): ascending array created")

		fmt.Println("Debug (makeRandomBoard l.40): the length of the ascending array for node", currentNode, "is", len(ascendingEdges))

		for _, i := range ascendingEdges {
			if !(visited[i]) {
				circuit = append(circuit, currentNode)
				currentNode = ascendingEdges[i]
				visited[currentNode] = true
				break
			}
		}
		counter += 1
	}

	// build a map that encode the incidence behavior of the nodes forming the circuit.
	fmt.Println("Debug (makeRandomBoard l.245): we are aboutto construct the circuit map")
	for i := range circuit {
		if i > 0 && i < n-1 {

			circuitMap[circuit[i]] = append(circuitMap[circuit[i]], circuit[i-1], circuit[i+1])
		} else if i == 0 {
			circuitMap[circuit[i]] = append(circuitMap[circuit[i]], circuit[i+1])
		} else {
			circuitMap[circuit[i]] = append(circuitMap[circuit[i]], circuit[i-1])
		}
	}

	return circuitMap
}

// sparsifyEdges uses a random process to remove edges from the edge map provided.
// The process looks at how far the associated coordinate is from the center and uses this distance
// to determine the probability that the node will lose and edge and which one.
// A particular spanning tree will be protected to ensure that the resulting graph is always connected
func sparsifyEdges(n int, edges map[int][]int, coordMap map[int][]float64) map[int][]int {

	safeEdges := make(map[int][]int)

	// declare all edges in the circuit to be safe from elimination
	// to ensure that the sparsified graph is always connected.
	circuit := getCircuit(n, edges)
	fmt.Println("Debug (makeRandomBoard l.276): our circuit is:", circuit)
	for i := 0; i < n; i++ {
		safeEdges[i] = circuit[i]
	}
	fmt.Println("Debug (makeRandomboard l.279): circuit has been constructed")

	for i := 0; i < n; i++ {

		var validEdges []int

		// iterate over all the edges attached to a node i and check to see if they are
		// safe (i.e. that they belong to the circuit computed earlier or they have
		// survived a roll already)
		// edges that are not safe may be subject to elimination when the rolling ovvurs.
		for _, edge := range edges[i] {
			isSafe := false
			for _, safeEdge := range safeEdges[i] {
				if edge == safeEdge {
					isSafe = true
					break
				}
			}
			if !isSafe {
				validEdges = append(validEdges, edge)
			}
		}

		fmt.Println("Debug: the number of valid edges is", len(validEdges))

		// if a node has more than 2 edges, then its edges will be subjected to
		// a random elimination process
		if len(validEdges) > 2 {
			//norm := euclideanNorm(coordMap[i])

			stillRolling := true
			for stillRolling || len(validEdges) > 2 {
				//TODO: tweak this until I get something I like.
				// Right now, being farther away from the center increases the likelihood that
				// a node will lose and edge as does having a lot of edges

				// prob := math.Pi / (2 * math.Pow(math.Atan(norm+0.05), float64(6-len(validEdges))))
				prob := float64(0)
				rand.Seed(time.Now().UnixNano())
				v := rand.Float64()

				if v < prob {
					randomIndex := rand.Intn(len(validEdges))
					validEdges = append(validEdges[:randomIndex], validEdges[randomIndex+1:]...)

				} else {
					stillRolling = false
				}

			}

			// all edges that survive this process are declared to be safe
			if len(validEdges) > 0 {
				for j := 0; j < len(validEdges)-1; j++ {

					safeEdges[validEdges[j]] = append(safeEdges[validEdges[j]], i)
					safeEdges[i] = append(safeEdges[i], validEdges[j])
				}
			}
		}
	}

	return safeEdges
}

/*
func addBoundaryEdges(n int, edges map[int][]int) map[int][]int {
//TODO make a map that will be guaranteed a planar representation.

	hexNums := getHexagonalNumbers(n)

	boundary = []int

	// TODO how to find the boundary?

return
}
*/

// makeRandomBoard begins by generating a hexagonal lattice on n points
// and then runs through a procedure to make edges more sparse
func makeRandomBoard(n int) boardTop {
	var ourTopology boardTop
	ourTopology.nodeCount = n

	tempEdges, coords := getHexagonalLattice(n)
	ourTopology.cartesianCoords = coords
	ourTopology.edges = sparsifyEdges(n, tempEdges, coords)

	return ourTopology
}
