package z89

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

// Deal 12 Black damage, then set opponent's max HP to their current HP.
// Cooldown 1.
var SkillScarcity = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Scarcity",
		IsUltimate: false,
		Colour:     game.ColourBlack,
	},
	Cooldown: 1,
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		c.Damage(opp, 12, game.ColourBlack)
		opp.SetMaxHP(opp.HP())
	},
}

// If opponent's ultimate is not available yet, delay it for 1 turn.
// Can't be delayed later than their 10th turn.
// Cooldown 2.
// Unlocks on turn 2.
var SkillIndifference = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Indifference",
		IsUltimate: false,
		Colour:     game.ColourCyan,
	},
	Cooldown:   2,
	UnlockTurn: 2,
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		if hasOpponentUltimateUnlocked(opp, turnState) {
			return
		}

		increaseUltimateSlow(opp)
	},
	IsAppropriate: func(c, opp *game.Character, turnState game.TurnState) bool {
		return !hasOpponentUltimateUnlocked(opp, turnState)
	},
}

// Deal 15 - (opponent's max HP - opponent's current HP) green damage.
var SkillGreenSphere = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Green Sphere",
		IsUltimate: false,
		Colour:     game.ColourGreen,
	},
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		dmg := 15 - (opp.MaxHP() - opp.HP())
		c.Damage(opp, dmg, game.ColourGreen)
	},
}

// Deal 40 - (opponent's max HP - 70) Blue damage.
// Unlocks on turn 9.
var SkillDespondency = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Despondency",
		IsUltimate: true,
		Colour:     game.ColourBlue,
	},
	Cooldown:   10,
	UnlockTurn: 9,
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		dmg := 40 - (opp.MaxHP() - 70)
		c.Damage(opp, dmg, game.ColourBlue)
	},
}

func hasOpponentUltimateUnlocked(opp *game.Character, turnState game.TurnState) bool {
	unlockCtx := game.TurnState{
		TurnNum:      opp.Skills()[3].UnlockTurn(),
		IsGoingFirst: true,
	}
	return turnState.IsAfter(unlockCtx)
}

func increaseUltimateSlow(opp *game.Character) {
	effSlow, ok := game.CharacterEffect[*EffectUltimateSlow](opp, EffectDescUltimateSlow)
	if ok {
		effSlow.Increase()
		return
	}

	opp.AddEffect(NewEffectUltimateSlow())
}
