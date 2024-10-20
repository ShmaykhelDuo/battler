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

			e := storyteller.NewEffectCannotUse(game.TurnState{}, tt.colour)

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
		name          string
		initturnState game.TurnState
		turnState     game.TurnState
		hasExpired    bool
	}{
		{
			name: "ImmediatelyAfterWhenGoingFirst",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			turnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			hasExpired: false,
		},
		{
			name: "ImmediatelyAfterTurnEndWhenGoingFirst",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			turnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: true,
				IsTurnEnd:    true,
			},
			hasExpired: false,
		},
		{
			name: "OpponentTurnWhenGoingFirst",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			turnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			hasExpired: false,
		},
		{
			name: "AfterOpponentTurnWhenGoingFirst",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			turnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
				IsTurnEnd:    true,
			},
			hasExpired: true,
		},
		{
			name: "ImmediatelyAfterWhenGoingSecond",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			turnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			hasExpired: false,
		},
		{
			name: "ImmediatelyAfterTurnEndWhenGoingSecond",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			turnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
				IsTurnEnd:    true,
			},
			hasExpired: false,
		},
		{
			name: "OpponentTurnWhenGoingSecond",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			turnState: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			hasExpired: false,
		},
		{
			name: "AfterOpponentTurnWhenGoingSecond",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			turnState: game.TurnState{
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

			eff := storyteller.NewEffectCannotUse(tt.initturnState, game.ColourNone)

			hasExpired := eff.HasExpired(tt.turnState)
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
		name          string
		initturnState game.TurnState
		turnState     game.TurnState
		hasExpired    bool
	}{
		{
			name: "ImmediatelyAfterWhenGoingFirst",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			turnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			hasExpired: false,
		},
		{
			name: "ImmediatelyAfterTurnEndWhenGoingFirst",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			turnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: true,
				IsTurnEnd:    true,
			},
			hasExpired: false,
		},
		{
			name: "OpponentTurnWhenGoingFirst",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			turnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			hasExpired: false,
		},
		{
			name: "AfterOpponentTurnWhenGoingFirst",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			turnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
				IsTurnEnd:    true,
			},
			hasExpired: true,
		},
		{
			name: "ImmediatelyAfterWhenGoingSecond",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			turnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			hasExpired: false,
		},
		{
			name: "ImmediatelyAfterTurnEndWhenGoingSecond",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			turnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
				IsTurnEnd:    true,
			},
			hasExpired: false,
		},
		{
			name: "OpponentTurnWhenGoingSecond",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			turnState: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			hasExpired: false,
		},
		{
			name: "AfterOpponentTurnWhenGoingSecond",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			turnState: game.TurnState{
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

			eff := storyteller.NewEffectControlled(tt.initturnState)

			hasExpired := eff.HasExpired(tt.turnState)
			assert.Equal(t, tt.hasExpired, hasExpired)
		})
	}
}
