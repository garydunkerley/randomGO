package game

type Event interface {
	isEvent()
}

/*
//TODO: when this event occurs, both players are informed
// of the pending challenge
type ChallengeIssued struct {
	// metadata
	gameID string

	// TODO determine proper format for recording time?
	// defaulting to a string for the time being
	timeIssued   string
	challengerID string

	// can be left blank so as to be a general challenge that
	// anyone can accept (possibly with rank restrixtions thoguh)
	challengedID string
	// how long until the challenge expires, measured in seconds, possibly infinite
	timeToRespond int

	// game parameters
	isRandom bool
	komi     float64
	size     int

	// positive to benefit the challenger, negative to
	// benefit the opponent
	handicap int

	isRanked bool

	// this is a struct that will determine how long the
	// main part of the game will be, whether there are
	// Byo-Yomi, how many, and how long.
	timeStructure
}


func (e ChallengeIssued) isEvent() {}


type ChallengeAccepted struct {
	gameID string
}

func (e ChallengeAccepted) isEvent() {}
*/

/*
// an event that triggers after both players have
// submitted their votes for whether the board should
// be rerolled.

type rerollVoteCompleted struct {
	gameID    string
	whiteVote bool
	blackVote bool
}

func (e rerollVoteCompleted) isEvent() {}
*/

// user entity is defined by a unique int id and a name string
type user struct {
	name string
	id   int
}

// mockUserDatabase entity tracks multiple users and ensures unique user ids
type mockUserDatabase struct {
	userLookup map[int]user // look up a user by ID
	nextUserID int          // track largest nonnegative ID that has not been used
}

// addUser creates a new user from given username and increments nextUserID
func (db *mockUserDatabase) addUser(username string) {
	newID := db.nextUserID
	newUser := user{name: username, id: newID}
	db.userLookup[newID] = newUser
	db.nextUserID++
}

// createUserDB instantiates a new mockUserDatabase
func createUserDB() mockUserDatabase {
	users := make(map[int]user)
	return mockUserDatabase{userLookup: users}
}

// gameInstance is defined by a unique int id and two users
// TODO: add something which connects to game logic
type gameInstance struct {
	black  user
	white  user
	gameID int
}

// mockGameDatabase entity tracks multiple games and ensures unique game ids
type mockGameDatabase struct {
	gameLookup map[int]gameInstance // look up a game by ID
	nextGameID int                  // track largest nonnegative ID that has not been used
}

// addGame creates a new game between given users and increments nextGameID
func (db *mockGameDatabase) addGame(blackPlayer user, whitePlayer user) {
	newID := db.nextGameID
	newGame := gameInstance{black: blackPlayer, white: whitePlayer, gameID: newID}
	db.gameLookup[newID] = newGame
	db.nextGameID++
}

// createGameDB instantiates a new mockGameDatabase
func createGameDB() mockGameDatabase {
	games := make(map[int]gameInstance)
	return mockGameDatabase{gameLookup: games}
}

// StoneUpdate describes a set of stone deltas, e.g. the consequence of
// playing a given move, loading an initial handicap, or loading a past game.
// TODO: Client's board state should be able to execute StoneUpdate as a command
type StoneUpdate struct {
	stones map[int]int8 // pairings of (stone id): (new color)
	err    error
}

// TODO: Client needs a command to request a StoneUpdate (i.e. try to play a move and receive StoneUpdate)

/*
type EnteredScoring struct {
	gameID string
}

func (e EnteredScoring) isEvent() {}

type ResumedFromScoring struct {
	gameID string
}

func (e ResumedFromScoring) isEvent() {}

type GameEnded struct {
	// which game?
	gameID string
}

func (e GameEnded) isEvent() {}

//TODO determine what other commands someone is likely to need to be able to issue to
// interact with the application
*/

//TODO: sort out command structuring and GameInput.
// Who's responsible for receiving inputs? What's the communication model?

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

/*
type ChatInput struct {
	commandIssuer string

	gameID string
	msg    string
}

func (c ChatInput) isCommand() {}

func (c ChatInput) issuer() string {
	return c.commandIssuer
}
*/

/*
// It stands to reason that when generating a board, one or more
// players may be unsatisfied with the outcome, especially since
// there are low probability boards that are really bad, like,
// almost every node is valence two. It would be good
// to allow for the board to be reconstructed some limited number
// of times before play begins to make sure everyone is
// reasonably satisfied with the goban they will be playing on
type gobanRerollVote struct {
	commandIssuer string

	vote bool
}

func (c gobanRerollVote) isCommand() {}

func (c gobanRerollVote) issuer() string {
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
	timeStructure
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
*/
type timeStructure struct {
	primaryTimeLength int // in seconds

	hasByoYomi    bool
	numberByoYomi int
	lengthByoYomi int // in seconds
}

/*
type Player struct {
	playerid       string
	profilePicture string // URL (probably to Imgur) to import profilePictures
	nationality    string
	activeGames    []string
	pastGames      []string
	rank           int

	// additional statistics, like percentage of games won, etc
}

type PlayerList struct {
	// a map taking playerids to their associated
	// Player structs
	players map[string]Player
}
*/

// TODO: rewrite this struct after determining communication model
type Game struct {
	gameID     string
	gameName   string
	whiteID    string
	blackID    string
	dateOfGame string
	isRanked   bool
	isRandom   bool
	size       int
	komi       float64
	handicap   int
	timeStructure

	// string is empty if game is ongoing
	victor string

	blackScore float64
	whiteScore float64

	// chatLog records userid and text input as strings.
	chatLog [][2]string

	// TODO: create the file format for this
	gameData string
}

/* NOTE: think this is deprecated by the mock dbs --jr
// we can construct these by looking at Games that are the subject of an initialization event
// but not an ending event.
type OngoingGames struct {
	ongoingGames []Game
}
*/
