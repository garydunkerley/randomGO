package main

import (
	game "github.com/garydunkerley/randomGO/internal/game"
)

func main() {
	squareBoard := game.MakeSquareBoard(5, 5)
	squareBoard.RunGame()
}
