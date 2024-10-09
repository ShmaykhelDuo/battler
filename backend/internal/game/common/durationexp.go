package common

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

// DurationExpirable is a mixin that allows expiring effects after specified turn.
type DurationExpirable struct {
	expCtx game.Context
}

// NewDurationExpirable returns a [DurationExpirable] with specified expiry turn.
func NewDurationExpirable(expCtx game.Context) DurationExpirable {
	return DurationExpirable{expCtx: expCtx}
}

// TurnsLeft returns the number of turns left until this effect expires.
func (e DurationExpirable) TurnsLeft(gameCtx game.Context) int {
	fullTurns := e.expCtx.TurnNum - gameCtx.TurnNum

	if e.expCtx.IsGoingFirst && !gameCtx.IsGoingFirst {
		fullTurns -= 1
	}

	return fullTurns + 1
}

// HasExpired reports whether the effect has expired.
func (e DurationExpirable) HasExpired(gameCtx game.Context) bool {
	return gameCtx.IsAfter(e.expCtx)
}
