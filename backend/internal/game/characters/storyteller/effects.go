package storyteller

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

// EffectDescCannotUse is a description of [EffectCannotUse]
var EffectDescCannotUse = game.EffectDescription{
	Name: "Can't use",
	Type: game.EffectTypeProhibiting,
}

// You can't use skills of the same colour you used last.
type EffectCannotUse struct {
	colour game.Colour
}

// NewEffectCannotUse returns a new [EffectCannotUse] of provided colour.
func NewEffectCannotUse(colour game.Colour) EffectCannotUse {
	return EffectCannotUse{colour: colour}
}

// Desc returns the effect's description.
func (e EffectCannotUse) Desc() game.EffectDescription {
	return EffectDescCannotUse
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
}

// Desc returns the effect's description.
func (e EffectControlled) Desc() game.EffectDescription {
	return EffectDescControlled
}
