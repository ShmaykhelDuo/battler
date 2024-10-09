package ruby_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/ruby"
	"github.com/stretchr/testify/assert"
)

func TestEffectDoubleDamage_ModifyDealtDamage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		dmg    int
		colour game.Colour
		out    int
	}{
		{
			name:   "Test1",
			dmg:    15,
			colour: game.ColourNone,
			out:    30,
		},
		{
			name:   "Test2",
			dmg:    24,
			colour: game.ColourRed,
			out:    48,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			eff := ruby.EffectDoubleDamage{}
			out := eff.ModifyDealtDamage(tt.dmg, tt.colour)
			assert.Equal(t, tt.out, out)
		})
	}
}

func TestEffectDoubleDamage_HasExpired(t *testing.T) {
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
			name: "NextTurnWhenGoingFirst",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			gameCtx: game.Context{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			hasExpired: false,
		},
		{
			name: "SecondTurnWhenGoingFirst",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			gameCtx: game.Context{
				TurnNum:      6,
				IsGoingFirst: true,
			},
			hasExpired: false,
		},
		{
			name: "AfterSecondTurnWhenGoingFirst",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			gameCtx: game.Context{
				TurnNum:      6,
				IsGoingFirst: true,
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
			name: "NextTurnWhenGoingSecond",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			gameCtx: game.Context{
				TurnNum:      5,
				IsGoingFirst: false,
			},
			hasExpired: false,
		},
		{
			name: "SecondTurnWhenGoingSecond",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			gameCtx: game.Context{
				TurnNum:      6,
				IsGoingFirst: false,
			},
			hasExpired: false,
		},
		{
			name: "AfterSecondTurnWhenGoingSecond",
			initGameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			gameCtx: game.Context{
				TurnNum:      6,
				IsGoingFirst: false,
				IsTurnEnd:    true,
			},
			hasExpired: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			eff := ruby.NewEffectDoubleDamage(tt.initGameCtx)

			hasExpired := eff.HasExpired(tt.gameCtx)
			assert.Equal(t, tt.hasExpired, hasExpired)
		})
	}
}

func TestEffectCannotHeal_IsHealAllowed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		heal int
	}{
		{
			name: "Test1",
			heal: 14,
		},
		{
			name: "Test1",
			heal: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			eff := ruby.EffectCannotHeal{}

			assert.False(t, eff.IsHealAllowed(tt.heal))
		})
	}
}
