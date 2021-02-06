package main

import (
	boards "github.com/garydunkerley/randomGO/internal/boards"
)

func main() {
	theBoard := boards.makeSquareBoard(5, 5)
	theBoard.runGame()

}
