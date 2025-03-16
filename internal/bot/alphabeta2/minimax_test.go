package alphabeta2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/bot/alphabeta2"
	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/milana"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/ruby"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/storyteller"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMinimax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		c     *game.CharacterData
		opp   *game.CharacterData
		depth int
		score int
		skill int
	}{
		{
			name:  "StorytellerRuby1",
			c:     storyteller.CharacterStoryteller,
			opp:   ruby.CharacterRuby,
			depth: 1,
			score: 0,
			skill: storyteller.SkillYourNumberIndex,
		},
		{
			name:  "StorytellerRuby2",
			c:     storyteller.CharacterStoryteller,
			opp:   ruby.CharacterRuby,
			depth: 2,
			score: -4,
			skill: storyteller.SkillYourNumberIndex,
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
			res, err := alphabeta2.Minimax(context.Background(), state, tt.depth)
			require.NoError(t, err, "error")
			assert.Equal(t, tt.score, res.Score, "score")
			assert.Equal(t, tt.skill, res.Skill, "skill")
		})
	}
}

func runMiniMax(b *testing.B, depth int) {
	c := game.NewCharacter(ruby.CharacterRuby)
	opp := game.NewCharacter(milana.CharacterMilana)

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

		alphabeta2.Minimax(context.Background(), state, depth)
	}
}

func BenchmarkMiniMax(b *testing.B) {
	for i := range 10 {
		b.Run(fmt.Sprintf("depth=%d", i+1), func(b *testing.B) {
			runMiniMax(b, i+1)
		})
	}
}
