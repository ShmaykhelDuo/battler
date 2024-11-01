package storyteller_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/storyteller"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSkillYourNumber_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		oppData   *game.CharacterData
		turnState game.TurnState
		hp        int
	}{
		{
			name: "Opponent1",
			oppData: &game.CharacterData{
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
			oppData: &game.CharacterData{
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

			s := c.Skills()[storyteller.SkillYourNumberIndex]

			err := s.Use(c, opp, tt.turnState)
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
		oppturnState    game.TurnState
		turnState       game.TurnState
		isAvailable     bool
	}{
		{
			name:            "OppUsedSkill",
			hasOppUsedSkill: true,
			oppturnState: game.TurnState{
				TurnNum: 1,
			},
			turnState: game.TurnState{
				TurnNum: 1,
			},
			isAvailable: true,
		},
		{
			name:            "OppNotUsedSkill",
			hasOppUsedSkill: false,
			turnState: game.TurnState{
				TurnNum: 1,
			},
			isAvailable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(storyteller.CharacterStoryteller)

			data := &game.CharacterData{}
			opp := game.NewCharacter(data)

			if tt.hasOppUsedSkill {
				oppS := game.NewSkill(&game.SkillData{
					Use: func(c, opp *game.Character, turnState game.TurnState) {},
				})
				oppS.Use(opp, c, tt.oppturnState)
			}

			s := c.Skills()[storyteller.SkillYourColourIndex]

			isAvailable := s.IsAvailable(c, opp, tt.turnState)
			assert.Equal(t, tt.isAvailable, isAvailable)
		})
	}
}

func TestSkillYourColour_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		oppData          *game.CharacterData
		oppPrevSkill     game.SkillDescription
		oppPrevturnState game.TurnState
		turnState        game.TurnState
		hp               int
		colour           game.Colour
	}{
		{
			name: "Opponent1",
			oppData: &game.CharacterData{
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
			oppPrevturnState: game.TurnState{TurnNum: 1},
			turnState:        game.TurnState{TurnNum: 2},
			hp:               102,
			colour:           game.ColourGreen,
		},
		{
			name: "Opponent2",
			oppData: &game.CharacterData{
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
			oppPrevturnState: game.TurnState{TurnNum: 1},
			turnState:        game.TurnState{TurnNum: 2},
			hp:               104,
			colour:           game.ColourOrange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(storyteller.CharacterStoryteller)
			opp := game.NewCharacter(tt.oppData)

			oppPrevSkillData := &game.SkillData{
				Desc: tt.oppPrevSkill,
				Use:  func(c, opp *game.Character, turnState game.TurnState) {},
			}
			oppPrevSkill := game.NewSkill(oppPrevSkillData)
			oppPrevSkill.Use(opp, c, tt.oppPrevturnState)

			s := c.Skills()[storyteller.SkillYourColourIndex]

			err := s.Use(c, opp, tt.turnState)
			require.NoError(t, err)

			assert.Equal(t, tt.hp, opp.HP(), "HP")

			eff, ok := game.CharacterEffect[storyteller.EffectCannotUse](opp, storyteller.EffectDescCannotUse)
			require.True(t, ok, "effect")

			assert.Equal(t, tt.colour, eff.Colour(), "colour")
			assert.Equal(t, 1, eff.TurnsLeft(tt.turnState.AddTurns(0, true)), "turns left opponent")
		})
	}
}

func TestSkillYourDream_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		oppData   *game.CharacterData
		maxHP     int
		dmg       int
		colour    game.Colour
		turnState game.TurnState
		hp        int
	}{
		{
			name: "Opponent1",
			oppData: &game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 10,
				},
			},
			maxHP:     119,
			dmg:       42,
			colour:    game.ColourBlack,
			turnState: game.TurnState{TurnNum: 3},
			hp:        112,
		},
		{
			name: "Opponent2",
			oppData: &game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 51,
				},
			},
			maxHP:     119,
			dmg:       12,
			colour:    game.ColourGreen,
			turnState: game.TurnState{TurnNum: 4},
			hp:        119,
		},
		{
			name: "Opponent3",
			oppData: &game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 119,
				},
			},
			maxHP:     119,
			dmg:       20,
			colour:    game.ColourCyan,
			turnState: game.TurnState{TurnNum: 5},
			hp:        106,
		},
		{
			name: "Opponent4",
			oppData: &game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 9,
				},
			},
			maxHP:     143,
			dmg:       24,
			colour:    game.ColourPink,
			turnState: game.TurnState{TurnNum: 6},
			hp:        118,
		},
		{
			name: "Opponent5",
			oppData: &game.CharacterData{
				Desc: game.CharacterDescription{
					Number: 8,
				},
			},
			maxHP:     6,
			turnState: game.TurnState{TurnNum: 9},
			hp:        6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(storyteller.CharacterStoryteller)
			opp := game.NewCharacter(tt.oppData)

			c.SetMaxHP(tt.maxHP)
			opp.Damage(c, tt.dmg, tt.colour)

			s := c.Skills()[storyteller.SkillYourDreamIndex]

			err := s.Use(c, opp, tt.turnState)
			require.NoError(t, err)

			assert.Equal(t, tt.hp, c.HP(), "HP")
		})
	}
}

func TestSkillMyStory_Use(t *testing.T) {
	t.Parallel()

	c := game.NewCharacter(storyteller.CharacterStoryteller)

	data := &game.CharacterData{}
	opp := game.NewCharacter(data)

	s := c.Skills()[storyteller.SkillMyStoryIndex]

	turnState := game.TurnState{
		TurnNum: 8,
	}
	err := s.Use(c, opp, turnState)
	require.NoError(t, err)

	eff, ok := game.CharacterEffect[storyteller.EffectControlled](opp, storyteller.EffectDescControlled)
	require.True(t, ok, "effect")

	assert.Equal(t, 1, eff.TurnsLeft(turnState.AddTurns(0, true)), "turns left opponent")
}
