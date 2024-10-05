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
