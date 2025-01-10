package gamestate

import "github.com/ShmaykhelDuo/battler/internal/game/characters/z89"

type EffectZ89UltimateSlow struct {
	Amount int `json:"amount"`
}

func NewEffectZ89UltimateSlow(e *z89.EffectUltimateSlow) EffectZ89UltimateSlow {
	return EffectZ89UltimateSlow{
		Amount: e.Amount(),
	}
}
