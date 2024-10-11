package minimax_test

import (
	"fmt"
	"testing"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/bot/minimax"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/ruby"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/storyteller"
)

func TestMiniMax(t *testing.T) {
	t.Parallel()

	c := game.NewCharacter(storyteller.CharacterStoryteller)
	opp := game.NewCharacter(ruby.CharacterRuby)

	gameCtx := game.Context{
		TurnNum:      1,
		IsGoingFirst: true,
	}

	minimax.MiniMax(c, opp, gameCtx, 1, 3, false)
}

func runMiniMax(b *testing.B, depth int) {
	c := game.NewCharacter(storyteller.CharacterStoryteller)
	opp := game.NewCharacter(ruby.CharacterRuby)

	gameCtx := game.Context{
		TurnNum:      1,
		IsGoingFirst: true,
	}

	for i := 0; i < b.N; i++ {
		clonedC, clonedOpp := game.Clone(c, opp)
		minimax.MiniMax(clonedC, clonedOpp, gameCtx, 1, depth, false)
	}
}

func BenchmarkMiniMax1(b *testing.B) {
	for i := range 5 {
		b.Run(fmt.Sprintf("%d", i+1), func(b *testing.B) {
			runMiniMax(b, i+1)
		})
	}
}
