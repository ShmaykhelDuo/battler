package structure

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

var CharacterStructure = game.CharacterData{
	Desc: game.CharacterDescription{
		Name:   "Structure",
		Number: 119,
	},
	DefaultHP: 119,
	Defences: map[game.Colour]int{
		game.ColourRed:    -2,
		game.ColourOrange: -1,
		game.ColourYellow: -2,
		game.ColourGreen:  0,
		game.ColourCyan:   2,
		game.ColourBlue:   -2,
		game.ColourViolet: 1,
		game.ColourPink:   0,
		game.ColourGray:   2,
		game.ColourBrown:  0,
		game.ColourBlack:  -3,
		game.ColourWhite:  1,
	},
	SkillData: [4]game.SkillData{
		SkillEShock,
		SkillIBoost,
		SkillSLayers,
		SkillLastChance,
	},
}
