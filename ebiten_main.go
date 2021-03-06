package main

import (
	game "github.com/garydunkerley/randomGO/internal/game"
	"math"

	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// GLOBAL PARAMETERS

// boardInfo holds all information relevant to the creation and arrangment of
// the various assets the game will be using during play
var boardInfo game.EbitenBoardInfo

// used for board topology generation
var isRandom bool

// stoneRadius will tell us how large our stone assets should be
var stoneRadius float64

// buttonRadius should be computed using geometric information
// about the graph embedding. It should be at least as big as the// pciture of the stone, but less than half the width of the shortest distance between two nodes in our graph.
var buttonRadius float64

// Variables corresponding to graphical assets
var wood *ebiten.Image
var graph *ebiten.Image
var whiteStone *ebiten.Image
var blackStone *ebiten.Image

type ebitenNode struct {
	picture    *ebiten.Image
	id         int
	xPos, yPos float64
}

var ebitenNodeMap map[int]ebitenNode

type Game struct{}

// TODO:
// Figure out how to pipe in the board generated by graphviz
// as well as the geometric data.
// I think we can use graphviz to output the raw .png file for the board lines.
// To get the geometric data, we're going to need to pipe in the // dot data.

// We also need to figure out how to get this to interface with
// the backend.

func (g *Game) Update() error {
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	width, height := wood.Size()
	return width + 500, height
}

func (g *Game) Draw(screen *ebiten.Image) {

	// Draws the wood for the board
	screen.DrawImage(wood, nil)

	// Adds the graphviz PNG asset
	// with suitable translation and rescaling
	// to ensure that it is centered

	gWidth, gHeight := graph.Size()
	wWidth, wHeight := wood.Size()
	rescale := 0.95

	rWidth := rescale * float64(wWidth) / float64(gWidth)
	rHeight := rescale * float64(wHeight) / float64(gHeight)

	tWidth := ((1 - rescale) * float64(wWidth)) / 2
	tHeight := ((1 - rescale) * float64(wHeight)) / 2

	graphPos := &ebiten.DrawImageOptions{}
	graphPos.GeoM.Scale(rWidth, rHeight)
	graphPos.GeoM.Translate(tWidth, tHeight)

	screen.DrawImage(graph, graphPos)
}

func init() {
	var err error

	// we initialize the board topology and generate the
	// PNG file for the graph defining our goban
	// we need to modify BuildRandomGame so that it will
	// save the PNG coordinates to boardInfo
	boardInfo, isRandom = game.BuildRandomGame(100)

	// populate all Ebiten assets with their relevant files
	wood, _, err = ebitenutil.NewImageFromFile("assets/woodgrain3.png")
	graph, _, err = ebitenutil.NewImageFromFile("assets/temporary/board.png")
	whiteStone, _, err = ebitenutil.NewImageFromFile("assets/go_stone_white.png")
	blackStone, _, err = ebitenutil.NewImageFromFile("assets/go_stone_black.png")

	if err != nil {
		log.Fatal(err)
	}

}

func isNodeClicked(x, y int, ebitenNodeMap map[int]ebitenNode) (int, bool) {
	for _, node := range ebitenNodeMap {
		if math.Sqrt(math.Pow(float64(x)-node.xPos, 2)+math.Pow(float64(y)-node.yPos, 2)) <= buttonRadius {
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				return node.id, true
			}
		}
	}

	return -1, false
}

/*
func getPlayerInput() {

	return
}

func placeStone() {
	if CursorPosition() ==

}

func update(screen *ebiten.Image) error {
	return

}
*/

func main() {

	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("it goes it goes it goes it goes it goes it goes it goes it goes")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
	game.StartGame(boardInfo, isRandom)
}
