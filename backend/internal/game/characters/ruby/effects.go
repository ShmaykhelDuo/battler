package ruby

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

// EffectDescDoubleDamage is a description of [EffectDoubleDamage]
var EffectDescDoubleDamage = game.EffectDescription{
	Name: "Double Damage",
}

// Doubles the damage you deal.
type EffectDoubleDamage struct {
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
}

// Prevents you from healing.
type EffectCannotHeal struct {
}

// Desc returns the effect's description.
func (e EffectCannotHeal) Desc() game.EffectDescription {
	return EffectDescCannotHeal
}

// IsHealAllowed reports whether the healing is allowed based on provided amount of healing.
func (e EffectCannotHeal) IsHealAllowed(heal int) bool {
	return false
}
