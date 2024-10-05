package euphoria_test

import (
	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/euphoria"
)

func euphoricSourceAmount(c *game.Character) int {
	eff := c.Effect(euphoria.EffectDescEuphoricSource)
	e, ok := eff.(*euphoria.EffectEuphoricSource)
	if !ok {
		return 0
	}

	return e.Amount()
}
