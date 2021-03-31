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

func getDeadNodes(gg GoGraph, cs chromaticStrings) map[*node]bool {
	dead := make(map[*node]bool)

	for _, i := range cs.black {
		for j := range i.stones {
			dead[gg.nodes[j]] = true
		}
	}
	for _, i := range cs.white {
		for j := range i.stones {
			dead[gg.nodes[j]] = true
		}
	}
	return dead

}

func getWorkingDirectory() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return cwd
}

// initStone takes a *node (see boardStructs.go) and initializes a graphviz *Node
func initStone(myNode *node, g *cgraph.Graph, dead map[*node]bool) *cgraph.Node {
	stone, err := g.CreateNode(strconv.Itoa(myNode.id))
	if err != nil {
		log.Fatal(err)
	}
	stone.SetStyle("filled")
	stone.SetShape("circle")

	stone.SetFixedSize(true)

	if myNode.color == 0 {
		stone.SetFillColor("burlywood3")
	} else if myNode.color == 1 {
		if dead[myNode] {
			stone.SetFillColor("grey1")
		} else {
			stone.SetFillColor("black")
			stone.SetFontColor("white")
		}
	} else {
		if dead[myNode] {
			stone.SetFillColor("grey100")
		} else {
			stone.SetFillColor("white")
			stone.SetFontColor("black")
		}
	}
	stone.SetFontSize(20)

	return stone
}

// constructAllStones iterates over the stones in a GoGraph, initializes their stones (associated *Node structs from graphviz) andconstructs a map relating GoGraph *nodes to graphviz *Nodes
func constructAllStones(gg GoGraph, g *cgraph.Graph, dead map[*node]bool) map[*node]*cgraph.Node {
	transChart := make(map[*node]*cgraph.Node)

	// for each node in our GoGraph, we initialize a stone
	// and create a map assignment relating GoGraph *nodes
	// to graphviz *Nodes
	for i := 0; i < len(gg.nodes); i++ {
		transChart[gg.nodes[i]] = initStone(gg.nodes[i], g, dead)
	}
	return transChart

}

// constructAllEdges iterates the initEdges function over each *node in the GoGraph
func constructAllEdges(gg GoGraph, t map[*node]*cgraph.Node, g *cgraph.Graph) {

	edgeMade := make(map[string]bool)

	for i := 0; i < len(gg.nodes); i++ {
		n := gg.nodes[i]
		for _, z := range n.neighbors {
			// these strings will be used to check
			// whether an edge has already been constructed or not.

			edgeString1 := strconv.Itoa(z.id) + "," + strconv.Itoa(n.id)

			edgeString2 := strconv.Itoa(n.id) + "," + strconv.Itoa(z.id)

			if edgeMade[edgeString1] {
				continue
			} else {
				s1 := t[n]
				s2 := t[z]
				edge, err := g.CreateEdge("some_name", s1, s2)
				if err != nil {
					log.Fatal(err)
				}
				edge.SetArrowHead("none")
				edgeMade[edgeString1] = true
				edgeMade[edgeString2] = true
			}
		}

	}

	return
}

func visualizeBoard(gg GoGraph, isRandom bool, dead map[*node]bool) {

	cwd := getWorkingDirectory()

	g := graphviz.New()

	if isRandom {
		g.SetLayout("neato")
	} else {
		g.SetLayout("osage")
	}

	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}
	graph.SetBackgroundColor("burlywood3")

	// construct the graphviz board

	defer func() {
		if err := g.Close(); err != nil {
			log.Fatal(err)
		}
		// g.Close
	}()

	trans := constructAllStones(gg, graph, dead)
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
	return

}
