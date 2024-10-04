package game

const (
	MinTurnNumber = 1  // minimum valid turn number
	MaxTurnNumber = 10 // maximum valid turn number
)

// Context contains any additional data needed to be used.
type Context struct {
	TurnNum      int  // the current character's turn number
	IsGoingFirst bool // whether the current character is going first
}
