package euphoria_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/euphoria"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSkillAmpleness_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		oppData              *game.CharacterData
		initMaxHP            int
		initOppMaxHP         int
		effs                 []game.Effect
		turnState            game.TurnState
		maxHP                int
		oppMaxHP             int
		euphoricSourceAmount int
	}{
		{
			name: "Opponent1",
			oppData: &game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 33,
				},
				DefaultHP: 113,
			},
			initMaxHP:            117,
			initOppMaxHP:         113,
			maxHP:                129,
			oppMaxHP:             125,
			euphoricSourceAmount: 12,
		},
		{
			name: "Opponent2",
			oppData: &game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 51,
				},
				DefaultHP: 114,
			},
			effs: []game.Effect{
				euphoria.NewEffectEuphoricSource(12),
			},
			initMaxHP:            129,
			initOppMaxHP:         104,
			maxHP:                141,
			oppMaxHP:             116,
			euphoricSourceAmount: 24,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(euphoria.CharacterEuphoria)
			opp := game.NewCharacter(tt.oppData)

			c.SetMaxHP(tt.initMaxHP)
			opp.SetMaxHP(tt.initOppMaxHP)

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			s := c.Skills()[euphoria.SkillAmplenessIndex]

			err := s.Use(c, opp, tt.turnState)
			require.NoError(t, err)

			assert.Equal(t, tt.maxHP, c.MaxHP(), "maximum HP")
			assert.Equal(t, tt.oppMaxHP, opp.MaxHP(), "opponent's maximum HP")

			assert.Equal(t, tt.euphoricSourceAmount, euphoricSourceAmount(c), "euphoric source amount")
		})
	}
}

func ultimateEarlyAmount(opp *game.Character) int {
	eff, ok := game.CharacterEffect[*euphoria.EffectUltimateEarly](opp, euphoria.EffectDescUltimateEarly)
	if !ok {
		return 0
	}

	return eff.Amount()
}

