package milana_test

import (
	"fmt"
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/milana"
	"github.com/stretchr/testify/assert"
)

func TestNewEffectStolenHP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		amount int
	}{
		{
			amount: 10,
		},
		{
			amount: 50,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.amount), func(t *testing.T) {
			t.Parallel()

			eff := milana.NewEffectStolenHP(tt.amount)

			assert.Equal(t, tt.amount, eff.Amount())
		})
	}
}

type dummyEffect struct {
	desc game.EffectDescription
}

// Desc returns the effect's description.
func (e dummyEffect) Desc() game.EffectDescription {
	return e.desc
}

// Clone returns a clone of the effect.
func (e dummyEffect) Clone() game.Effect {
	return e
}

func TestEffectMintMist_IsEffectAllowed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		eff       game.Effect
		isAllowed bool
	}{
		{
			name: "NotProhibiting",
			eff: dummyEffect{
				desc: game.EffectDescription{
					Type: game.EffectTypeBasic,
				},
			},
			isAllowed: true,
		},
		{
			name: "Prohibiting",
			eff: dummyEffect{
				desc: game.EffectDescription{
					Type: game.EffectTypeProhibiting,
				},
			},
			isAllowed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			eff := milana.EffectMintMist{}

			isAllowed := eff.IsEffectAllowed(tt.eff)
			assert.Equal(t, tt.isAllowed, isAllowed)
		})
	}
}
