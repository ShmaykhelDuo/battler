package minimax_test

import (
	"fmt"
	"testing"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/bot/minimax"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/milana"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/ruby"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/storyteller"
	"github.com/stretchr/testify/assert"
)

func TestMiniMax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		c        game.CharacterData
		opp      game.CharacterData
		depth    int
		score    int
		strategy []int
	}{
		{
			name:     "StorytellerRuby1",
			c:        storyteller.CharacterStoryteller,
			opp:      ruby.CharacterRuby,
			depth:    1,
			score:    0,
			strategy: []int{0, 1},
		},
		{
			name:     "StorytellerRuby2",
			c:        storyteller.CharacterStoryteller,
			opp:      ruby.CharacterRuby,
			depth:    2,
			score:    -4,
			strategy: []int{0, 0, 0, 1},
		},
	}

	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(tt.c)
			opp := game.NewCharacter(tt.opp)

			turnState := game.TurnState{
				TurnNum:      1,
				IsGoingFirst: true,
			}

			score, strategy := minimax.MiniMax(c, opp, turnState, 1, tt.depth, false)
			assert.Equal(t, tt.score, score, "score")
			assert.Equal(t, tt.strategy, strategy)
		})
	}
}

func ExampleMiniMax() {
	c := game.NewCharacter(ruby.CharacterRuby)
	opp := game.NewCharacter(milana.CharacterMilana)

	opp.Skills()[0].Use(c, game.TurnState{
		TurnNum:      1,
		IsGoingFirst: true,
	})
	c.Skills()[0].Use(opp, game.TurnState{
		TurnNum:      1,
		IsGoingFirst: false,
	})
	opp.Skills()[2].Use(c, game.TurnState{
		TurnNum:      2,
		IsGoingFirst: true,
	})

	turnState := game.TurnState{
		TurnNum:      2,
		IsGoingFirst: false,
	}
	minimax.MiniMax(c, opp, turnState, 1, 8, false)
}

func runMiniMax(b *testing.B, depth int) {
	c := game.NewCharacter(storyteller.CharacterStoryteller)
	opp := game.NewCharacter(ruby.CharacterRuby)

	turnState := game.TurnState{
		TurnNum:      1,
		IsGoingFirst: true,
	}

	for i := 0; i < b.N; i++ {
		clonedC, clonedOpp := game.Clone(c, opp)
		minimax.MiniMax(clonedC, clonedOpp, turnState, 1, depth, false)
	}
}

func BenchmarkMiniMax1(b *testing.B) {
	for i := range 5 {
		b.Run(fmt.Sprintf("%d", i+1), func(b *testing.B) {
			runMiniMax(b, i+1)
		})
	}
}
