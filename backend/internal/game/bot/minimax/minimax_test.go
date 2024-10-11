package minimax_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/bot/minimax"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/speed"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/storyteller"
)

func TestMiniMax(t *testing.T) {
	t.Parallel()

	c := game.NewCharacter(storyteller.CharacterStoryteller)
	opp := game.NewCharacter(speed.CharacterSpeed)

	gameCtx := game.Context{
		TurnNum:      1,
		IsGoingFirst: true,
	}

	minimax.MiniMax(c, opp, gameCtx, 1, 5, false)
}
