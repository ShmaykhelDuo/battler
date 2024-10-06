package milana_test

import (
	"fmt"
	"testing"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/milana"
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

func TestEffectStolenHP_Increase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		initAmount int
		incrAmount int
		amount     int
	}{
		{
			initAmount: 5,
			incrAmount: 5,
			amount:     10,
		},
		{
			initAmount: 4,
			incrAmount: 3,
			amount:     7,
		},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%d+%d=%d", tt.initAmount, tt.incrAmount, tt.amount)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			eff := milana.NewEffectStolenHP(tt.initAmount)
			eff.Increase(tt.incrAmount)
			assert.Equal(t, tt.amount, eff.Amount())
		})
	}
}

func TestEffectStolenHP_Decrease(t *testing.T) {
	t.Parallel()

	tests := []struct {
		initAmount int
		decrAmount int
		amount     int
	}{
		{
			initAmount: 10,
			decrAmount: 5,
			amount:     5,
		},
		{
			initAmount: 7,
			decrAmount: 3,
			amount:     4,
		},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%d-%d=%d", tt.initAmount, tt.decrAmount, tt.amount)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			eff := milana.NewEffectStolenHP(tt.initAmount)
			eff.Decrease(tt.decrAmount)
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
