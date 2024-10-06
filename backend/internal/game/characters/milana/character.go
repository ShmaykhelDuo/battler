package milana

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

var CharacterMilana = game.CharacterData{
	Desc: game.CharacterDescription{
		Name:   "Milana",
		Number: 51,
	},
	DefaultHP: 114,
	Defences: map[game.Colour]int{
		game.ColourRed:    0,
		game.ColourOrange: -1,
		game.ColourYellow: 1,
		game.ColourGreen:  1,
		game.ColourCyan:   1,
		game.ColourBlue:   0,
		game.ColourViolet: -1,
		game.ColourPink:   0,
		game.ColourGray:   0,
		game.ColourBrown:  0,
		game.ColourBlack:  -2,
		game.ColourWhite:  2,
	},
	SkillData: [4]game.SkillData{
		SkillRoyalMove,
		SkillComposure,
		SkillMintMist,
		SkillPride,
	},
}
