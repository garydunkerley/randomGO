package game

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/goccy/go-graphviz"
	"github.com/skratchdot/open-golang/open"
)

type nodeStoneTranslation struct {
	n2s map[*node]*Node
}

func getWorkingDirectory() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return cwd
}

// initStone takes a *node (see boardStructs.go) and initializes a graphviz *Node
func (g Graph) initStone(myNode *node) *Node {

	stone, err := g.CreateNode(strconv.Itoa(myNode.id))
	if err != nil {
		log.Fatal(err)
	}
	stone.SetStyle("filled")
	if myNode.color == 0 {
		stone.SetFillColor("grey")
	} else if myNode.color == 1 {
		stone.SetFillColor("black")
	} else {
		stone.SetFilleColor("white")
	}

	return stone
}

// constructAllStones iterates over the stones in a GoGraph, initializes their stones (associated *Node structs from graphviz) andconstructs a map relating GoGraph *nodes to graphviz *Nodes
func (g Graph) constructAllStones(gg GoGraph) map[*node]*Node {
	transChart := make(map[*node]*Node)

	// for each node in our GoGraph, we initialize a stone
	// and create a map assignment relating GoGraph *nodes
	// to graphviz *Nodes
	for i := 0; i < len(gg.nodes); i++ {
		transChart[gg.nodes[i]] = g.initStone(gg.nodes[i])
	}
	return transChart

}

// constructAllEdges iterates the initEdges function over each *node in the GoGraph
func (g Graph) constructAllEdges(gg GoGraph, t nodeStoneTranslation) {
	for i := 0; i < len(gg.nodes); i++ {
		n := gg.nodes[i]
		initEdges(n, t)
	}

	return
}

// initEdges takes a GoGraph node pointer and our struct associating these node pointers to graphviz Nodes to construct an edge connecting each neighbor
func (g Graph) initEdges(n *node, t nodeStoneTranslation) {
	s1 := t.n2s[n]
	var s2 *Node
	for z := range n.neighbors {
		s2 = t.n2s[z]
		edge, err := g.CreateEdge("some_name", s1, s2)
		if err != nil {
			log.Fatal(err)
		}
		edge.SetArrowHead("none")
	}

	return
}

func visualizeBoard(gg GoGraph) {

	var trans nodeStoneTranslation

	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}

	// construct the graphviz board

	defer func() {
		if err := g.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close
	}()

	trans.n2s = graph.constructAllStones(gg)
	graph.constructAllEdges(gg, trans.n2s)

	var buf bytes.Buffer
	if err := g.Render(graph, graphviz.PNG, &buf); err != nil {
		log.Fatal(Err)
	}

	// edit so that the directory is what you want it to be
	if err := g.RenderFilename(graph, graphviz.PNG, cwd+"/board.png"); err != nil {
		log.Fatal(err)
	}
	open.Run(cwd + "/board.png")

}
