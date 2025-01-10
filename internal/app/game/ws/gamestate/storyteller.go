package gamestate

import (
	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/storyteller"
)

type EffectStorytellerCannotUse struct {
	TurnsLeft int    `json:"turns_left"`
	Colour    Colour `json:"colour"`
}

func NewEffectStorytellerCannotUse(e storyteller.EffectCannotUse, turnState game.TurnState) EffectStorytellerCannotUse {
	return EffectStorytellerCannotUse{
		TurnsLeft: e.TurnsLeft(turnState),
		Colour:    NewColour(e.Colour()),
	}
}

type EffectStorytellerControlled struct {
	TurnsLeft int `json:"turns_left"`
}

func NewEffectStorytellerControlled(e storyteller.EffectControlled, turnState game.TurnState) EffectStorytellerControlled {
	return EffectStorytellerControlled{
		TurnsLeft: e.TurnsLeft(turnState),
	}
}
