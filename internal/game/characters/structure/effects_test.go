package structure_test

import (
	"fmt"
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/structure"
	"github.com/stretchr/testify/assert"
)

func TestNewEffectIBoost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		amount int
	}{
		{
			amount: 5,
		},
		{
			amount: 15,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.amount), func(t *testing.T) {
			t.Parallel()

			eff := structure.NewEffectIBoost(tt.amount)

			assert.Equal(t, tt.amount, eff.Amount())
		})
	}
}

func TestEffectIBoost_Increase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		initAmount int
		amount     int
	}{
		{
			initAmount: 5,
			amount:     10,
		},
		{
			initAmount: 10,
			amount:     15,
		},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%d+5=%d", tt.initAmount, tt.amount)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			eff := structure.NewEffectIBoost(tt.initAmount)

			eff.Increase()
			assert.Equal(t, tt.amount, eff.Amount())
		})
	}
}

func TestNewEffectSLayers(t *testing.T) {
	t.Parallel()

	threshold := 5

	eff := structure.NewEffectSLayers(game.TurnState{}, threshold)
	assert.Equal(t, threshold, eff.Threshold())
}

func TestEffectSLayers_ModifyTakenDamage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		threshold int
		dmg       int
		out       int
	}{
		{
			name:      "AboveThreshold",
			threshold: 15,
			dmg:       16,
			out:       16,
		},
		{
			name:      "AtThreshold",
			threshold: 15,
			dmg:       15,
			out:       0,
		},
		{
			name:      "BelowThreshold",
			threshold: 15,
			dmg:       13,
			out:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			eff := structure.NewEffectSLayers(game.TurnState{}, tt.threshold)

			out := eff.ModifyTakenDamage(tt.dmg, game.ColourNone)
			assert.Equal(t, tt.out, out)
		})
	}
}

func TestEffectLastChance_OnTurnEnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		prevDmg       int
		initturnState game.TurnState
		turnState     game.TurnState
		hp            int
	}{
		{
			name:    "SameTurn",
			prevDmg: 36,
			initturnState: game.TurnState{
				TurnNum:      7,
				IsGoingFirst: true,
			},
			turnState: game.TurnState{
				TurnNum:      7,
				IsGoingFirst: true,
				IsTurnEnd:    true,
			},
			hp: 83,
		},
		{
			name:    "OpponentsTurnGoingFirst",
			prevDmg: 36,
			initturnState: game.TurnState{
				TurnNum:      7,
				IsGoingFirst: true,
			},
			turnState: game.TurnState{
				TurnNum:      7,
				IsGoingFirst: false,
				IsTurnEnd:    true,
			},
			hp: 119,
		},
		{
			name:    "OpponentsTurnGoingSecond",
			prevDmg: 36,
			initturnState: game.TurnState{
				TurnNum:      7,
				IsGoingFirst: false,
			},
			turnState: game.TurnState{
				TurnNum:      8,
				IsGoingFirst: true,
				IsTurnEnd:    true,
			},
			hp: 119,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(structure.CharacterStructure)
			opp := game.NewCharacter(&game.CharacterData{})

			opp.Damage(c, tt.prevDmg, game.ColourNone)

			eff := structure.NewEffectLastChance(tt.initturnState)

			eff.OnTurnEnd(c, opp, tt.turnState)

			assert.Equal(t, tt.hp, c.HP())
		})
	}
}
