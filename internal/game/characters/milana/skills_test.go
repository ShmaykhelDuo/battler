package milana_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/milana"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func stolenAmount(c *game.Character) int {
	stolenHP, ok := game.CharacterEffect[milana.EffectStolenHP](c, milana.EffectDescStolenHP)
	if !ok {
		return 0
	}

	return stolenHP.Amount()
}

func TestSkillRoyalMove_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		oppData      *game.CharacterData
		effs         []game.Effect
		turnState    game.TurnState
		hp           int
		stolenAmount int
	}{
		{
			name: "Basic",
			oppData: &game.CharacterData{
				DefaultHP: 119,
				Defences: map[game.Colour]int{
					game.ColourGreen: -2,
				},
			},
			hp:           105,
			stolenAmount: 14,
		},
		{
			name: "HasStolenHP",
			oppData: &game.CharacterData{
				DefaultHP: 119,
				Defences: map[game.Colour]int{
					game.ColourGreen: -2,
				},
			},
			effs: []game.Effect{
				milana.NewEffectStolenHP(15),
			},
			hp:           105,
			stolenAmount: 29,
		},
		{
			name: "HasMintMist",
			oppData: &game.CharacterData{
				DefaultHP: 119,
				Defences: map[game.Colour]int{
					game.ColourGreen: -2,
				},
			},
			effs: []game.Effect{
				milana.EffectMintMist{},
			},
			hp:           97,
			stolenAmount: 22,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(milana.CharacterMilana)
			opp := game.NewCharacter(tt.oppData)

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			s := c.Skills()[milana.SkillRoyalMoveIndex]

			err := s.Use(c, opp, tt.turnState)
			require.NoError(t, err)

			assert.Equal(t, tt.hp, opp.HP(), "opponent's HP")
			assert.Equal(t, tt.stolenAmount, stolenAmount(c), "stolen amount")
		})
	}
}

func TestSkillComposure_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		prevDmg      int
		effs         []game.Effect
		hp           int
		stolenAmount int
	}{
		{
			name:    "EnoughStolenHP",
			prevDmg: 34,
			effs: []game.Effect{
				milana.NewEffectStolenHP(8),
			},
			hp:           100,
			stolenAmount: 2,
		},
		{
			name:    "NotEnoughStolenHP",
			prevDmg: 34,
			effs: []game.Effect{
				milana.NewEffectStolenHP(5),
			},
			hp:           85,
			stolenAmount: 0,
		},
		{
			name:    "FullHeal",
			prevDmg: 15,
			effs: []game.Effect{
				milana.NewEffectStolenHP(8),
			},
			hp:           114,
			stolenAmount: 2,
		},
		{
			name:    "EnoughStolenHPWithMintMist",
			prevDmg: 34,
			effs: []game.Effect{
				milana.NewEffectStolenHP(11),
				milana.EffectMintMist{},
			},
			hp:           110,
			stolenAmount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(milana.CharacterMilana)
			opp := game.NewCharacter(&game.CharacterData{})

			opp.Damage(c, tt.prevDmg, game.ColourNone)

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			s := c.Skills()[milana.SkillComposureIndex]
			err := s.Use(c, opp, game.TurnState{})
			require.NoError(t, err)

			assert.Equal(t, tt.hp, c.HP(), "HP")
			assert.Equal(t, tt.stolenAmount, stolenAmount(c), "stolen amount")
		})
	}
}

func TestSkillMintMist_Use(t *testing.T) {
	t.Parallel()

	c := game.NewCharacter(milana.CharacterMilana)
	opp := game.NewCharacter(&game.CharacterData{})

	s := c.Skills()[milana.SkillMintMistIndex]
	err := s.Use(c, opp, game.TurnState{})
	require.NoError(t, err)

	eff, ok := game.CharacterEffect[milana.EffectMintMist](c, milana.EffectDescMintMist)
	require.True(t, ok, "effect")

	assert.Equal(t, 2, eff.TurnsLeft(game.TurnState{}.AddTurns(1, false)), "turns left")
}

func TestSkillPride_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		oppData *game.CharacterData
		effs    []game.Effect
		hp      int
	}{
		{
			name: "1",
			oppData: &game.CharacterData{
				DefaultHP: 119,
				Defences: map[game.Colour]int{
					game.ColourCyan: -1,
				},
			},
			effs: []game.Effect{
				milana.NewEffectStolenHP(25),
			},
			hp: 93,
		},
		{
			name: "2",
			oppData: &game.CharacterData{
				DefaultHP: 119,
				Defences: map[game.Colour]int{
					game.ColourCyan: -1,
				},
			},
			effs: []game.Effect{
				milana.NewEffectStolenHP(73),
			},
			hp: 45,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(milana.CharacterMilana)
			opp := game.NewCharacter(tt.oppData)

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			s := c.Skills()[milana.SkillPrideIndex]

			err := s.Use(c, opp, game.TurnState{TurnNum: 8})
			require.NoError(t, err)

			assert.Equal(t, tt.hp, opp.HP(), "opponent's HP")
			assert.Equal(t, 0, stolenAmount(c), "stolen amount")
		})
	}
}
