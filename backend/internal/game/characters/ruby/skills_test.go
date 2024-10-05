package ruby_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/ruby"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSkillDance_Use(t *testing.T) {
	t.Parallel()

	c := game.NewCharacter(ruby.CharacterRuby)

	data := game.CharacterData{}
	opp := game.NewCharacter(data)

	s := c.Skills()[0]

	gameCtx := game.Context{}
	err := s.Use(opp, gameCtx)
	require.NoError(t, err)

	eff := c.Effect(ruby.EffectDescDoubleDamage)
	require.NotNil(t, eff, "effect")

	_, ok := eff.(ruby.EffectDoubleDamage)
	require.True(t, ok, "effect type")
}

func TestSkillRage_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		oppData game.CharacterData
		gameCtx game.Context
		oppHP   int
	}{
		{
			name: "Opponent1",
			oppData: game.CharacterData{
				DefaultHP: 111,
				Defences: map[game.Colour]int{
					game.ColourRed: -1,
				},
			},
			gameCtx: game.Context{
				TurnNum: 3,
			},
			oppHP: 92,
		},
		{
			name: "Opponent2",
			oppData: game.CharacterData{
				DefaultHP: 119,
				Defences: map[game.Colour]int{
					game.ColourRed: -1,
				},
			},
			gameCtx: game.Context{
				TurnNum: 6,
			},
			oppHP: 106,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(ruby.CharacterRuby)
			opp := game.NewCharacter(tt.oppData)

			s := c.Skills()[1]
			err := s.Use(opp, tt.gameCtx)
			require.NoError(t, err)

			assert.Equal(t, tt.oppHP, opp.HP())
		})
	}
}

func assertEffectCannotHeal(t *testing.T, c *game.Character, name string) {
	t.Helper()

	eff := c.Effect(ruby.EffectDescCannotHeal)
	if !assert.NotNil(t, eff, "%s's effect", name) {
		return
	}

	_, ok := eff.(ruby.EffectCannotHeal)
	assert.True(t, ok, "%s's effect type", name)
}

func TestSkillStop_Use(t *testing.T) {
	t.Parallel()

	c := game.NewCharacter(ruby.CharacterRuby)

	data := game.CharacterData{}
	opp := game.NewCharacter(data)

	s := c.Skills()[2]

	gameCtx := game.Context{}
	err := s.Use(opp, gameCtx)
	require.NoError(t, err)

	assertEffectCannotHeal(t, c, "character")
	assertEffectCannotHeal(t, opp, "opponent")
}

func TestSkillExecute_IsAvailable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		oppData     game.CharacterData
		prevDmg     int
		effs        []game.Effect
		gameCtx     game.Context
		isAvailable bool
	}{
		{
			name: "AboveThreshold",
			oppData: game.CharacterData{
				DefaultHP: 111,
			},
			prevDmg:     99,
			isAvailable: false,
		},
		{
			name: "BelowThreshold",
			oppData: game.CharacterData{
				DefaultHP: 111,
			},
			prevDmg:     100,
			isAvailable: true,
		},
		{
			name: "AboveThresholdWithCannotHeal",
			oppData: game.CharacterData{
				DefaultHP: 111,
			},
			effs: []game.Effect{
				ruby.EffectCannotHeal{},
			},
			prevDmg:     88,
			isAvailable: false,
		},
		{
			name: "BelowThresholdWithCannotHeal",
			oppData: game.CharacterData{
				DefaultHP: 111,
			},
			effs: []game.Effect{
				ruby.EffectCannotHeal{},
			},
			prevDmg:     89,
			isAvailable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(ruby.CharacterRuby)
			opp := game.NewCharacter(tt.oppData)

			c.Damage(opp, tt.prevDmg, game.ColourNone)

			for _, eff := range tt.effs {
				c.AddEffect(eff)
			}

			s := c.Skills()[3]

			isAvailable := s.IsAvailable(opp, tt.gameCtx)
			assert.Equal(t, tt.isAvailable, isAvailable)
		})
	}
}

func TestSkillExecute_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		oppData game.CharacterData
		prevDmg int
		effs    []game.Effect
		gameCtx game.Context
	}{
		{
			name: "BelowThreshold",
			oppData: game.CharacterData{
				DefaultHP: 111,
			},
			prevDmg: 100,
		},
		{
			name: "BelowThresholdWithCannotHeal",
			oppData: game.CharacterData{
				DefaultHP: 111,
			},
			effs: []game.Effect{
				ruby.EffectCannotHeal{},
			},
			prevDmg: 89,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(ruby.CharacterRuby)
			opp := game.NewCharacter(tt.oppData)

			c.Damage(opp, tt.prevDmg, game.ColourNone)

			for _, eff := range tt.effs {
				c.AddEffect(eff)
			}

			s := c.Skills()[3]

			err := s.Use(opp, tt.gameCtx)
			require.NoError(t, err)

			assert.Equal(t, 0, opp.HP(), "opponent's HP")
		})
	}
}
