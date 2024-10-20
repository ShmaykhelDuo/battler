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
			name: "NextTurnWhenGoingFirst",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			turnState: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			hasExpired: false,
		},
		{
			name: "SecondTurnWhenGoingFirst",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			turnState: game.TurnState{
				TurnNum:      6,
				IsGoingFirst: true,
			},
			hasExpired: false,
		},
		{
			name: "AfterSecondTurnWhenGoingFirst",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: true,
			},
			turnState: game.TurnState{
				TurnNum:      6,
				IsGoingFirst: true,
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
			name: "NextTurnWhenGoingSecond",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			turnState: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: false,
			},
			hasExpired: false,
		},
		{
			name: "SecondTurnWhenGoingSecond",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			turnState: game.TurnState{
				TurnNum:      6,
				IsGoingFirst: false,
			},
			hasExpired: false,
		},
		{
			name: "AfterSecondTurnWhenGoingSecond",
			initturnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			turnState: game.TurnState{
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

			eff := ruby.NewEffectDoubleDamage(tt.initturnState)

			hasExpired := eff.HasExpired(tt.turnState)
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
