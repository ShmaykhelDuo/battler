package minimax_test

import (
	"fmt"
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/bot/minimax"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/speed"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/storyteller"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

func BenchmarkMiniMaxBot(b *testing.B) {
	state := match.GameState{
		Character:  game.NewCharacter(storyteller.CharacterStoryteller),
		Opponent:   game.NewCharacter(speed.CharacterSpeed),
		TurnState:  game.NewTurnState(1),
		PlayerTurn: true,
	}

	for depth := 1; depth < 5; depth++ {
		name := fmt.Sprintf("Depth%d", depth)
		b.Run(name, func(b *testing.B) {
			bot := minimax.NewBot(depth)

			for i := 0; i < b.N; i++ {
				bot.SendState(state)
				bot.RequestSkill()
			}
		})
	}
}
