package storyteller

import (
	"math"

	"github.com/ShmaykhelDuo/battler/internal/game"
)

const (
	SkillYourNumberIndex int = iota
	SkillYourColourIndex
	SkillYourDreamIndex
	SkillMyStoryIndex
)

// Deal 10 + the remainder of your opponent's number divided by 7 Orange damage.
var SkillYourNumber = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Your Number",
		IsUltimate: false,
		Colour:     game.ColourOrange,
	},
	Use: func(c, opp *game.Character, turnState game.TurnState) {
		dmg := 10 + opp.Desc().Number%7
		c.Damage(opp, dmg, game.ColourOrange)
	},
}

// Next turn, your opponent can't use the skills of the same colour they used last. Deal 15 damage of that colour.
// Unlocks when your opponent uses a skill.
// Cooldown: 1.
var SkillYourColour = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Your Colour",
		IsUltimate: false,
		Colour:     game.ColourWhite,
	},
	Cooldown: 1,
	Use: func(c, opp *game.Character, turnState game.TurnState) {
		colour := opp.LastUsedSkill().Desc().Colour
		opp.AddEffect(NewEffectCannotUse(turnState, colour))
		c.Damage(opp, 15, colour)
	},
	IsAvailable: func(c, opp *game.Character, turnState game.TurnState) bool {
		return opp.LastUsedSkill() != nil
	},
}

// Heal for (your max HP - your opponent's number) / your turn number.
// If your opponent's number is more than 83, subtract a flat number as if it was 83.
var SkillYourDream = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Your Dream",
		IsUltimate: false,
		Colour:     game.ColourViolet,
	},
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		heal := yourDreamHeal(c, opp, turnState)
		c.Heal(heal)
	},
	IsAppropriate: func(c, opp *game.Character, turnState game.TurnState) bool {
		heal := yourDreamHeal(c, opp, turnState)
		return !(c.CanHeal(heal) || c.HP() >= c.MaxHP())
	},
}

func yourDreamHeal(c *game.Character, opp *game.Character, turnState game.TurnState) int {
	num := opp.Desc().Number
	if num > 83 {
		num = 83
	}

	return int(math.Ceil(float64(c.MaxHP()-num) / float64(turnState.TurnNum)))
}

// Next turn, you decide which skills your opponent uses.
// Unlocks on turn 7.
// Cooldown 1.
var SkillMyStory = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "My Story",
		IsUltimate: true,
		Colour:     game.ColourBlue,
	},
	Cooldown:   1,
	UnlockTurn: 7,
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		opp.AddEffect(NewEffectControlled(turnState))
	},
}
