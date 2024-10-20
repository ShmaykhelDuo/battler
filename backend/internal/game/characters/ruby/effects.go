package ruby

import (
	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/common"
)

// EffectDescDoubleDamage is a description of [EffectDoubleDamage]
var EffectDescDoubleDamage = game.EffectDescription{
	Name: "Double Damage",
	Type: game.EffectTypeBasic,
}

// Doubles the damage you deal.
type EffectDoubleDamage struct {
	common.DurationExpirable
}

func NewEffectDoubleDamage(turnState game.TurnState) EffectDoubleDamage {
	return EffectDoubleDamage{
		DurationExpirable: common.NewDurationExpirable(turnState.AddTurns(2, false)),
	}
}

// Desc returns the effect's description.
func (e EffectDoubleDamage) Desc() game.EffectDescription {
	return EffectDescDoubleDamage
}

// ModifyDealtDamage returns the modified amount of damage based on provided amount of damage and damage colour.
func (e EffectDoubleDamage) ModifyDealtDamage(dmg int, colour game.Colour) int {
	return dmg * 2
}

// EffectDescCannotHeal is a description of [EffectCannotHeal]
var EffectDescCannotHeal = game.EffectDescription{
	Name: "Can't Heal",
	Type: game.EffectTypeProhibiting,
}

// Prevents you from healing.
type EffectCannotHeal struct {
	common.DurationExpirable
}

func NewEffectCannotHeal(turnState game.TurnState, isOpp bool) EffectCannotHeal {
	var expCtx game.TurnState
	if isOpp {
		expCtx = turnState.AddTurns(0, true)
	} else {
		expCtx = turnState.AddTurns(1, false)
	}

	return EffectCannotHeal{
		DurationExpirable: common.NewDurationExpirable(expCtx),
	}
}

// Desc returns the effect's description.
func (e EffectCannotHeal) Desc() game.EffectDescription {
	return EffectDescCannotHeal
}

// IsHealAllowed reports whether the healing is allowed based on provided amount of healing.
func (e EffectCannotHeal) IsHealAllowed(heal int) bool {
	return false
}
