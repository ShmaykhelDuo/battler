package z89

import "github.com/ShmaykhelDuo/battler/internal/game"

var CharacterZ89 = &game.CharacterData{
	Desc: game.CharacterDescription{
		Name:   "Z89",
		Number: 8,
	},
	DefaultHP: 111,
	Defences: map[game.Colour]int{
		game.ColourRed:    -1,
		game.ColourOrange: 0,
		game.ColourYellow: -1,
		game.ColourGreen:  2,
		game.ColourCyan:   3,
		game.ColourBlue:   2,
		game.ColourViolet: -1,
		game.ColourPink:   -1,
		game.ColourGray:   0,
		game.ColourBrown:  0,
		game.ColourBlack:  2,
		game.ColourWhite:  -1,
	},
	SkillData: [4]*game.SkillData{
		SkillScarcity,
		SkillIndifference,
		SkillGreenSphere,
		SkillDespondency,
	},
}
