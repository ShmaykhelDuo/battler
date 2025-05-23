package structure_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/structure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSkillEShock_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		oppData *game.CharacterData
		effs    []game.Effect
		hp      int
	}{
		{
			name: "NotBoosted",
			oppData: &game.CharacterData{
				DefaultHP: 119,
				Defences: map[game.Colour]int{
					game.ColourCyan: -1,
				},
			},
			hp: 113,
		},
		{
			name: "BoostedOnce",
			oppData: &game.CharacterData{
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
			oppData: &game.CharacterData{
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

			s := c.Skills()[structure.SkillEShockIndex]
			err := s.Use(c, opp, game.TurnState{})
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
			opp := game.NewCharacter(&game.CharacterData{})

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			s := c.Skills()[structure.SkillIBoostIndex]
			isAvailable := s.IsAvailable(c, opp, game.TurnState{})
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
			opp := game.NewCharacter(&game.CharacterData{})

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			s := c.Skills()[structure.SkillIBoostIndex]
			err := s.Use(c, opp, game.TurnState{})
			require.NoError(t, err)

			boost, ok := game.CharacterEffect[*structure.EffectIBoost](c, structure.EffectDescIBoost)
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
			opp := game.NewCharacter(&game.CharacterData{})

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			s := c.Skills()[structure.SkillSLayersIndex]
			err := s.Use(c, opp, game.TurnState{})
			require.NoError(t, err)

			layers, ok := game.CharacterEffect[structure.EffectSLayers](c, structure.EffectDescSLayers)
			require.True(t, ok, "effect")

			assert.Equal(t, tt.threshold, layers.Threshold(), "threshold")
			assert.Equal(t, 1, layers.TurnsLeft(game.TurnState{}.AddTurns(0, true)), "turns left")
		})
	}
}

func TestSkillLastChance_Use(t *testing.T) {
	t.Parallel()

	c := game.NewCharacter(structure.CharacterStructure)
	opp := game.NewCharacter(&game.CharacterData{})

	turnState := game.TurnState{TurnNum: 7}

	s := c.Skills()[structure.SkillLastChanceIndex]
	err := s.Use(c, opp, turnState)
	require.NoError(t, err)

	eff, ok := game.CharacterEffect[structure.EffectLastChance](c, structure.EffectDescLastChance)
	require.True(t, ok, "effect")

	assert.Equal(t, 1, eff.TurnsLeft(turnState.AddTurns(0, true)), "turns left")
}
