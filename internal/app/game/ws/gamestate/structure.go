package gamestate

import (
	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/structure"
)

type EffectStructureIBoost struct {
	Amount int `json:"amount"`
}

func NewEffectStructureIBoost(e *structure.EffectIBoost) EffectStructureIBoost {
	return EffectStructureIBoost{
		Amount: e.Amount(),
	}
}

type EffectStructureSLayers struct {
	TurnsLeft int `json:"turns_left"`
	Threshold int `json:"threshold"`
}

func NewEffectStructureSLayers(e structure.EffectSLayers, turnState game.TurnState) EffectStructureSLayers {
	return EffectStructureSLayers{
		TurnsLeft: e.TurnsLeft(turnState),
		Threshold: e.Threshold(),
	}
}

type EffectStructureLastChance struct {
	TurnsLeft int `json:"turns_left"`
}

func NewEffectStructureLastChance(e structure.EffectLastChance, turnState game.TurnState) EffectStructureLastChance {
	return EffectStructureLastChance{
		TurnsLeft: e.TurnsLeft(turnState),
	}
}
