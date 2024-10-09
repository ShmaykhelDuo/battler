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

var EffectDescMintMist = game.EffectDescription{
	Name: "Mint Mist",
}

// Your opponent can't use debuffs on you. Your Royal Move and Composure become stronger.
type EffectMintMist struct {
	common.DurationExpirable
}

func NewEffectMintMist(gameCtx game.Context) EffectMintMist {
	return EffectMintMist{
		DurationExpirable: common.NewDurationExpirable(gameCtx.AddTurns(2, false)),
	}
}

// Desc returns the effect's description.
func (e EffectMintMist) Desc() game.EffectDescription {
	return EffectDescMintMist
}

// IsEffectAllowed reports whether the effect can be applied to a character.
func (e EffectMintMist) IsEffectAllowed(eff game.Effect) bool {
	return eff.Desc().Type != game.EffectTypeProhibiting
}
