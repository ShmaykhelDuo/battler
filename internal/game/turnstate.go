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

// NewTurnState returns a new [TurnState] with provided turn number.
func NewTurnState(turnNum int) TurnState {
	return TurnState{TurnNum: turnNum}
}

func StartTurnState() TurnState {
	return TurnState{
		TurnNum:      MinTurnNumber,
		IsGoingFirst: true,
	}
}

// WithGoingFirst returns a copy of the state with going first set.
func (c TurnState) WithGoingFirst(isGoingFirst bool) TurnState {
	c.IsGoingFirst = isGoingFirst
	return c
}

// WithTurnEnd returns a copy of the state with turn end set.
func (c TurnState) WithTurnEnd() TurnState {
	c.IsTurnEnd = true
	return c
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

func (c TurnState) Next() TurnState {
	return c.AddTurns(0, true)
}
