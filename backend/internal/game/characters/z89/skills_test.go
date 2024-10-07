package z89_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/z89"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSkillScarcity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		oppData game.CharacterData
		gameCtx game.Context
		hp      int
	}{
		{
			name: "Opponent1",
			oppData: game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 33,
				},
				DefaultHP: 113,
				Defences: map[game.Colour]int{
					game.ColourBlack: 2,
				},
			},
			hp: 103,
		},
		{
			name: "Opponent2",
			oppData: game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 9,
				},
				DefaultHP: 117,
				Defences: map[game.Colour]int{
					game.ColourBlack: 0,
				},
			},
			hp: 105,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(z89.CharacterZ89)
			opp := game.NewCharacter(tt.oppData)

			s := c.Skills()[0]

			err := s.Use(opp, tt.gameCtx)
			require.NoError(t, err)

			assert.Equal(t, tt.hp, opp.HP(), "HP")
			assert.Equal(t, tt.hp, opp.MaxHP(), "maximum HP")
		})
	}
}

func TestSkillIndifference_IsAvailable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		oppData     game.CharacterData
		effs        []game.Effect
		gameCtx     game.Context
		isAvailable bool
	}{
		{
			name: "Basic",
			oppData: game.CharacterData{
				SkillData: [4]game.SkillData{
					3: {
						Desc: game.SkillDescription{
							IsUltimate: true,
						},
						UnlockTurn: 7,
					},
				},
			},
			gameCtx:     game.Context{TurnNum: 5},
			isAvailable: true,
		},
		{
			name: "AlreadyUnlocked1",
			oppData: game.CharacterData{
				SkillData: [4]game.SkillData{
					3: {
						Desc: game.SkillDescription{
							IsUltimate: true,
						},
						UnlockTurn: 7,
					},
				},
			},
			gameCtx: game.Context{
				TurnNum:      8,
				IsGoingFirst: true,
			},
			isAvailable: false,
		},
		{
			name: "AlreadyUnlocked2",
			oppData: game.CharacterData{
				SkillData: [4]game.SkillData{
					3: {
						Desc: game.SkillDescription{
							IsUltimate: true,
						},
						UnlockTurn: 7,
					},
				},
			},
			gameCtx: game.Context{
				TurnNum:      8,
				IsGoingFirst: false,
			},
			isAvailable: false,
		},
		{
			name: "JustBeforeUnlocked",
			oppData: game.CharacterData{
				SkillData: [4]game.SkillData{
					3: {
						Desc: game.SkillDescription{
							IsUltimate: true,
						},
						UnlockTurn: 7,
					},
				},
			},
			gameCtx: game.Context{
				TurnNum:      7,
				IsGoingFirst: true,
			},
			isAvailable: true,
		},
		{
			name: "JustAfterUnlocked",
			oppData: game.CharacterData{
				SkillData: [4]game.SkillData{
					3: {
						Desc: game.SkillDescription{
							IsUltimate: true,
						},
						UnlockTurn: 7,
					},
				},
			},
			gameCtx: game.Context{
				TurnNum:      7,
				IsGoingFirst: false,
			},
			isAvailable: false,
		},
		{
			name: "JustBeforeUnlockedWithEffect",
			oppData: game.CharacterData{
				SkillData: [4]game.SkillData{
					3: {
						Desc: game.SkillDescription{
							IsUltimate: true,
						},
						UnlockTurn: 7,
					},
				},
			},
			effs: []game.Effect{
				z89.NewEffectUltimateSlow(),
			},
			gameCtx: game.Context{
				TurnNum:      8,
				IsGoingFirst: true,
			},
			isAvailable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(z89.CharacterZ89)
			opp := game.NewCharacter(tt.oppData)

			for _, e := range tt.effs {
				opp.AddEffect(e)
			}

			s := c.Skills()[1]

			isAvailable := s.IsAvailable(opp, tt.gameCtx)
			assert.Equal(t, tt.isAvailable, isAvailable)
		})
	}
}

