package z89_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/z89"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSkillScarcity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		oppData   game.CharacterData
		turnState game.TurnState
		hp        int
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

			s := c.Skills()[z89.SkillScarcityIndex]

			err := s.Use(opp, tt.turnState)
			require.NoError(t, err)

			assert.Equal(t, tt.hp, opp.HP(), "HP")
			assert.Equal(t, tt.hp, opp.MaxHP(), "maximum HP")
		})
	}
}

func ultimateSlowAmount(opp *game.Character) int {
	e, ok := game.CharacterEffect[*z89.EffectUltimateSlow](opp, z89.EffectDescUltimateSlow)
	if !ok {
		return 0
	}

	return e.Amount()
}

func TestSkillIndifference_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		oppData   game.CharacterData
		effs      []game.Effect
		turnState game.TurnState
		effAmount int
	}{
		{
			name: "NoEffectNotUnlocked",
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
			turnState: game.TurnState{TurnNum: 5},
			effAmount: 1,
		},
		{
			name: "NoEffectAlreadyUnlocked",
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
			turnState: game.TurnState{TurnNum: 8},
			effAmount: 0,
		},
		{
			name: "NoEffectJustBeforeUnlocked",
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
			turnState: game.TurnState{
				TurnNum:      7,
				IsGoingFirst: true,
			},
			effAmount: 1,
		},
		{
			name: "NoEffectJustAfterUnlocked",
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
			turnState: game.TurnState{
				TurnNum:      7,
				IsGoingFirst: false,
			},
			effAmount: 0,
		},
		{
			name: "WithEffectNotUnlocked",
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
			turnState: game.TurnState{TurnNum: 5},
			effAmount: 2,
		},
		{
			name: "WithEffectUnlocked",
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
			turnState: game.TurnState{TurnNum: 9},
			effAmount: 1,
		},
		{
			name: "WithEffectJustBeforeUnlocked",
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
			turnState: game.TurnState{
				TurnNum:      8,
				IsGoingFirst: true,
			},
			effAmount: 2,
		},
		{
			name: "WithEffectJustAfterUnlocked",
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
			turnState: game.TurnState{
				TurnNum:      8,
				IsGoingFirst: false,
			},
			effAmount: 1,
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

			s := c.Skills()[z89.SkillIndifferenceIndex]

			err := s.Use(opp, tt.turnState)
			require.NoError(t, err)

			assert.Equal(t, tt.effAmount, ultimateSlowAmount(opp), "ultimate slow amount")
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
		turnState  game.TurnState
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

			s := c.Skills()[z89.SkillGreenSphereIndex]

			err := s.Use(opp, tt.turnState)
			require.NoError(t, err)

			assert.Equal(t, tt.hp, opp.HP())
		})
	}
}

func TestSkillDespondency_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		oppData   game.CharacterData
		maxHP     int
		turnState game.TurnState
		hp        int
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
			maxHP:     101,
			turnState: game.TurnState{TurnNum: 9},
			hp:        92,
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
			maxHP:     81,
			turnState: game.TurnState{TurnNum: 9},
			hp:        52,
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
			maxHP:     66,
			turnState: game.TurnState{TurnNum: 9},
			hp:        22,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(z89.CharacterZ89)
			opp := game.NewCharacter(tt.oppData)

			opp.SetMaxHP(tt.maxHP)

			s := c.Skills()[z89.SkillDespondencyIndex]

			err := s.Use(opp, tt.turnState)
			require.NoError(t, err)

			assert.Equal(t, tt.hp, opp.HP())
		})
	}
}
