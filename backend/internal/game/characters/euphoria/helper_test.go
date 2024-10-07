package euphoria_test

import (
	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/euphoria"
)

func euphoricSourceAmount(c *game.Character) int {
	e, ok := game.CharacterEffect[*euphoria.EffectEuphoricSource](c)
	if !ok {
		return 0
	}

	return e.Amount()
}
