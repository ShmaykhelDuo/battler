package storyteller_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/storyteller"
	"github.com/stretchr/testify/assert"
)

func TestEffectCannotUse_IsSkillAvailable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		colour      game.Colour
		skill       game.SkillData
		isAvailable bool
	}{
		{
			name:   "ColoursMatch",
			colour: game.ColourViolet,
			skill: game.SkillData{
				Desc: game.SkillDescription{
					Colour: game.ColourViolet,
				},
			},
			isAvailable: false,
		},
		{
			name:   "ColoursNotMatch",
			colour: game.ColourViolet,
			skill: game.SkillData{
				Desc: game.SkillDescription{
					Colour: game.ColourGreen,
				},
			},
			isAvailable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := storyteller.NewEffectCannotUse(game.Context{}, tt.colour)

			data := game.CharacterData{}
			c := game.NewCharacter(data)
			s := game.NewSkill(c, tt.skill)

			isAvailable := e.IsSkillAvailable(s)
			assert.Equal(t, tt.isAvailable, isAvailable)
		})
	}
}

func TestEffectCannotUse_HasExpired(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		initGameCtx game.Context
		gameCtx     game.Context
		hasExpired  bool
	}{
		{
			name: "ImmediatelyAfterWhenGoingFirst",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			gameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			hasExpired: false,
		},
		{
			name: "ImmediatelyAfterTurnEndWhenGoingFirst",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			gameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: true,
				IsTurnEnd:    true,
			},
			hasExpired: false,
		},
		{
			name: "OpponentTurnWhenGoingFirst",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			gameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			hasExpired: false,
		},
		{
			name: "AfterOpponentTurnWhenGoingFirst",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			gameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
				IsTurnEnd:    true,
			},
			hasExpired: true,
		},
		{
			name: "ImmediatelyAfterWhenGoingSecond",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			gameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			hasExpired: false,
		},
		{
			name: "ImmediatelyAfterTurnEndWhenGoingSecond",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			gameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
				IsTurnEnd:    true,
			},
			hasExpired: false,
		},
		{
			name: "OpponentTurnWhenGoingSecond",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			gameCtx: game.Context{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			hasExpired: false,
		},
		{
			name: "AfterOpponentTurnWhenGoingSecond",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			gameCtx: game.Context{
				TurnNum:      5,
				IsGoingFirst: true,
				IsTurnEnd:    true,
			},
			hasExpired: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			eff := storyteller.NewEffectCannotUse(tt.initGameCtx, game.ColourNone)

			hasExpired := eff.HasExpired(tt.gameCtx)
			assert.Equal(t, tt.hasExpired, hasExpired)
		})
	}
}

func TestEffectControlled_HasTakenControl(t *testing.T) {
	t.Parallel()

	eff := storyteller.EffectControlled{}

	assert.True(t, eff.HasTakenControl())
}

func TestEffectControlled_HasExpired(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		initGameCtx game.Context
		gameCtx     game.Context
		hasExpired  bool
	}{
		{
			name: "ImmediatelyAfterWhenGoingFirst",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			gameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			hasExpired: false,
		},
		{
			name: "ImmediatelyAfterTurnEndWhenGoingFirst",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			gameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: true,
				IsTurnEnd:    true,
			},
			hasExpired: false,
		},
		{
			name: "OpponentTurnWhenGoingFirst",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			gameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			hasExpired: false,
		},
		{
			name: "AfterOpponentTurnWhenGoingFirst",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			gameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
				IsTurnEnd:    true,
			},
			hasExpired: true,
		},
		{
			name: "ImmediatelyAfterWhenGoingSecond",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			gameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			hasExpired: false,
		},
		{
			name: "ImmediatelyAfterTurnEndWhenGoingSecond",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			gameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
				IsTurnEnd:    true,
			},
			hasExpired: false,
		},
		{
			name: "OpponentTurnWhenGoingSecond",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			gameCtx: game.Context{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			hasExpired: false,
		},
		{
			name: "AfterOpponentTurnWhenGoingSecond",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			gameCtx: game.Context{
				TurnNum:      5,
				IsGoingFirst: true,
				IsTurnEnd:    true,
			},
			hasExpired: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			eff := storyteller.NewEffectControlled(tt.initGameCtx)

			hasExpired := eff.HasExpired(tt.gameCtx)
			assert.Equal(t, tt.hasExpired, hasExpired)
		})
	}
}
