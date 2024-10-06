package milana

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

var EffectDescStolenHP = game.EffectDescription{
	Name: "Stolen HP",
}

// Damage dealt by Royal Move. You can spend it on Composure heal or Pride damage.
type EffectStolenHP struct {
	amount int
}

func NewEffectStolenHP(amount int) *EffectStolenHP {
	return &EffectStolenHP{amount: amount}
}

// Desc returns the effect's description.
func (e *EffectStolenHP) Desc() game.EffectDescription {
	return EffectDescStolenHP
}

func (e *EffectStolenHP) Amount() int {
	return e.amount
}

func (e *EffectStolenHP) Increase(amount int) {
	e.amount += amount
}

func (e *EffectStolenHP) Decrease(amount int) {
	e.amount -= amount
}

var EffectDescMintMist = game.EffectDescription{
	Name: "Mint Mist",
}

// Your opponent can't use debuffs on you. Your Royal Move and Composure become stronger.
type EffectMintMist struct {
}

// Desc returns the effect's description.
func (e EffectMintMist) Desc() game.EffectDescription {
	return EffectDescMintMist
}

// IsEffectAllowed reports whether the effect can be applied to a character.
func (e EffectMintMist) IsEffectAllowed(eff game.Effect) bool {
	return eff.Desc().Type != game.EffectTypeProhibiting
}
