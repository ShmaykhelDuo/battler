package game

const (
	MinTurnNumber = 1  // minimum valid turn number
	MaxTurnNumber = 10 // maximum valid turn number
)

// Context contains any additional data needed to be used.
type Context struct {
	TurnNum      int  // the current character's turn number
	IsGoingFirst bool // whether the current character is going first
	IsTurnEnd    bool // whether it is the end of turn
}

func TurnCtx(turnNum int) Context {
	return Context{TurnNum: turnNum}
}

func (c Context) WithGoingFirst(isGoingFirst bool) Context {
	c.IsGoingFirst = isGoingFirst
	return c
}

func (c Context) WithTurnEnd() Context {
	c.IsTurnEnd = true
	return c
}

// IsAfter reports whether the current context is indicating later game time than provided context.
func (c Context) IsAfter(other Context) bool {
	if c.TurnNum != other.TurnNum {
		return c.TurnNum > other.TurnNum
	}

	if c.IsGoingFirst != other.IsGoingFirst {
		return other.IsGoingFirst
	}

	return c.IsTurnEnd && !other.IsTurnEnd
}

// AddTurns returns the context of provided number of turns ahead.
// Is isOpponentsTurn is true, the returned context is of the nearest opponent's turn.
func (c Context) AddTurns(turns int, isOpponentsTurn bool) Context {
	c.TurnNum += turns

	if isOpponentsTurn {
		if !c.IsGoingFirst {
			c.TurnNum++
		}

		c.IsGoingFirst = !c.IsGoingFirst
	}

	return c
}
