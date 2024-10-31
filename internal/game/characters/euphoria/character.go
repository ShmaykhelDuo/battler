package euphoria

import "github.com/ShmaykhelDuo/battler/internal/game"

var CharacterEuphoria = &game.CharacterData{
	Desc: game.CharacterDescription{
		Name:   "Euphoria",
		Number: 9,
	},
	DefaultHP: 117,
	Defences: map[game.Colour]int{
		game.ColourRed:    0,
		game.ColourOrange: 2,
		game.ColourYellow: 0,
		game.ColourGreen:  0,
		game.ColourCyan:   -3,
		game.ColourBlue:   0,
		game.ColourViolet: 0,
		game.ColourPink:   3,
		game.ColourGray:   0,
		game.ColourBrown:  -4,
		game.ColourBlack:  0,
		game.ColourWhite:  0,
	},
	SkillData: [4]*game.SkillData{
		SkillAmpleness,
		SkillExuberance,
		SkillPinkSphere,
		SkillEuphoria,
	},
}
