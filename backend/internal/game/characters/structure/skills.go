package structure

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

// Deal Cyan damage. Base damage is 5, gets to 10, 15 and 20 when boosted by I Boost.
var SkillEShock = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "E-Shock",
		IsUltimate: false,
		Colour:     game.ColourCyan,
	},
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		dmg := 5

		boost, ok := game.CharacterEffect[*EffectIBoost](c, EffectDescIBoost)
		if ok {
			dmg += boost.Amount()
		}

		c.Damage(opp, dmg, game.ColourCyan)
	},
}

// Boost your S Layers threshold by 5 and E-Shock damage by 5.
// Can only be used three times in a match.
var SkillIBoost = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "I Boost",
		IsUltimate: false,
		Colour:     game.ColourViolet,
	},
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		boost, ok := game.CharacterEffect[*EffectIBoost](c, EffectDescIBoost)
		if !ok {
			c.AddEffect(NewEffectIBoost(5))
			return
		}

		boost.Increase()
	},
	IsAvailable: func(c *game.Character, opp *game.Character, turnState game.TurnState) bool {
		boost, ok := game.CharacterEffect[*EffectIBoost](c, EffectDescIBoost)
		if !ok {
			return true
		}

		return boost.Amount() < 15
	},
}

// Next turn, your opponent can't damage you unless they deal more than a certain threshold.
// Thresholds are 5, 10, 15 and 20.
// Gain 1 Defense against all colours but Black.
var SkillSLayers = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "S Layers",
		IsUltimate: false,
		Colour:     game.ColourGray,
	},
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		threshold := 5

		boost, ok := game.CharacterEffect[*EffectIBoost](c, EffectDescIBoost)
		if ok {
			threshold += boost.Amount()
		}

		c.AddEffect(NewEffectSLayers(turnState, threshold))
	},
}

// If you survive your opponent's next turn, fully heal.
// Unlocks on turn 7.
// Can only be used once per match.
var SkillLastChance = game.SkillData{
	Desc:       game.SkillDescription{},
	Cooldown:   10,
	UnlockTurn: 7,
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		c.AddEffect(NewEffectLastChance(turnState))
	},
}
