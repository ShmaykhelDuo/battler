package storyteller_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/storyteller"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSkillYourNumber_Use(t *testing.T) {
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
					game.ColourOrange: 0,
				},
			},
			hp: 98,
		},
		{
			name: "Opponent2",
			oppData: game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 9,
				},
				DefaultHP: 117,
				Defences: map[game.Colour]int{
					game.ColourOrange: 2,
				},
			},
			hp: 107,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(storyteller.CharacterStoryteller)
			opp := game.NewCharacter(tt.oppData)

			s := c.Skills()[0]

			err := s.Use(opp, tt.gameCtx)
			require.NoError(t, err)

			assert.Equal(t, tt.hp, opp.HP())
		})
	}
}

func TestSkillYourColour_IsAvailable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		hasOppUsedSkill bool
		oppGameCtx      game.Context
		gameCtx         game.Context
		isAvailable     bool
	}{
		{
			name:            "OppUsedSkill",
			hasOppUsedSkill: true,
			oppGameCtx: game.Context{
				TurnNum: 1,
			},
			gameCtx: game.Context{
				TurnNum: 1,
			},
			isAvailable: true,
		},
		{
			name:            "OppNotUsedSkill",
			hasOppUsedSkill: false,
			gameCtx: game.Context{
				TurnNum: 1,
			},
			isAvailable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(storyteller.CharacterStoryteller)

			data := game.CharacterData{}
			opp := game.NewCharacter(data)

			if tt.hasOppUsedSkill {
				oppS := game.NewSkill(opp, game.SkillData{
					Use: func(c, opp *game.Character, gameCtx game.Context) {},
				})
				oppS.Use(c, tt.oppGameCtx)
			}

			s := c.Skills()[1]

			isAvailable := s.IsAvailable(opp, tt.gameCtx)
			assert.Equal(t, tt.isAvailable, isAvailable)
		})
	}
}

func TestSkillYourColour_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		oppData        game.CharacterData
		oppPrevSkill   game.SkillDescription
		oppPrevGameCtx game.Context
		gameCtx        game.Context
		hp             int
		colour         game.Colour
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
			oppPrevSkill: game.SkillDescription{
				Colour: game.ColourGreen,
			},
			oppPrevGameCtx: game.Context{TurnNum: 1},
			gameCtx:        game.Context{TurnNum: 2},
			hp:             102,
			colour:         game.ColourGreen,
		},
		{
			name: "Opponent2",
			oppData: game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 9,
				},
				DefaultHP: 117,
				Defences: map[game.Colour]int{
					game.ColourOrange: 2,
				},
			},
			oppPrevSkill: game.SkillDescription{
				Colour: game.ColourOrange,
			},
			oppPrevGameCtx: game.Context{TurnNum: 1},
			gameCtx:        game.Context{TurnNum: 2},
			hp:             104,
			colour:         game.ColourOrange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(storyteller.CharacterStoryteller)
			opp := game.NewCharacter(tt.oppData)

			oppPrevSkillData := game.SkillData{
				Desc: tt.oppPrevSkill,
				Use:  func(c, opp *game.Character, gameCtx game.Context) {},
			}
			oppPrevSkill := game.NewSkill(opp, oppPrevSkillData)
			oppPrevSkill.Use(c, tt.oppPrevGameCtx)

			s := c.Skills()[1]

			err := s.Use(opp, tt.gameCtx)
			require.NoError(t, err)

			assert.Equal(t, tt.hp, opp.HP(), "HP")

			eff, ok := game.CharacterEffect[storyteller.EffectCannotUse](opp)
			require.True(t, ok, "effect")

			assert.Equal(t, tt.colour, eff.Colour(), "colour")
		})
	}
}

func TestSkillYourDream_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		oppData game.CharacterData
		maxHP   int
		dmg     int
		colour  game.Colour
		gameCtx game.Context
		hp      int
	}{
		{
			name: "Opponent1",
			oppData: game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 10,
				},
			},
			maxHP:   119,
			dmg:     42,
			colour:  game.ColourBlack,
			gameCtx: game.Context{TurnNum: 3},
			hp:      112,
		},
		{
			name: "Opponent2",
			oppData: game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 51,
				},
			},
			maxHP:   119,
			dmg:     12,
			colour:  game.ColourGreen,
			gameCtx: game.Context{TurnNum: 4},
			hp:      119,
		},
		{
			name: "Opponent3",
			oppData: game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 119,
				},
			},
			maxHP:   119,
			dmg:     20,
			colour:  game.ColourCyan,
			gameCtx: game.Context{TurnNum: 5},
			hp:      106,
		},
		{
			name: "Opponent4",
			oppData: game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 9,
				},
			},
			maxHP:   143,
			dmg:     24,
			colour:  game.ColourPink,
			gameCtx: game.Context{TurnNum: 6},
			hp:      118,
		},
		{
			name: "Opponent5",
			oppData: game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 8,
				},
			},
			maxHP:   6,
			gameCtx: game.Context{TurnNum: 9},
			hp:      6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(storyteller.CharacterStoryteller)
			opp := game.NewCharacter(tt.oppData)

			c.SetMaxHP(tt.maxHP)
			opp.Damage(c, tt.dmg, tt.colour)

			s := c.Skills()[2]

			err := s.Use(opp, tt.gameCtx)
			require.NoError(t, err)

			assert.Equal(t, tt.hp, c.HP(), "HP")
		})
	}
}

func TestSkillMyStory_Use(t *testing.T) {
	t.Parallel()

	c := game.NewCharacter(storyteller.CharacterStoryteller)

	data := game.CharacterData{}
	opp := game.NewCharacter(data)

	s := c.Skills()[3]

	gameCtx := game.Context{
		TurnNum: 8,
	}
	err := s.Use(opp, gameCtx)
	require.NoError(t, err)

	_, ok := game.CharacterEffect[storyteller.EffectControlled](opp)
	require.True(t, ok, "effect")
}
