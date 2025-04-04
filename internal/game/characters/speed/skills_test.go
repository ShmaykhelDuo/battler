package speed_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/speed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func greenTokensNumber(c *game.Character) int {
	tokens, ok := game.CharacterEffect[speed.EffectGreenTokens](c, speed.EffectDescGreenTokens)
	if !ok {
		return 0
	}

	return tokens.Amount()
}

func blackTokensNumber(c *game.Character) int {
	tokens, ok := game.CharacterEffect[speed.EffectBlackTokens](c, speed.EffectDescBlackTokens)
	if !ok {
		return 0
	}

	return tokens.Amount()
}

func runSkill(t *testing.T, skillIndex int, effs []game.Effect, turnState game.TurnState) (c, opp *game.Character) {
	t.Helper()

	c = game.NewCharacter(speed.CharacterSpeed)

	data := &game.CharacterData{}
	opp = game.NewCharacter(data)

	for _, e := range effs {
		c.AddEffect(e)
	}

	s := c.Skills()[skillIndex]

	err := s.Use(c, opp, turnState)
	require.NoError(t, err)

	return
}

func testGainGreenToken(t *testing.T, skillIndex int) {
	t.Run("GainGreenToken", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name      string
			effs      []game.Effect
			turnState game.TurnState
			number    int
		}{
			{
				name:   "Basic",
				number: 1,
			},
			{
				name: "LowTokens",
				effs: []game.Effect{
					speed.NewEffectGreenTokens(3),
				},
				number: 4,
			},
			{
				name: "FullTokens",
				effs: []game.Effect{
					speed.NewEffectGreenTokens(5),
				},
				number: 5,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				c, _ := runSkill(t, skillIndex, tt.effs, tt.turnState)

				assert.Equal(t, tt.number, greenTokensNumber(c), "tokens")
			})
		}
	})
}

func TestSkillRun_Use(t *testing.T) {
	t.Parallel()

	t.Run("ReduceDamage", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name                  string
			effs                  []game.Effect
			turnState             game.TurnState
			damageReductionAmount int
		}{
			{
				name:                  "NoneApplied",
				damageReductionAmount: 5,
			},
			{
				name: "SomeApplied",
				effs: []game.Effect{
					speed.NewEffectDamageReduced(5),
				},
				damageReductionAmount: 10,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				c, _ := runSkill(t, speed.SkillRunIndex, tt.effs, tt.turnState)

				red, ok := game.CharacterEffect[*speed.EffectDamageReduced](c, speed.EffectDescDamageReduced)
				require.True(t, ok, "effect")

				assert.Equal(t, tt.damageReductionAmount, red.Amount(), "reduction amount")
			})
		}
	})

	testGainGreenToken(t, speed.SkillRunIndex)
}

func TestSkillWeaken_Use(t *testing.T) {
	t.Parallel()

	t.Run("ReduceDefence", func(t *testing.T) {
		t.Parallel()

		_, opp := runSkill(t, speed.SkillWeakenIndex, nil, game.TurnState{})

		_, ok := game.CharacterEffect[speed.EffectDefenceReduced](opp, speed.EffectDescDefenceReduced)
		require.True(t, ok, "effect")
	})

	t.Run("GainBlackToken", func(t *testing.T) {
		tests := []struct {
			name      string
			effs      []game.Effect
			turnState game.TurnState
			number    int
		}{
			{
				name:   "Basic",
				number: 1,
			},
			{
				name: "LowTokens",
				effs: []game.Effect{
					speed.NewEffectBlackTokens(3),
				},
				number: 4,
			},
			{
				name: "FullTokens",
				effs: []game.Effect{
					speed.NewEffectBlackTokens(5),
				},
				number: 5,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				c, _ := runSkill(t, speed.SkillWeakenIndex, tt.effs, tt.turnState)

				assert.Equal(t, tt.number, blackTokensNumber(c), "tokens")
			})
		}
	})
}

func TestSkillSpeed_Use(t *testing.T) {
	t.Parallel()

	t.Run("SpeedUp", func(t *testing.T) {
		t.Parallel()

		c, _ := runSkill(t, speed.SkillSpeedIndex, nil, game.TurnState{})

		eff, ok := game.CharacterEffect[speed.EffectSpedUp](c, speed.EffectDescSpedUp)
		require.True(t, ok, "effect")

		assert.Equal(t, 1, eff.TurnsLeft(game.TurnState{}.AddTurns(1, false)), "turns left")
	})

	testGainGreenToken(t, speed.SkillSpeedIndex)
}

func TestSkillStab_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		oppData   *game.CharacterData
		effs      []game.Effect
		turnState game.TurnState
		hp        int
	}{
		{
			name: "Opponent1",
			oppData: &game.CharacterData{
				DefaultHP: 119,
				Defences: map[game.Colour]int{
					game.ColourGreen: -2,
					game.ColourBlack: -2,
				},
			},
			effs: []game.Effect{
				speed.NewEffectGreenTokens(2),
				speed.NewEffectBlackTokens(4),
			},
			hp: 79,
		},
		{
			name: "Opponent2",
			oppData: &game.CharacterData{
				DefaultHP: 114,
				Defences: map[game.Colour]int{
					game.ColourGreen: 1,
					game.ColourBlack: -2,
				},
			},
			effs: []game.Effect{
				speed.NewEffectGreenTokens(3),
				speed.NewEffectBlackTokens(2),
			},
			hp: 83,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(speed.CharacterSpeed)
			opp := game.NewCharacter(tt.oppData)

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			s := c.Skills()[speed.SkillStabIndex]

			err := s.Use(c, opp, tt.turnState)
			require.NoError(t, err)

			assert.Equal(t, tt.hp, opp.HP())
		})
	}
}
