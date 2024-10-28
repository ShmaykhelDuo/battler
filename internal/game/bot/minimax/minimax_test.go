package minimax_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/bot/minimax"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/milana"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/ruby"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/storyteller"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiniMax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		c        game.CharacterData
		opp      game.CharacterData
		depth    int
		score    int
		strategy match.SkillLog
	}{
		{
			name:  "StorytellerRuby1",
			c:     storyteller.CharacterStoryteller,
			opp:   ruby.CharacterRuby,
			depth: 1,
			score: 0,
			// strategy: []int{0, 1},
			strategy: match.SkillLog{
				game.NewTurnState(1).WithGoingFirst(true):  []int{storyteller.SkillYourNumberIndex},
				game.NewTurnState(1).WithGoingFirst(false): []int{ruby.SkillRageIndex},
			},
		},
		{
			name:  "StorytellerRuby2",
			c:     storyteller.CharacterStoryteller,
			opp:   ruby.CharacterRuby,
			depth: 2,
			score: -4,
			// strategy: []int{0, 0, 0, 1},
			strategy: match.SkillLog{
				game.NewTurnState(1).WithGoingFirst(true):  []int{storyteller.SkillYourNumberIndex},
				game.NewTurnState(1).WithGoingFirst(false): []int{ruby.SkillDanceIndex},
				game.NewTurnState(2).WithGoingFirst(true):  []int{storyteller.SkillYourNumberIndex},
				game.NewTurnState(2).WithGoingFirst(false): []int{ruby.SkillRageIndex},
			},
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

			state := match.GameState{
				Character:  c,
				Opponent:   opp,
				TurnState:  turnState,
				SkillsLeft: 1,
				SkillLog:   make(match.SkillLog),
				PlayerTurn: true,
				AsOpp:      false,
			}
			score, strategy, _, err := minimax.MiniMax(context.Background(), state, tt.depth)
			require.NoError(t, err, "error")
			assert.Equal(t, tt.score, score, "score")
			assert.Equal(t, tt.strategy, strategy, "strategy")
		})
	}
}

func ExampleMiniMax() {
	c := game.NewCharacter(ruby.CharacterRuby)
	opp := game.NewCharacter(milana.CharacterMilana)

	opp.Skills()[milana.SkillRoyalMoveIndex].Use(c, game.TurnState{
		TurnNum:      1,
		IsGoingFirst: true,
	})
	c.Skills()[ruby.SkillDanceIndex].Use(opp, game.TurnState{
		TurnNum:      1,
		IsGoingFirst: false,
	})
	opp.Skills()[milana.SkillMintMistIndex].Use(c, game.TurnState{
		TurnNum:      2,
		IsGoingFirst: true,
	})

	turnState := game.TurnState{
		TurnNum:      2,
		IsGoingFirst: false,
	}

	state := match.GameState{
		Character:  c,
		Opponent:   opp,
		TurnState:  turnState,
		SkillsLeft: 1,
		SkillLog: match.SkillLog{
			game.NewTurnState(1).WithGoingFirst(true):  []int{milana.SkillRoyalMoveIndex},
			game.NewTurnState(1).WithGoingFirst(false): []int{ruby.SkillDanceIndex},
			game.NewTurnState(2).WithGoingFirst(true):  []int{milana.SkillMintMistIndex},
		},
		PlayerTurn: true,
		AsOpp:      false,
	}

	minimax.MiniMax(context.Background(), state, 8)
}

func runMiniMax(b *testing.B, depth int) {
	c := game.NewCharacter(storyteller.CharacterStoryteller)
	opp := game.NewCharacter(ruby.CharacterRuby)

	turnState := game.TurnState{
		TurnNum:      1,
		IsGoingFirst: true,
	}

	for i := 0; i < b.N; i++ {
		clonedC := c.Clone()
		clonedOpp := opp.Clone()

		state := match.GameState{
			Character:  clonedC,
			Opponent:   clonedOpp,
			TurnState:  turnState,
			SkillsLeft: 1,
			SkillLog:   make(match.SkillLog),
			PlayerTurn: true,
			AsOpp:      false,
		}

		minimax.MiniMax(context.Background(), state, depth)
	}
}

func BenchmarkMiniMax1(b *testing.B) {
	for i := range 5 {
		b.Run(fmt.Sprintf("%d", i+1), func(b *testing.B) {
			runMiniMax(b, i+1)
		})
	}
}
