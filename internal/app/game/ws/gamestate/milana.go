package gamestate

import (
	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/milana"
)

type EffectMilanaStolenHP struct {
	Amount int `json:"amount"`
}

func NewEffectMilanaStolenHP(e milana.EffectStolenHP) EffectMilanaStolenHP {
	return EffectMilanaStolenHP{
		Amount: e.Amount(),
	}
}

type EffectMilanaMintMist struct {
	TurnsLeft int `json:"turns_left"`
}

func NewEffectMilanaMintMist(e milana.EffectMintMist, turnState game.TurnState) EffectMilanaMintMist {
	return EffectMilanaMintMist{
		TurnsLeft: e.TurnsLeft(turnState),
	}
}
