package structure_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/structure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSkillEShock_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		oppData game.CharacterData
		effs    []game.Effect
		hp      int
	}{
		{
			name: "NotBoosted",
			oppData: game.CharacterData{
				DefaultHP: 119,
				Defences: map[game.Colour]int{
					game.ColourCyan: -1,
				},
			},
			hp: 113,
		},
		{
			name: "BoostedOnce",
			oppData: game.CharacterData{
				DefaultHP: 119,
				Defences: map[game.Colour]int{
					game.ColourCyan: -1,
				},
			},
			effs: []game.Effect{
				structure.NewEffectIBoost(5),
			},
			hp: 108,
		},
		{
			name: "BoostedThrice",
			oppData: game.CharacterData{
				DefaultHP: 119,
				Defences: map[game.Colour]int{
					game.ColourCyan: -1,
				},
			},
			effs: []game.Effect{
				structure.NewEffectIBoost(15),
			},
			hp: 98,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(structure.CharacterStructure)
			opp := game.NewCharacter(tt.oppData)

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			s := c.Skills()[0]
			err := s.Use(opp, game.Context{})
			require.NoError(t, err)

			assert.Equal(t, tt.hp, opp.HP())
		})
	}
}

func TestSkillIBoost_IsAvailable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		effs        []game.Effect
		isAvailable bool
	}{
		{
			name:        "Basic",
			isAvailable: true,
		},
		{
			name: "BoostedTwice",
			effs: []game.Effect{
				structure.NewEffectIBoost(10),
			},
			isAvailable: true,
		},
		{
			name: "BoostedThrice",
			effs: []game.Effect{
				structure.NewEffectIBoost(15),
			},
			isAvailable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(structure.CharacterStructure)
			opp := game.NewCharacter(game.CharacterData{})

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			s := c.Skills()[1]
			isAvailable := s.IsAvailable(opp, game.Context{})
			assert.Equal(t, tt.isAvailable, isAvailable)
		})
	}
}

func TestSkillIBoost_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		effs   []game.Effect
		amount int
	}{
		{
			name:   "Basic",
			amount: 5,
		},
		{
			name: "BoostedOnce",
			effs: []game.Effect{
				structure.NewEffectIBoost(5),
			},
			amount: 10,
		},
		{
			name: "BoostedTwice",
			effs: []game.Effect{
				structure.NewEffectIBoost(10),
			},
			amount: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(structure.CharacterStructure)
			opp := game.NewCharacter(game.CharacterData{})

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			s := c.Skills()[1]
			err := s.Use(opp, game.Context{})
			require.NoError(t, err)

			boost, ok := game.CharacterEffect[*structure.EffectIBoost](c)
			require.True(t, ok, "effect")

			assert.Equal(t, tt.amount, boost.Amount())
		})
	}
}

func TestSkillSLayers_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		effs      []game.Effect
		threshold int
	}{
		{
			name:      "Basic",
			threshold: 5,
		},
		{
			name: "BoostedOnce",
			effs: []game.Effect{
				structure.NewEffectIBoost(5),
			},
			threshold: 10,
		},
		{
			name: "BoostedTwice",
			effs: []game.Effect{
				structure.NewEffectIBoost(10),
			},
			threshold: 15,
		},
		{
			name: "BoostedThrice",
			effs: []game.Effect{
				structure.NewEffectIBoost(15),
			},
			threshold: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(structure.CharacterStructure)
			opp := game.NewCharacter(game.CharacterData{})

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			s := c.Skills()[2]
			err := s.Use(opp, game.Context{})
			require.NoError(t, err)

			layers, ok := game.CharacterEffect[structure.EffectSLayers](c)
			require.True(t, ok, "effect")

			assert.Equal(t, tt.threshold, layers.Threshold(), "threshold")
		})
	}
}

func TestSkillLastChance_Use(t *testing.T) {
	t.Parallel()

	c := game.NewCharacter(structure.CharacterStructure)
	opp := game.NewCharacter(game.CharacterData{})

	s := c.Skills()[3]
	err := s.Use(opp, game.Context{TurnNum: 7})
	require.NoError(t, err)

	_, ok := game.CharacterEffect[structure.EffectLastChance](c)
	require.True(t, ok, "effect")
}
