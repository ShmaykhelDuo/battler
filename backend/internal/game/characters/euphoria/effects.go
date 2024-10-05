package euphoria

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

// EffectDescEuphoricSource is a description of [EffectEuphoricSource]
var EffectDescEuphoricSource = game.EffectDescription{
	Name: "Euphoric Source",
}

// Euphoric Source gives your Pink Sphere additional damage as well as well as healing while in Euphoria.
type EffectEuphoricSource struct {
	amount int
}

// NewEffectEuphoricSource returns a new [EffectEuphoricSource] of provided amount.
func NewEffectEuphoricSource(amount int) *EffectEuphoricSource {
	return &EffectEuphoricSource{amount: amount}
}

// Desc returns the effect's description.
func (e *EffectEuphoricSource) Desc() game.EffectDescription {
	return EffectDescEuphoricSource
}

// Amount returns the Euphoric Source's amount.
func (e *EffectEuphoricSource) Amount() int {
	return e.amount
}

// Increase increases the Euphoric Source's amount by specified amount.
func (e *EffectEuphoricSource) Increase(amount int) {
	e.amount += amount
}

// Decrease decreases the Euphoric Source's amount by specified amount.
func (e *EffectEuphoricSource) Decrease(amount int) {
	e.amount -= amount
}

// EffectDescUltimateEarly is a description of [EffectUltimateEarly]
var EffectDescUltimateEarly = game.EffectDescription{
	Name: "Ultimate Early",
}

// Your ultimate will unlock this many turns earlier than normal.
type EffectUltimateEarly struct {
	amount int
}

// NewEffectUltimateEarly returns a new [EffectUltimateEarly].
func NewEffectUltimateEarly() *EffectUltimateEarly {
	return &EffectUltimateEarly{amount: 1}
}

// Desc returns the effect's description.
func (e *EffectUltimateEarly) Desc() game.EffectDescription {
	return EffectDescUltimateEarly
}

// Amount returns the amount of shift.
func (e *EffectUltimateEarly) Amount() int {
	return e.amount
}

// Increase increases the amount of shift.
func (e *EffectUltimateEarly) Increase() {
	e.amount++
}

// ModifySkillUnlockTurn returns the modified turn number when skill is to be unlocked.
func (e *EffectUltimateEarly) ModifySkillUnlockTurn(s *game.Skill, unlockTurn int) int {
	if s.Desc().IsUltimate {
		unlockTurn -= e.amount
	}

	return unlockTurn
}

var EffectDescEuphoricHeal = game.EffectDescription{
	Name: "Euphoric Heal",
}

// You heal from Euphoric Source at the end of each turn, but Source gets rapidly depleted.
type EffectEuphoricHeal struct {
}

// Desc returns the effect's description.
func (e EffectEuphoricHeal) Desc() game.EffectDescription {
	return EffectDescEuphoricHeal
}

// OnTurnEnd executes the end-of-turn action.
func (e EffectEuphoricHeal) OnTurnEnd(c *game.Character, opp *game.Character, gameCtx game.Context) {
	eff := c.Effect(EffectDescEuphoricSource)
	source, ok := eff.(*EffectEuphoricSource)
	if !ok {
		return
	}

	amount := source.Amount()
	c.Heal(amount)
	opp.Heal(amount)

	if amount < 9 {
		source.Decrease(amount)
		return
	}
	source.Decrease(9)
}
