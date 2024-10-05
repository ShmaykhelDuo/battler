package euphoria

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

// Increases Euphoric Source and everyone's max HP by 12.
// Cooldown: 1.
var SkillAmpleness = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Ampleness",
		IsUltimate: false,
		Colour:     game.ColourOrange,
	},
	Cooldown: 1,
	Use: func(c *game.Character, opp *game.Character, gameCtx game.Context) {
		increaseMaxHP(c, opp, 12)
		increaseEuphoricSource(c, 12)
	},
}

// If your opponent's ultimate is not unlocked yet, add 10 to Euphoric Source, everyone's max HP and your current HP.
// Also, your opponent's ultimate unlocks 1 turn earlier.
// If it already is unlocked, add 20 instead.
// Cooldown 2.
var SkillExuberance = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Exuberance",
		IsUltimate: false,
		Colour:     game.ColourOrange,
	},
	Cooldown: 2,
	Use: func(c *game.Character, opp *game.Character, gameCtx game.Context) {
		amount := 20

		if !isSkillUnlocked(gameCtx, opp.Skills()[3]) {
			amount = 10
			increaseUltimateEarly(opp)
		}

		increaseMaxHP(c, opp, amount)
		c.Heal(amount)
		increaseEuphoricSource(c, amount)
	},
}

// Deal 12 Pink damage. Also, increase everyone's max HP by 12.
var SkillPinkSphere = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Pink Sphere",
		IsUltimate: false,
		Colour:     game.ColourPink,
	},
	Use: func(c *game.Character, opp *game.Character, gameCtx game.Context) {
		c.Damage(opp, 12, game.ColourPink)
		increaseMaxHP(c, opp, 12)
	},
}

// Heal everyone by the amount in Euphoric Source at the end each turn.
// Every turn, Source gets depleted by 9. Lasts until the end of the game.
// Starting turn: 4.
var SkillEuphoria = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Euphoria",
		IsUltimate: true,
		Colour:     game.ColourPink,
	},
	Cooldown:   20,
	UnlockTurn: 4,
	Use: func(c *game.Character, opp *game.Character, gameCtx game.Context) {
		c.AddEffect(EffectEuphoricHeal{})
	},
}

func increaseMaxHP(c, opp *game.Character, amount int) {
	c.SetMaxHP(c.MaxHP() + amount)
	opp.SetMaxHP(opp.MaxHP() + amount)
}

func increaseEuphoricSource(c *game.Character, amount int) {
	eff := c.Effect(EffectDescEuphoricSource)
	source, ok := eff.(*EffectEuphoricSource)
	if ok {
		source.Increase(amount)
		return
	}

	c.AddEffect(NewEffectEuphoricSource(amount))
}

func isSkillUnlocked(gameCtx game.Context, s *game.Skill) bool {
	if gameCtx.TurnNum == s.UnlockTurn() {
		return !gameCtx.IsGoingFirst
	}

	return gameCtx.TurnNum > s.UnlockTurn()
}

func increaseUltimateEarly(opp *game.Character) {
	eff := opp.Effect(EffectDescUltimateEarly)
	source, ok := eff.(*EffectUltimateEarly)
	if ok {
		source.Increase()
		return
	}

	opp.AddEffect(NewEffectUltimateEarly())
}
