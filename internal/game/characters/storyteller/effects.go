package storyteller

import (
	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/common"
)

// EffectDescCannotUse is a description of [EffectCannotUse]
var EffectDescCannotUse = game.EffectDescription{
	Name: "Can't use",
	Type: game.EffectTypeProhibiting,
}

// You can't use skills of the same colour you used last.
type EffectCannotUse struct {
	common.DurationExpirable
	colour game.Colour
}

// NewEffectCannotUse returns a new [EffectCannotUse] of provided colour.
func NewEffectCannotUse(turnState game.TurnState, colour game.Colour) EffectCannotUse {
	return EffectCannotUse{
		DurationExpirable: common.NewDurationExpirable(turnState.AddTurns(0, true)),
		colour:            colour,
	}
}

// Desc returns the effect's description.
func (e EffectCannotUse) Desc() game.EffectDescription {
	return EffectDescCannotUse
}

// Clone returns a clone of the effect.
func (e EffectCannotUse) Clone() game.Effect {
	return e
}

// Colour returns the forbidden skills' colour.
func (e EffectCannotUse) Colour() game.Colour {
	return e.colour
}

// IsSkillAvailable reports whether the skill can be used.
func (e EffectCannotUse) IsSkillAvailable(s *game.Skill) bool {
	return s.Desc().Colour != e.colour
}

// EffectDescControlled is a description of [EffectControlled]
var EffectDescControlled = game.EffectDescription{
	Name: "Controlled",
	Type: game.EffectTypeControl,
}

// This turn, your opponent chooses which skills you use.
type EffectControlled struct {
	common.DurationExpirable
}

func NewEffectControlled(turnState game.TurnState) EffectControlled {
	return EffectControlled{
		DurationExpirable: common.NewDurationExpirable(turnState.AddTurns(0, true)),
	}
}

// Desc returns the effect's description.
func (e EffectControlled) Desc() game.EffectDescription {
	return EffectDescControlled
}

// Clone returns a clone of the effect.
func (e EffectControlled) Clone() game.Effect {
	return e
}

// HasTakenControl reports whether the opponent has taken control over the character.
func (e EffectControlled) HasTakenControl() bool {
	return true
}
