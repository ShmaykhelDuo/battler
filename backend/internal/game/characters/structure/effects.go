package structure

import (
	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/common"
)

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
	common.DurationExpirable
	threshold int
}

func NewEffectSLayers(gameCtx game.Context, threshold int) EffectSLayers {
	return EffectSLayers{
		DurationExpirable: common.NewDurationExpirable(gameCtx.AddTurns(0, true)),
		threshold:         threshold,
	}
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
	common.DurationExpirable
	healCtx game.Context
}

func NewEffectLastChance(gameCtx game.Context) EffectLastChance {
	endCtx := gameCtx.AddTurns(0, true)

	return EffectLastChance{
		DurationExpirable: common.NewDurationExpirable(endCtx),
		healCtx:           endCtx,
	}
}

// Desc returns the effect's description.
func (e EffectLastChance) Desc() game.EffectDescription {
	return EffectDescLastChance
}

// OnTurnEnd executes the end-of-turn action.
func (e EffectLastChance) OnTurnEnd(c *game.Character, opp *game.Character, gameCtx game.Context) {
	if gameCtx.IsAfter(e.healCtx) {
		c.Heal(c.MaxHP())
	}
}
