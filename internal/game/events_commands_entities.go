package game

/*
type Event interface {
	isEvent()
}

type Command interface {
	playerID()
	// will track who is issuing the command
	// should return nil if player is not signed in
	inGame()
	// determines whether the command was issued while
	// the player in question is in a game and, if so,
	// which one
	isValid()
	// using the first two pieces of information
	// this method should determine whether the
	// requested command is available to the user.
	// for instance, you cannot place a stone on a board
	// if you are not in a game.
}

// A game is started between two players
func (e BeganGame) isEvent() {}

// A user command requests that a stone is placed at a
// certain location, after which te dead stones are
// removed and new stone groups are created or joined
func (e BoardUpdate) isEvent() {}

//

// both players
func (e GameEnded) isEvent() {}

// the BeganGame struct
type BeganGame struct {
	white  string
	black  string
	gameID string
}

type BoardUpdate struct {
	newWhiteStone  int
	newBlackStone  int
	newBlankStones []int
	bool
}
*/
