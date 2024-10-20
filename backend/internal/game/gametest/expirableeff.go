package gametest

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

var EffectDescExpirable = game.EffectDescription{
	Name: "Expirable",
}

type EffectExpirable struct {
	expired bool
}

func NewEffectExpirable(expired bool) *EffectExpirable {
	return &EffectExpirable{expired: expired}
}

// Desc returns the effect's description.
func (e *EffectExpirable) Desc() game.EffectDescription {
	return EffectDescExpirable
}

func (e *EffectExpirable) Expire() {
	e.expired = true
}

// HasExpired reports whether the effect has expired.
func (e *EffectExpirable) HasExpired(turnState game.TurnState) bool {
	return e.expired
}
