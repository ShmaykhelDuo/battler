package speed

import "github.com/ShmaykhelDuo/battler/internal/game"

var CharacterSpeed = &game.CharacterData{
	Desc: game.CharacterDescription{
		Name:   "Speed",
		Number: 33,
	},
	DefaultHP: 113,
	Defences: map[game.Colour]int{
		game.ColourRed:    0,
		game.ColourOrange: 0,
		game.ColourYellow: 0,
		game.ColourGreen:  4,
		game.ColourCyan:   0,
		game.ColourBlue:   0,
		game.ColourViolet: -2,
		game.ColourPink:   0,
		game.ColourGray:   0,
		game.ColourBrown:  0,
		game.ColourBlack:  2,
		game.ColourWhite:  -2,
	},
	SkillData: [4]*game.SkillData{
		SkillRun,
		SkillWeaken,
		SkillSpeed,
		SkillStab,
	},
}
