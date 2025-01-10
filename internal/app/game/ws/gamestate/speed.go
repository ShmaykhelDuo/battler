package gamestate

import (
	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/speed"
)

type EffectSpeedGreenTokens struct {
	Amount int `json:"amount"`
}

func NewEffectSpeedGreenTokens(e speed.EffectGreenTokens) EffectSpeedGreenTokens {
	return EffectSpeedGreenTokens{
		Amount: e.Amount(),
	}
}

type EffectSpeedBlackTokens struct {
	Amount int `json:"amount"`
}

func NewEffectSpeedBlackTokens(e speed.EffectBlackTokens) EffectSpeedBlackTokens {
	return EffectSpeedBlackTokens{
		Amount: e.Amount(),
	}
}

type EffectSpeedDamageReduced struct {
	Amount int `json:"amount"`
}

func NewEffectSpeedDamageReduced(e *speed.EffectDamageReduced) EffectSpeedDamageReduced {
	return EffectSpeedDamageReduced{
		Amount: e.Amount(),
	}
}

type EffectSpeedDefenceReduced struct {
}

type EffectSpeedSpedUp struct {
	TurnsLeft int `json:"turns_left"`
}

func NewEffectSpeedSpedUp(e speed.EffectSpedUp, turnState game.TurnState) EffectSpeedSpedUp {
	return EffectSpeedSpedUp{
		TurnsLeft: e.TurnsLeft(turnState),
	}
}
