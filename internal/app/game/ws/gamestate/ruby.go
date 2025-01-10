package gamestate

import (
	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/ruby"
)

type EffectRubyDoubleDamage struct {
	TurnsLeft int `json:"turns_left"`
}

func NewEffectRubyDoubleDamage(e ruby.EffectDoubleDamage, turnState game.TurnState) EffectRubyDoubleDamage {
	return EffectRubyDoubleDamage{
		TurnsLeft: e.TurnsLeft(turnState),
	}
}

type EffectRubyCannotHeal struct {
	TurnsLeft int `json:"turns_left"`
}

func NewEffectRubyCannotHeal(e ruby.EffectCannotHeal, turnState game.TurnState) EffectRubyCannotHeal {
	return EffectRubyCannotHeal{
		TurnsLeft: e.TurnsLeft(turnState),
	}
}