func TestSkillIndifference_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		oppData   game.CharacterData
		effs      []game.Effect
		effAmount int
	}{
		{
			name: "NoEffect",
			oppData: game.CharacterData{
				SkillData: [4]game.SkillData{
					3: {
						Desc: game.SkillDescription{
							IsUltimate: true,
						},
						UnlockTurn: 7,
					},
				},
			},
			effAmount: 1,
		},
		{
			name: "Effect",
			oppData: game.CharacterData{
				SkillData: [4]game.SkillData{
					3: {
						Desc: game.SkillDescription{
							IsUltimate: true,
						},
						UnlockTurn: 7,
					},
				},
			},
			effs: []game.Effect{
				z89.NewEffectUltimateSlow(),
			},
			effAmount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(z89.CharacterZ89)
			opp := game.NewCharacter(tt.oppData)

			for _, e := range tt.effs {
				opp.AddEffect(e)
			}

			s := c.Skills()[1]

			gameCtx := game.Context{TurnNum: 5}
			err := s.Use(opp, gameCtx)
			require.NoError(t, err)

			e, ok := game.CharacterEffect[*z89.EffectUltimateSlow](opp)
			require.True(t, ok, "effect")

			assert.Equal(t, tt.effAmount, e.Amount(), "amount")
		})
	}
}

func TestSkillGreenSphere_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		oppData    game.CharacterData
		prevDmg    int
		prevColour game.Colour
		gameCtx    game.Context
		hp         int
	}{
		{
			name: "Opponent1",
			oppData: game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 33,
				},
				DefaultHP: 113,
				Defences: map[game.Colour]int{
					game.ColourGreen: 4,
				},
			},
			hp: 102,
		},
		{
			name: "Opponent2",
			oppData: game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 9,
				},
				DefaultHP: 117,
				Defences: map[game.Colour]int{
					game.ColourGreen: 0,
					game.ColourBlack: 0,
				},
			},
			prevDmg:    12,
			prevColour: game.ColourBlack,
			hp:         102,
		},
		{
			name: "Opponent3",
			oppData: game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 51,
				},
				DefaultHP: 114,
				Defences: map[game.Colour]int{
					game.ColourGreen: 1,
					game.ColourBlack: -2,
				},
			},
			prevDmg:    24,
			prevColour: game.ColourBlack,
			hp:         88,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(z89.CharacterZ89)
			opp := game.NewCharacter(tt.oppData)

			c.Damage(opp, tt.prevDmg, tt.prevColour)

			s := c.Skills()[2]

			err := s.Use(opp, tt.gameCtx)
			require.NoError(t, err)

			assert.Equal(t, tt.hp, opp.HP())
		})
	}
}

func TestSkillDespondency_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		oppData game.CharacterData
		maxHP   int
		gameCtx game.Context
		hp      int
	}{
		{
			name: "Opponent1",
			oppData: game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 33,
				},
				DefaultHP: 113,
				Defences: map[game.Colour]int{
					game.ColourBlue: 0,
				},
			},
			maxHP:   101,
			gameCtx: game.Context{TurnNum: 9},
			hp:      92,
		},
		{
			name: "Opponent2",
			oppData: game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 9,
				},
				DefaultHP: 117,
				Defences: map[game.Colour]int{
					game.ColourBlue: 0,
				},
			},
			maxHP:   81,
			gameCtx: game.Context{TurnNum: 9},
			hp:      52,
		},
		{
			name: "Opponent3",
			oppData: game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 51,
				},
				DefaultHP: 114,
				Defences: map[game.Colour]int{
					game.ColourGreen: 1,
					game.ColourBlack: -2,
				},
			},
			maxHP:   66,
			gameCtx: game.Context{TurnNum: 9},
			hp:      22,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(z89.CharacterZ89)
			opp := game.NewCharacter(tt.oppData)

			opp.SetMaxHP(tt.maxHP)

			s := c.Skills()[3]

			err := s.Use(opp, tt.gameCtx)
			require.NoError(t, err)

			assert.Equal(t, tt.hp, opp.HP())
		})
	}
}
