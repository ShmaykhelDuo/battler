package milana

import (
	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/common"
)

var EffectDescStolenHP = game.EffectDescription{
	Name: "Stolen HP",
}

// Damage dealt by Royal Move. You can spend it on Composure heal or Pride damage.
type EffectStolenHP struct {
	*common.Collectible
}

func NewEffectStolenHP(amount int) EffectStolenHP {
	return EffectStolenHP{
		Collectible: common.NewCollectible(amount),
	}
}

// Desc returns the effect's description.
func (e EffectStolenHP) Desc() game.EffectDescription {
	return EffectDescStolenHP
}

// Clone returns a clone of the effect.
func (e EffectStolenHP) Clone() game.Effect {
	return NewEffectStolenHP(e.Amount())
}

var EffectDescMintMist = game.EffectDescription{
	Name: "Mint Mist",
}

// Your opponent can't use debuffs on you. Your Royal Move and Composure become stronger.
type EffectMintMist struct {
	common.DurationExpirable
}

func NewEffectMintMist(turnState game.TurnState) EffectMintMist {
	return EffectMintMist{
		DurationExpirable: common.NewDurationExpirable(turnState.AddTurns(2, false)),
	}
}

// Desc returns the effect's description.
func (e EffectMintMist) Desc() game.EffectDescription {
	return EffectDescMintMist
}

// Clone returns a clone of the effect.
func (e EffectMintMist) Clone() game.Effect {
	return e
}

// IsEffectAllowed reports whether the effect can be applied to a character.
func (e EffectMintMist) IsEffectAllowed(eff game.Effect) bool {
	return eff.Desc().Type != game.EffectTypeProhibiting
}
