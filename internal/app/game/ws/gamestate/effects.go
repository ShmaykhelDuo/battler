package gamestate

import (
	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/euphoria"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/milana"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/ruby"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/speed"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/storyteller"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/structure"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/z89"
)

type Effect struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func NewEffect(e game.Effect, turnState game.TurnState) Effect {
	res := Effect{
		Type: e.Desc().Name,
	}

	switch eff := e.(type) {
	case storyteller.EffectCannotUse:
		res.Data = NewEffectStorytellerCannotUse(eff, turnState)
	case storyteller.EffectControlled:
		res.Data = NewEffectStorytellerControlled(eff, turnState)
	case *z89.EffectUltimateSlow:
		res.Data = NewEffectZ89UltimateSlow(eff)
	case euphoria.EffectEuphoricSource:
		res.Data = NewEffectEuphoriaEuphoricSource(eff)
	case *euphoria.EffectUltimateEarly:
		res.Data = NewEffectEuphoriaUltimateEarly(eff)
	case euphoria.EffectEuphoricHeal:
		res.Data = EffectEuphoriaEuphoricHeal{}
	case ruby.EffectDoubleDamage:
		res.Data = NewEffectRubyDoubleDamage(eff, turnState)
	case ruby.EffectCannotHeal:
		res.Data = NewEffectRubyCannotHeal(eff, turnState)
	case speed.EffectGreenTokens:
		res.Data = NewEffectSpeedGreenTokens(eff)
	case speed.EffectBlackTokens:
		res.Data = NewEffectSpeedBlackTokens(eff)
	case *speed.EffectDamageReduced:
		res.Data = NewEffectSpeedDamageReduced(eff)
	case speed.EffectDefenceReduced:
		res.Data = EffectSpeedDefenceReduced{}
	case speed.EffectSpedUp:
		res.Data = NewEffectSpeedSpedUp(eff, turnState)
	case milana.EffectStolenHP:
		res.Data = NewEffectMilanaStolenHP(eff)
	case milana.EffectMintMist:
		res.Data = NewEffectMilanaMintMist(eff, turnState)
	case *structure.EffectIBoost:
		res.Data = NewEffectStructureIBoost(eff)
	case structure.EffectSLayers:
		res.Data = NewEffectStructureSLayers(eff, turnState)
	case structure.EffectLastChance:
		res.Data = NewEffectStructureLastChance(eff, turnState)
	}

	return res
}
