package ruby

import "github.com/ShmaykhelDuo/battler/internal/game"

var CharacterRuby = &game.CharacterData{
	Desc: game.CharacterDescription{
		Name:   "Ruby",
		Number: 10,
	},
	DefaultHP: 110,
	Defences: map[game.Colour]int{
		game.ColourRed:    4,
		game.ColourOrange: -1,
		game.ColourYellow: 2,
		game.ColourGreen:  0,
		game.ColourCyan:   1,
		game.ColourBlue:   -2,
		game.ColourViolet: 0,
		game.ColourPink:   0,
		game.ColourGray:   0,
		game.ColourBrown:  1,
		game.ColourBlack:  0,
		game.ColourWhite:  0,
	},
	SkillData: [4]*game.SkillData{
		SkillDance,
		SkillRage,
		SkillStop,
		SkillExecute,
	},
}
