package structure

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

var EffectDescIBoost = game.EffectDescription{
	Name: "I Boost",
	Type: game.EffectTypeNumeric,
}

// Boosts your Electric Shock damage and Layers defense.
type EffectIBoost struct {
	amount int
}

func NewEffectIBoost(amount int) *EffectIBoost {
	return &EffectIBoost{amount: amount}
}

// Desc returns the effect's description.
func (e *EffectIBoost) Desc() game.EffectDescription {
	return EffectDescIBoost
}

func (e *EffectIBoost) Amount() int {
	return e.amount
}

func (e *EffectIBoost) Increase() {
	e.amount += 5
}

var EffectDescSLayers = game.EffectDescription{
	Name: "S Layers",
	Type: game.EffectTypeBuff,
}

type EffectSLayers struct {
	threshold int
}

func NewEffectSLayers(threshold int) EffectSLayers {
	return EffectSLayers{threshold: threshold}
}

// Desc returns the effect's description.
func (e EffectSLayers) Desc() game.EffectDescription {
	return EffectDescSLayers
}

func (e EffectSLayers) Threshold() int {
	return e.threshold
}

// ModifyTakenDamage returns the modified amount of damage based on provided amount of damage and damage colour.
func (e EffectSLayers) ModifyTakenDamage(dmg int, colour game.Colour) int {
	if dmg <= e.threshold {
		return 0
	}

	return dmg
}

var EffectDescLastChance = game.EffectDescription{
	Name: "Last Chance",
	Type: game.EffectTypeState,
}

// If you survive your opponent's next turn, fully heals you.
type EffectLastChance struct {
}

// Desc returns the effect's description.
func (e EffectLastChance) Desc() game.EffectDescription {
	return EffectDescLastChance
}
