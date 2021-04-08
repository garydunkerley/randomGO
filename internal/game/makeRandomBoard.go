package game

import (
	//	"fmt"
	"math"
	"math/rand"
	"strconv"
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
func euclideanNorm(x [2]float64) float64 {
	var sum float64
	for _, i := range x {
		sum += math.Pow(i, 2)
	}
	return math.Sqrt(sum)
}

// vecSum just takes two floating point arrays of length two and finds their coordinate-wise sum
func vecSum(x [2]float64, y [2]float64) [2]float64 {
	var z [2]float64
	for i := 0; i < len(x); i++ {
		z[i] = x[i] + y[i]
	}
	return z
}

// getHexagonalNumbers computes a slice of centered hexagonal numbers less than or equal to n
// For more information, see: https://oeis.org/A003215
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
func getHexagonalLattice(n int) (map[int][]int, map[int][2]float64) {

	theta := 2 * math.Pi / 6
	multiplier := float64(4)

	edges := make(map[int][]int)
	coordMap := make(map[int][2]float64)

	// roomsMade will track how far we have iterated along
	// direction to create a layer in Pascal's triangle.
	cellsMade := 1

	// counter tracks the number of nodes
	// that have been generated so far
	counter := 1

	// the first node will be centered at the coordinate 0,0

	currentVec := [2]float64{0, 0}
	coordMap[0] = currentVec

	// we will also initialize the direction of generation
	// as being trivial
	genDirection := [2]float64{1, 0}

	hexNums := getCenteredHexagonalNumbers(n)
	hexIterate := 0

	corner := 1
	cornerCount := 1

	for counter < n {
		levelingNumber := int(math.Min(float64(hexIterate+1), float64(len(hexNums)-1)))

		if counter == hexNums[levelingNumber] {
			// move along the current direction to the
			// node that was seen previously

			currentVec = vecSum(currentVec, genDirection)
			multiplier = float64(4)
			genDirection = [2]float64{1, 0}

			if hexIterate < len(hexNums)-1 {
				hexIterate += 1
			}

		}

		newVec := vecSum(currentVec, genDirection)
		coordMap[counter] = newVec

		currentVec = newVec

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

			multiplier = float64((int(multiplier) + 1) % 6)
			genDirection = [2]float64{-math.Cos(multiplier * theta), -math.Sin(multiplier * theta)}
			cellsMade = 0

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

			cellsMade += 1

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

			cellsMade += 1
		}

		counter += 1
	}

	return edges, coordMap
}

// getCircuit takes our collection of edges and outputs an
// encoding of a spanning tree
func getCircuitMap(n int, edges map[int][]int) map[int][]int {

	var circuit []int

	circuitMap := make(map[int][]int)
	visited := make(map[int]bool)

	currentNode := n - 1

	for len(circuit) < n {
		circuit = append(circuit, currentNode)
		visited[currentNode] = true

		//TODO: from the neighbors of the current node, find the one that
		// minimizes the distance and hasn't been visited yet
		var unvisited []int
		for i := 0; i < len(edges[currentNode]); i++ {
			if !(visited[edges[currentNode][i]]) {
				unvisited = append(unvisited, edges[currentNode][i])
			}
		}
		smallest := math.Inf(1)
		for _, i := range unvisited {
			if math.Sqrt(math.Pow(float64(currentNode-i), 2)) < math.Sqrt(math.Pow(float64(currentNode)-smallest, 2)) {
				smallest = float64(i)
			}
		}

		currentNode = int(smallest)
	}

	// build a map that encode the incidence behavior of the nodes forming the circuit.
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
func sparsifyEdges(n int, edges map[int][]int, coordMap map[int][2]float64) map[int][]int {

	safeEdges, circuitMap := initSafeEdges(n, edges)

	newEdges := removeRandomCandidates(coordMap, safeEdges, edges, circuitMap)

	return newEdges
}

// initSafeEdges looks at the circuit for the lattice and stores all edges in a map
// so that they will not be removed later.
func initSafeEdges(n int, edges map[int][]int) (map[string]bool, map[int][]int) {

	safeEdges := make(map[string]bool)

	// declare all edges in the circuit to be safe from elimination
	// to ensure that the sparsified graph is always connected.
	circuitMap := getCircuitMap(n, edges)

	for i := 0; i < n; i++ {
		for _, j := range circuitMap[i] {
			newSafeEdge := strconv.Itoa(i) + "," + strconv.Itoa(j)
			safeEdges[newSafeEdge] = true

		}
	}
	return safeEdges, circuitMap

}

// removeRandomCandidates will locate the candidate edges that may be subject to the random
// elimination process and then performs said process.
func removeRandomCandidates(coordMap map[int][2]float64, safeEdges map[string]bool, edges map[int][]int, circuitMap map[int][]int) map[int][]int {

	newEdges := make(map[int][]int)
	newEdges = circuitMap

	newSafeEdges := make(map[string]bool)
	newSafeEdges = safeEdges

	for i := 0; i < len(coordMap); i++ {

		candidateEdges := getRemovalCandidates(i, newSafeEdges, edges)
		// if a node has an unsafe edge, then its edges will be subjected to
		// a random elimination process

		if len(candidateEdges) > 0 {
			norm := euclideanNorm(coordMap[i])

			stillRolling := true

			for stillRolling && len(candidateEdges) > 0 {

				//TODO: tweak this until I get something I like.
				// Right now, being farther away from the center increases the likelihood that
				// a node will lose an edge, as does having a lot of edges

				prob := math.Pow((2*math.Atan(norm))/math.Pi, float64(4-len(candidateEdges)))
				// prob := float64(0)
				rand.Seed(time.Now().UnixNano())
				v := rand.Float64()

				if v < prob {
					randomIndex := rand.Intn(len(candidateEdges))

					candidateEdges = append(candidateEdges[:randomIndex], candidateEdges[randomIndex+1:]...)
				} else {
					stillRolling = false
				}

			}

			// all edges that survive this process are declared to be safe
			if len(candidateEdges) > 0 {
				for _, j := range candidateEdges {
					// initialize the coordinate strings for the safe edge
					newSafeString1 := strconv.Itoa(i) + "," + strconv.Itoa(j)
					newSafeString2 := strconv.Itoa(j) + "," + strconv.Itoa(i)

					// declare them safe
					newSafeEdges[newSafeString1] = true
					newSafeEdges[newSafeString2] = true

					// add the newly safe edges to the new edge map for each node
					newEdges[i] = append(newEdges[i], j)
					newEdges[j] = append(newEdges[j], i)
				}
			}
		}
	}

	return newEdges

}

// getRemovalCandidates takes a node and creates a list of adjacent nodes whose associated edges are not safe
func getRemovalCandidates(i int, safeEdges map[string]bool, edges map[int][]int) []int {
	var candidateEdges []int

	// iterate over all the edges attached to a node i and check to see if they are
	// safe (i.e. that they belong to the circuit computed earlier or they have
	// survived a roll already)
	// edges that are not safe may be subject to elimination when the rolling occurs.
	for _, edge := range edges[i] {
		stringToCheck := strconv.Itoa(i) + "," + strconv.Itoa(edge)
		if !safeEdges[stringToCheck] {
			candidateEdges = append(candidateEdges, edge)
		}
	}
	return candidateEdges

}

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
