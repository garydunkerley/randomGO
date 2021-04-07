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

func BenchmarkGetLatticeAndSparsify(b *testing.B) {
	for i := 0; i < b.N; i++ {
		edges, coords := getHexagonalLattice(100)
		_ = sparsifyEdges(100, edges, coords)
	}
}

func BenchmarkGetCircuit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		edges, _ := getHexagonalLattice(100)
		_ = getCircuit(100, edges)
	}
}
