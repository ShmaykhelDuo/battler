package storyteller

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

var CharacterStoryteller = game.CharacterData{
	Desc: game.CharacterDescription{
		Name:   "Storyteller",
		Number: 1,
	},
	DefaultHP: 119,
	Defences: map[game.Colour]int{
		game.ColourRed:    -1,
		game.ColourOrange: 1,
		game.ColourYellow: 0,
		game.ColourGreen:  -2,
		game.ColourCyan:   -1,
		game.ColourBlue:   1,
		game.ColourViolet: 1,
		game.ColourPink:   0,
		game.ColourGray:   -1,
		game.ColourBrown:  -1,
		game.ColourBlack:  -2,
		game.ColourWhite:  1,
	},
	SkillData: [4]game.SkillData{
		SkillYourNumber,
		SkillYourColour,
		SkillYourDream,
		SkillMyStory,
	},
}
