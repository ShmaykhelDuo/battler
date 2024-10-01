package storyteller_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/storyteller"
	"github.com/stretchr/testify/assert"
)

func TestEffectCannotUse_IsSkillAvailable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		colour      game.Colour
		skill       game.SkillData
		isAvailable bool
	}{
		{
			name:   "ColoursMatch",
			colour: game.ColourViolet,
			skill: game.SkillData{
				Desc: game.SkillDescription{
					Colour: game.ColourViolet,
				},
			},
			isAvailable: false,
		},
		{
			name:   "ColoursNotMatch",
			colour: game.ColourViolet,
			skill: game.SkillData{
				Desc: game.SkillDescription{
					Colour: game.ColourGreen,
				},
			},
			isAvailable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := storyteller.NewEffectCannotUse(tt.colour)

			data := game.CharacterData{}
			c := game.NewCharacter(data)
			s := game.NewSkill(c, tt.skill)

			isAvailable := e.IsSkillAvailable(s)
			assert.Equal(t, tt.isAvailable, isAvailable)
		})
	}
}
