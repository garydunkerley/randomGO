package game

type Event interface {
	isEvent()
}

type ChallengeIssued struct {
	// metadata

	// TODO determine proper format for recording time?
	// defaulting to a string for the time being
	timeIssued   string
	challengerID string
	challengedID string
	// how long until the challenge expires, measured in seconds
	timeToRespond int

	// game parameters
	isRandom bool
	komi     float64
	size     int
	isRanked bool
	timer    timeStructure
}

func (e ChallengeIssued) isEvent() {}

// gameInitialized is an event that records who is black, who is white, and what the gameID is
type gameInitialized struct {
	black  string
	white  string
	gameID string
}

func (e gameInitialized) isEvent() {}

// this struct has less information than the
// boardState struct which is used internally and
// is intended to interface between the backend and
// Ebiten
type BoardUpdate struct {
	newStone      int
	newStoneColor int8

	newBlankStones []int

	blackCaptures int
	whiteCaptures int
	numberOfPlays int

	consecutivePasses int
}

func (e BoardUpdate) isEvent() {}

type EnteredScoring struct {
	gameID string

	blackCaptures int
	whiteCapture  int
}

func (e EnteredScoring) isEvent() {}

type GameEnded struct {
	// which game?
	gameID string
	// who is black?
	black string
	// who is white?
	white string

	whoWon string
}

func (e GameEnded) isEvent() {}

//TODO determine what other commands someone is likely to need to be able to issue to
// interact with the application

type Command interface {
	isCommand()

	// this method will declare who is
	// issuing the command
	issuer()
}

type GameInput struct {
	commandIssuer string

	// to which game is this command being issued
	gameid string

	// what is the command? Values between 0 and the number of nodes will
	// be interpreted as a stone placement.
	// -1 corresponds to a pass
	// -2 corresponds to forfeiture
	input int
}

func (c GameInput) isCommand() {}

func (c GameInput) issuer() string {
	return c.commandIssuer
}

type ChatInput struct {
	commandIssuer string

	gameID string
	msg    string
}

func (c ChatInput) isCommand() {}

func (c ChatInput) issuer() string {
	return c.commandIssuer
}

type ChallengePlayer struct {
	commandIssuer string

	// can be empty to indicate that
	// anyone may accept the challenge
	challengedID string

	isRandom bool
	komi     float64
	size     int
	handicap int
	timer    timeStructure
}

func (c ChallengePlayer) isCommand() {}

func (c ChallengePlayer) issuer() string {
	return c.commandIssuer
}

type RespondToChallenge struct {
	commandIssuer string

	// who sent the challenge?
	challengerID string

	// what is its ID?
	// challengeID will later pass to game ID
	challengeID string

	// was the challenge accepted?
	challengeAccepted bool
}

func (c RespondToChallenge) isCommand() {}

func (c RespondToChallenge) issuer() string {
	return c.commandIssuer
}

type timeStructure struct {
	primaryTimeLength int // in seconds

	hasByoYomi    bool
	numberByoYomi int
	lengthByoYomi int // in seconds
}

type Player struct {
	playerid     string
	currentGames []string
	pastGames    []string
	rank         int
}

type PlayerList struct {
	// a map taking playerids to their associated
	// Player structs
	players map[string]Player
}

type Game struct {
	gameID     string
	whiteID    string
	blackID    string
	dateOfGame string
	isRanked   bool

	// string is empty if game is ongoing
	victor string

	// chatLog records userid and text input as strings.
	chatLog [][2]string

	// TODO: create the file format for this
	gameData string

	isRandom bool
	size     int
	komi     float64
	handicap int
	timeStructure
}

// we can construct these by looking at Games that
// do not have a victor yet.
type OngoingGames struct {
	currentGames []Game
}
