package gamestate

import "github.com/ShmaykhelDuo/battler/internal/game/characters/euphoria"

type EffectEuphoriaEuphoricSource struct {
	Amount int `json:"amount"`
}

func NewEffectEuphoriaEuphoricSource(e euphoria.EffectEuphoricSource) EffectEuphoriaEuphoricSource {
	return EffectEuphoriaEuphoricSource{
		Amount: e.Amount(),
	}
}

type EffectEuphoriaUltimateEarly struct {
	Amount int `json:"amount"`
}

func NewEffectEuphoriaUltimateEarly(e *euphoria.EffectUltimateEarly) EffectEuphoriaUltimateEarly {
	return EffectEuphoriaUltimateEarly{
		Amount: e.Amount(),
	}
}

type EffectEuphoriaEuphoricHeal struct {
}
