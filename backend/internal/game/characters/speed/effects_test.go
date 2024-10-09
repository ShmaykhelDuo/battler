package speed_test

import (
	"fmt"
	"maps"
	"testing"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/speed"
	"github.com/stretchr/testify/assert"
)

func TestNewEffectGreenTokens(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		number int
	}{
		{
			name:   "1",
			number: 1,
		},
		{
			name:   "15",
			number: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			eff := speed.NewEffectGreenTokens(tt.number)

			assert.Equal(t, speed.EffectDescGreenTokens, eff.Desc(), "description")
			assert.Equal(t, tt.number, eff.Amount(), "number")
		})
	}
}

func TestNewEffectBlackTokens(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		number int
	}{
		{
			name:   "1",
			number: 1,
		},
		{
			name:   "15",
			number: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			eff := speed.NewEffectBlackTokens(tt.number)

			assert.Equal(t, speed.EffectDescBlackTokens, eff.Desc(), "description")
			assert.Equal(t, tt.number, eff.Amount(), "number")
		})
	}
}

func TestNewEffectDamageReduced(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		amount int
	}{
		{
			name:   "One",
			amount: 1,
		},
		{
			name:   "Five",
			amount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			eff := speed.NewEffectDamageReduced(tt.amount)

			assert.Equal(t, tt.amount, eff.Amount())
		})
	}
}

func TestEffectDamageReduced_Increase(t *testing.T) {
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

			eff := speed.NewEffectDamageReduced(tt.initAmount)
			eff.Increase(tt.incrAmount)
			assert.Equal(t, tt.amount, eff.Amount())
		})
	}
}

func TestEffectDamageReduced_ModifyTakenDamage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		amount int
		in     int
		colour game.Colour
		out    int
	}{
		{
			name:   "1",
			amount: 5,
			in:     10,
			colour: game.ColourGreen,
			out:    5,
		},
		{
			name:   "2",
			amount: 3,
			in:     32,
			colour: game.ColourBlue,
			out:    29,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			eff := speed.NewEffectDamageReduced(tt.amount)

			out := eff.ModifyTakenDamage(tt.in, tt.colour)
			assert.Equal(t, tt.out, out)
		})
	}
}

func TestEffectDamageReduced_HasExpired(t *testing.T) {
	t.Parallel()

	t.Run("NotExpiredByDefault", func(t *testing.T) {
		t.Parallel()

		eff := speed.NewEffectDamageReduced(5)
		assert.False(t, eff.HasExpired(game.Context{}))
	})

	t.Run("ExpiredAfterAttack", func(t *testing.T) {
		t.Parallel()

		eff := speed.NewEffectDamageReduced(5)
		eff.ModifyTakenDamage(25, game.ColourNone)

		assert.True(t, eff.HasExpired(game.Context{}))
	})
}

func TestEffectDefenceReduced_ModifyDefences(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   map[game.Colour]int
		out  map[game.Colour]int
	}{
		{
			name: "1",
			in:   map[game.Colour]int{},
			out: map[game.Colour]int{
				game.ColourGreen: -1,
			},
		},
		{
			name: "2",
			in: map[game.Colour]int{
				game.ColourBlue:  2,
				game.ColourGreen: 4,
			},
			out: map[game.Colour]int{
				game.ColourBlue:  2,
				game.ColourGreen: 3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			eff := speed.EffectDefenceReduced{}

			res := maps.Clone(tt.in)
			eff.ModifyDefences(res)
			assert.Equal(t, tt.out, res)
		})
	}
}

func TestEffectSpedUp_SkillsPerTurn(t *testing.T) {
	t.Parallel()

	eff := speed.EffectSpedUp{}

	assert.Equal(t, 2, eff.SkillsPerTurn())
}

func TestEffectSpedUp_IsSkillAvailable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		skill       game.SkillData
		isAvailable bool
	}{
		{
			name: "Ultimate",
			skill: game.SkillData{
				Desc: game.SkillDescription{
					IsUltimate: true,
				},
			},
			isAvailable: false,
		},
		{
			name: "NotUltimate",
			skill: game.SkillData{
				Desc: game.SkillDescription{
					IsUltimate: false,
				},
			},
			isAvailable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			eff := speed.EffectSpedUp{}

			data := game.CharacterData{}
			c := game.NewCharacter(data)
			s := game.NewSkill(c, tt.skill)

			isAvailable := eff.IsSkillAvailable(s)
			assert.Equal(t, tt.isAvailable, isAvailable)
		})
	}
}
