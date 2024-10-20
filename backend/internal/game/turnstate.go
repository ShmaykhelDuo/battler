package game

const (
	MinTurnNumber = 1  // minimum valid turn number
	MaxTurnNumber = 10 // maximum valid turn number
)

// TurnState contains any additional data needed to be used.
type TurnState struct {
	TurnNum      int  // the current character's turn number
	IsGoingFirst bool // whether the current character is going first
	IsTurnEnd    bool // whether it is the end of turn
}

// IsAfter reports whether the current context is indicating later game time than provided context.
func (c TurnState) IsAfter(other TurnState) bool {
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
func (c TurnState) AddTurns(turns int, isOpponentsTurn bool) TurnState {
	c.TurnNum += turns

	if isOpponentsTurn {
		if !c.IsGoingFirst {
			c.TurnNum++
		}

		c.IsGoingFirst = !c.IsGoingFirst
	}

	return c
}
