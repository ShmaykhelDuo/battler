package common

import "github.com/ShmaykhelDuo/battler/internal/game"

// DurationExpirable is a mixin that allows expiring effects after specified turn.
type DurationExpirable struct {
	expCtx game.TurnState
}

// NewDurationExpirable returns a [DurationExpirable] with specified expiry turn.
func NewDurationExpirable(expCtx game.TurnState) DurationExpirable {
	return DurationExpirable{expCtx: expCtx}
}

// TurnsLeft returns the number of turns left until this effect expires.
func (e DurationExpirable) TurnsLeft(turnState game.TurnState) int {
	fullTurns := e.expCtx.TurnNum - turnState.TurnNum

	if e.expCtx.IsGoingFirst && !turnState.IsGoingFirst {
		fullTurns -= 1
	}

	return fullTurns + 1
}

// HasExpired reports whether the effect has expired.
func (e DurationExpirable) HasExpired(turnState game.TurnState) bool {
	return turnState.IsAfter(e.expCtx)
}
