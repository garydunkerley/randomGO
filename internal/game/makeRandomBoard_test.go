package game

import "testing"

func BenchmarkMakeRandomBoard(b *testing.B) {
	for i := 0; i < b.N; i++ {
		makeRandomBoard(100)
	}
}

func BenchmarkGetHexagonalLattice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getHexagonalLattice(100)
	}
}

func BenchmarkSparsify100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		edges, coords := getHexagonalLattice(400)
		_ = sparsifyEdges(400, edges, coords)
	}
}

/*
func BenchMarkRemoveRandomCandidates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		edges, coordMap := getHexagonalLattice(100)
		safeEdges := initSafeEdges(100, edges)

		safeEdges = removeRandomCandidates(coordMap, safeEdges, edges)
	}
}
*/

func BenchmarkGetCircuit(b *testing.B) {

	for i := 0; i < b.N; i++ {
		edges, _ := getHexagonalLattice(100)
		_ = getCircuitMap(100, edges)
	}
}
