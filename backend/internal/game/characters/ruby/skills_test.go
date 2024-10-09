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

	eff, ok := game.CharacterEffect[ruby.EffectDoubleDamage](c)
	require.True(t, ok, "effect")

	assert.Equal(t, 2, eff.TurnsLeft(gameCtx.AddTurns(1, false)), "turns left next turn")
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

func assertEffectCannotHeal(t *testing.T, c *game.Character, gameCtx game.Context, isOpp bool, name string) {
	t.Helper()

	eff, ok := game.CharacterEffect[ruby.EffectCannotHeal](c)
	if !assert.True(t, ok, "%s's effect", name) {
		return
	}

	checkCtx := gameCtx.AddTurns(1, false)
	if isOpp {
		checkCtx = gameCtx.AddTurns(0, true)
	}

	assert.Equal(t, 1, eff.TurnsLeft(checkCtx), "%s's effect turns left", name)
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

	assertEffectCannotHeal(t, c, gameCtx, false, "character")
	assertEffectCannotHeal(t, opp, gameCtx, true, "opponent")
}

func TestSkillExecute_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		oppData game.CharacterData
		prevDmg int
		effs    []game.Effect
		gameCtx game.Context
		hp      int
	}{
		{
			name: "AboveThreshold",
			oppData: game.CharacterData{
				DefaultHP: 111,
			},
			prevDmg: 99,
			hp:      12,
		},
		{
			name: "BelowThreshold",
			oppData: game.CharacterData{
				DefaultHP: 111,
			},
			prevDmg: 100,
			hp:      0,
		},
		{
			name: "AboveThresholdWithCannotHeal",
			oppData: game.CharacterData{
				DefaultHP: 111,
			},
			effs: []game.Effect{
				ruby.EffectCannotHeal{},
			},
			prevDmg: 88,
			hp:      23,
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
			hp:      0,
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

			assert.Equal(t, tt.hp, opp.HP(), "opponent's HP")
		})
	}
}
