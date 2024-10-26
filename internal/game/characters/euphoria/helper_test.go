package euphoria_test

import (
	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/euphoria"
)

func euphoricSourceAmount(c *game.Character) int {
	e, ok := game.CharacterEffect[euphoria.EffectEuphoricSource](c, euphoria.EffectDescEuphoricSource)
	if !ok {
		return 0
	}

	return e.Amount()
}
