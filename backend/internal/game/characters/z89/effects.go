package z89

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

// EffectDescUltimateSlow is a description of [EffectUltimateSlow]
var EffectDescUltimateSlow = game.EffectDescription{
	Name: "Ultimate Slow",
	Type: game.EffectTypeDebuff,
}

// Your ultimate will unlock this many turns later than normal.
type EffectUltimateSlow struct {
	amount int
}

// NewEffectUltimateSlow returns a new [EffectUltimateSlow] with the delay amount of 1.
func NewEffectUltimateSlow() *EffectUltimateSlow {
	return &EffectUltimateSlow{amount: 1}
}

// Desc returns the effect's description.
func (e *EffectUltimateSlow) Desc() game.EffectDescription {
	return EffectDescUltimateSlow
}

// Amount returns the amount of delay.
func (e *EffectUltimateSlow) Amount() int {
	return e.amount
}

// Increase increments the amount of delay.
func (e *EffectUltimateSlow) Increase() {
	e.amount++
}

// ModifySkillUnlockTurn returns the modified turn number when skill is to be unlocked.
func (e *EffectUltimateSlow) ModifySkillUnlockTurn(s *game.Skill, unlockTurn int) int {
	if s.Desc().IsUltimate {
		unlockTurn += e.amount
	}
	return unlockTurn
}
