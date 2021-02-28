package game

import (
	"bytes"
	"log"
	"os"
	"strconv"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"

	"github.com/skratchdot/open-golang/open"
)

/*
type nodeStoneTranslation struct {
	n2s map[*node]*cgraph.Node
}
*/

func getWorkingDirectory() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return cwd
}

// initStone takes a *node (see boardStructs.go) and initializes a graphviz *Node
func initStone(myNode *node, g *cgraph.Graph) *cgraph.Node {

	stone, err := g.CreateNode(strconv.Itoa(myNode.id))
	if err != nil {
		log.Fatal(err)
	}
	stone.SetStyle("filled")
	stone.SetShape("circle")

	stone.SetFixedSize(true)

	if myNode.color == 0 {
		stone.SetFillColor("grey")
	} else if myNode.color == 1 {
		stone.SetFillColor("black")
	} else {
		stone.SetFillColor("white")
	}

	return stone
}

// constructAllStones iterates over the stones in a GoGraph, initializes their stones (associated *Node structs from graphviz) andconstructs a map relating GoGraph *nodes to graphviz *Nodes
func constructAllStones(gg GoGraph, g *cgraph.Graph) map[*node]*cgraph.Node {
	transChart := make(map[*node]*cgraph.Node)

	// for each node in our GoGraph, we initialize a stone
	// and create a map assignment relating GoGraph *nodes
	// to graphviz *Nodes
	for i := 0; i < len(gg.nodes); i++ {
		transChart[gg.nodes[i]] = initStone(gg.nodes[i], g)
	}
	return transChart

}

// constructAllEdges iterates the initEdges function over each *node in the GoGraph
func constructAllEdges(gg GoGraph, t map[*node]*cgraph.Node, g *cgraph.Graph) {
	for i := 0; i < len(gg.nodes); i++ {
		n := gg.nodes[i]
		initEdges(n, t, g)
	}

	return
}

// initEdges takes a GoGraph node pointer and our struct associating these node pointers to graphviz Nodes to construct an edge connecting each neighbor
func initEdges(n *node, t map[*node]*cgraph.Node, g *cgraph.Graph) {
	s1 := t[n]
	var s2 *cgraph.Node
	for _, z := range n.neighbors {
		s2 = t[z]
		edge, err := g.CreateEdge("some_name", s1, s2)
		if err != nil {
			log.Fatal(err)
		}
		edge.SetArrowHead("none")
	}

	return
}

func visualizeBoard(gg GoGraph) {

	var trans map[*node]*cgraph.Node

	cwd := getWorkingDirectory()

	g := graphviz.New()
	g.SetLayout("osage")

	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}

	// construct the graphviz board

	defer func() {
		if err := g.Close(); err != nil {
			log.Fatal(err)
		}
		// g.Close
	}()

	trans = constructAllStones(gg, graph)
	constructAllEdges(gg, trans, graph)

	var buf bytes.Buffer
	if err := g.Render(graph, graphviz.PNG, &buf); err != nil {
		log.Fatal(err)
	}

	// edit so that the directory is what you want it to be
	if err := g.RenderFilename(graph, graphviz.PNG, cwd+"/board.png"); err != nil {
		log.Fatal(err)
	}
	open.Run(cwd + "/board.png")

}