func TestSkillExuberance_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		oppData              *game.CharacterData
		effs                 []game.Effect
		oppEffs              []game.Effect
		turnState            game.TurnState
		hp                   int
		maxHP                int
		oppMaxHP             int
		euphoricSourceAmount int
		ultimateEarlyAmount  int
	}{
		{
			name: "UltimateNotUnlockedNoEuphoricSource",
			oppData: &game.CharacterData{
				DefaultHP: 100,
				SkillData: [4]*game.SkillData{
					3: {
						Desc: game.SkillDescription{
							IsUltimate: true,
						},
						UnlockTurn: 7,
					},
				},
			},
			turnState:            game.TurnState{TurnNum: 4},
			hp:                   127,
			maxHP:                127,
			oppMaxHP:             110,
			euphoricSourceAmount: 10,
			ultimateEarlyAmount:  1,
		},
		{
			name: "UltimateNotUnlockedHasEuphoricSource",
			oppData: &game.CharacterData{
				DefaultHP: 100,
				SkillData: [4]*game.SkillData{
					3: {
						Desc: game.SkillDescription{
							IsUltimate: true,
						},
						UnlockTurn: 7,
					},
				},
			},
			effs: []game.Effect{
				euphoria.NewEffectEuphoricSource(12),
			},
			turnState:            game.TurnState{TurnNum: 4},
			hp:                   127,
			maxHP:                127,
			oppMaxHP:             110,
			euphoricSourceAmount: 22,
			ultimateEarlyAmount:  1,
		},
		{
			name: "UltimateNotUnlockedNoEuphoricSourceEarly",
			oppData: &game.CharacterData{
				DefaultHP: 100,
				SkillData: [4]*game.SkillData{
					3: {
						Desc: game.SkillDescription{
							IsUltimate: true,
						},
						UnlockTurn: 7,
					},
				},
			},
			oppEffs: []game.Effect{
				euphoria.NewEffectUltimateEarly(),
			},
			hp:                   127,
			maxHP:                127,
			oppMaxHP:             110,
			turnState:            game.TurnState{TurnNum: 4},
			euphoricSourceAmount: 10,
			ultimateEarlyAmount:  2,
		},
		{
			name: "UltimateUnlockedNoEuphoricSource",
			oppData: &game.CharacterData{
				DefaultHP: 100,
				SkillData: [4]*game.SkillData{
					3: {
						Desc: game.SkillDescription{
							IsUltimate: true,
						},
						UnlockTurn: 7,
					},
				},
			},
			hp:                   137,
			maxHP:                137,
			oppMaxHP:             120,
			turnState:            game.TurnState{TurnNum: 8},
			euphoricSourceAmount: 20,
		},
		{
			name: "UltimateJustBeforeUnlockedNoEuphoricSource",
			oppData: &game.CharacterData{
				DefaultHP: 100,
				SkillData: [4]*game.SkillData{
					3: {
						Desc: game.SkillDescription{
							IsUltimate: true,
						},
						UnlockTurn: 7,
					},
				},
			},
			turnState: game.TurnState{
				TurnNum:      7,
				IsGoingFirst: true,
			},
			hp:                   127,
			maxHP:                127,
			oppMaxHP:             110,
			euphoricSourceAmount: 10,
			ultimateEarlyAmount:  1,
		},
		{
			name: "UltimateJustAfterUnlockedNoEuphoricSource",
			oppData: &game.CharacterData{
				DefaultHP: 100,
				SkillData: [4]*game.SkillData{
					3: {
						Desc: game.SkillDescription{
							IsUltimate: true,
						},
						UnlockTurn: 7,
					},
				},
			},
			turnState: game.TurnState{
				TurnNum:      7,
				IsGoingFirst: false,
			},
			hp:                   137,
			maxHP:                137,
			oppMaxHP:             120,
			euphoricSourceAmount: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(euphoria.CharacterEuphoria)
			opp := game.NewCharacter(tt.oppData)

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			for _, e := range tt.oppEffs {
				opp.AddEffect(e)
			}

			s := c.Skills()[euphoria.SkillExuberanceIndex]

			err := s.Use(c, opp, tt.turnState)
			require.NoError(t, err)

			assert.Equal(t, tt.hp, c.HP(), "HP")
			assert.Equal(t, tt.maxHP, c.MaxHP(), "maxiumum HP")
			assert.Equal(t, tt.oppMaxHP, opp.MaxHP(), "opponent's maximum HP")
			assert.Equal(t, tt.euphoricSourceAmount, euphoricSourceAmount(c), "euphoric source amount")
			assert.Equal(t, tt.ultimateEarlyAmount, ultimateEarlyAmount(opp), "ultimate early amount")
		})
	}
}

func TestSkillPinkSphere_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		oppData   *game.CharacterData
		turnState game.TurnState
		oppHP     int
		maxHP     int
		oppMaxHP  int
	}{
		{
			name: "Opponent1",
			oppData: &game.CharacterData{
				DefaultHP: 113,
			},
			oppHP:    101,
			maxHP:    129,
			oppMaxHP: 125,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(euphoria.CharacterEuphoria)
			opp := game.NewCharacter(tt.oppData)

			s := c.Skills()[euphoria.SkillPinkSphereIndex]

			err := s.Use(c, opp, tt.turnState)
			require.NoError(t, err)

			assert.Equal(t, tt.oppHP, opp.HP(), "opponent's HP")
			assert.Equal(t, tt.maxHP, c.MaxHP(), "maximum HP")
			assert.Equal(t, tt.oppMaxHP, opp.MaxHP(), "opponent's maximum HP")
		})
	}
}

func TestSkillEuphoria_Use(t *testing.T) {
	t.Parallel()

	c := game.NewCharacter(euphoria.CharacterEuphoria)

	data := &game.CharacterData{}
	opp := game.NewCharacter(data)

	s := c.Skills()[euphoria.SkillEuphoriaIndex]

	turnState := game.TurnState{TurnNum: 4}
	err := s.Use(c, opp, turnState)
	require.NoError(t, err)

	_, ok := game.CharacterEffect[euphoria.EffectEuphoricHeal](c, euphoria.EffectDescEuphoricHeal)
	require.True(t, ok, "ultimate heal effect")
}
