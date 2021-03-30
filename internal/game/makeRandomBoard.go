package game

import (
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
func euclideanNorm(x []float64) (int, error) {
	var sum float64
	for i := 0; i < len(x); i++ {
		sum += x[i]**2	
	}
	return math.Sqrt(sum)
}

// getHexagonalNumbers computes a slice of hexagonal numbers that are
// less then or equal to the input n
func getHexagonalNumbers(n int) []int {
	var hexNums []int
	searching := true
	count := 1

	if n > 0 {
		hexNums = append(hexNums, 1)
		for searching {
			nextHex := sum(hexNums) + count*6  
			if nextHex =< n {
				hexNums = append(hexNums, nextHex)
				count += 1
			} else {
				searching = false	
			}
		}
		
	} else {
		fmt.Println("Please input a positive integer")
	}

	return hexNums
}

// generateHexagonalLattice takes an integer and outputs the information
// nececessary to create a hexagonal goban
func generateHexagonalLattice(n int) (map[int][]int, map[int][]float64) {

	theta := math.Pi / 6
	multiplier := 0

	edges := make(map[int][]int)
	coordMap := make(map[int][]float64)

	ourTopology.nodeCount = n

	// roomsMade will track how far we have iterated along
	// direction to create a layer in Pascal's triangle.
	cellsMade := 1

	// counter tracks the number of nodes
	// that have been generated so far
	counter := 1

	// the first node will be centered at the coordinate 0,0

	currentVec = []float64{0, 0}
	coordMap[0] = currentVec
	existsAlready[0] = true

	// we will also initialize the direction of generation 
	// as being trivial
	genDirection := []float64{0, 0}

	hexNums := getHexagonalNumbers(n) 
	hexIterate := 0

	for counter < n {
		if counter == hexNums[hexIterate] {
			// move along the current direction to the
			// node that was seen previously
			currentVec = currentVec + genDirection
			multiplier = 0
			genDirection = []float64{1, 0}

			if hexIterate < len(hexNums)-1 {
				hexIterate += 1
			}

		}

		newVec := currentVec + genDirection
		coordMap[counter] = newVec

		// After establishing the existence of a new node,
		// we attach edges between it and nodes we have already created.

		// if counter happens to be a hexagonal number and greater than 1, 
		// we ensure there is an edge between it and the previous hexagonal number.
		if counter == hexNums[hexIterate-1] && counter > 1 {
	
			edges[counter] = append(edges[counter], hexNums[hexIterate-1])
			edges[hexNums[hexIterate-1]] = append(edges[hexNums[hexIterate-1]], counter)
			
		} else {
			
			if hexIterate == 1 {
				edges[counter] = append(edges[counter], 0)
			} else {
				first := counter - (6 * layerCount)
				second := counter - (6*layercount - 1)

				edges[counter] = append(edges[counter], first, second)
				edges[first] = append(edges[first], counter)
				edges[second] = append(edges[second], counter)
			}
			edges[counter] = append(edges[counter], counter-1)
			edges[counter-1] append(edges[counter-1], counter)

		}

		cellsMade += 1
		if cellsMade == hexIterate {

			multiplier = (multiplier + 1) % 6
			genDirection = []float64{math.Sin(multiplier * theta), math.Cos(multiplier * theta)}
			cellsMade = 0
			
		}
	}
	return edges, coordMap
}

// makeAscending rearranges an integer slice so that its elements are in ascending order
func makeAscending(mySlice []int) []int {
	targetLength = len(mySlice)

	var ascendingSlice []int
	for len(ascendingSlice) < targetLength {
		currentMinIndex := 0
		// iterate over the elements of your slice
		// when you find an element smaller than the 
		// current minimum, record it and move on.
		for i := 0; i< len(mySlice); i++ {
			if mySlice[i] =< mySlice[currentMinIndex] {
				currentMinIndex = i
			}
		}
		// add the smallest element to the ascending slice
		// and then remove it from the original
		ascendingSlice = append(ascendingSlice, mySlice[smallestIndex])
		mySlice = append(mySlice[:smallestIndex], mySlice[smallestIndex+1:])
	}
	

}

// getCircuit takes our collection of edges and outputs an
// encoding of a spanning tree that is also a path
func getCircuit(n int, edges map[int][]int) map[int][]int{
	circuit = []
	circuitMap := make(map[int][]int)
	visited := make(map[int]bool)
	

	currentNode := 0
	
	for len(visited) < n {
		ascendingEdges = makeAscending(edges[currentNode])
		
		for i:= 0; i<len(ascendingEdges){
			if !(visited[i]) {
				circuit = append(circuit, currentNode)
				currentNode = ascendingEdges[i]
				visited[currentNode]
				break
			}
		}
	}

	for i := 0; i<len(circuit); i++ {
		if i>0 && i<n-1 {
			circuitMap[circuit[i]] = append(circuitMap[circuit[i]], circuit[i-1], circuit[i+1]) 
		} else if i==0 {
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
	
	newEdges := make(map[int][]int)
	safeEdges := make(map[int][]int)
	

	curcuit := getCircuit(n, edges)

	for i := 0; i < n; i++ {
		// TODO: The probability that an edge is deleted should be determined in part by how far 
		// away the node in question is from the "center" of the board (0).
		// We need a function that is bounded above by 1, and increases as 
		// the distance of our node from zero increases
		norm = euclideanNorm(coordMap[i])
		prob = math.Pi / 2*(math.Atan(norm + 0.05))
		
		validEdges = []
		dontInclude := make(map[int]bool)

		for j := 0; j<len(edges[i]); j++ {
			for k := 0; k < len(safeEdges[i]); k++ {
				if safeEdges[i][k] == edges[i][j] {
					dontInclude[edges[i][j]]=true
					break
				}
			}
			if !dontInclude[edges[i][j]] {
				validEdges = append(validEdges, edges[i][j])
			}
		}

		// if a node has more than 1 edge, then its edges will be subjected to 
		// a random elimination process
		if len(validEdges)>=2 {

			stillRolling := true
			for stillRolling {

				rand.Seed(time.Now().UnixNano())
				v := rand.Float64()
		
				if v < prob {
					randomIndex := rand.Intn(len(validEdges))
					validEdges = append(validEdges[:randomIndex], validEdges[randomIndex+1:])
				
				} else {
					stillRolling = false
				}

			}

			for j := 0 ; j < len(validEdges)-1 ; j++ {
				newEdges[i]=append(newEdges[i],j)
				newEdges[j]=append(newedges[j],i)

				// We declare any edges that survive this process to be safe 
				// and not in danger of being eliminated later
				safeEdges[j] = append(safeEdges[j],i)
				safeEdges[i] = append(safeEdges[i],j)
			
			}
		}
	}
	return newEdges
}


func addBoundaryEdges(n int, edges map[int][]int) map[int][]int {
//TODO make a map that will be guaranteed a planar representation.

	hexNums := getHexagonalNumbers(n)

	boundary = []int

	// TODO how to find the boundary?




return
}

func makeRandomBoard(n int) randomBoardTop {
	var ourTopology randomBoardTop
	ourTopology.nodeCount = n

	tempEdges, ourTopology.coords := getHexagonalLattice(n)
	
	ourTopology.edges = tempEdges

	return ourTopology
}
