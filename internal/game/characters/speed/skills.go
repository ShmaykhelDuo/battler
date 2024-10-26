package speed

import (
	"github.com/ShmaykhelDuo/battler/internal/game"
)

// Your opponent's next attack will deal 5 less damage. Gain a Green Token.
var SkillRun = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Run",
		IsUltimate: false,
		Colour:     game.ColourGreen,
	},
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		redAmount := 5

		red, ok := game.CharacterEffect[*EffectDamageReduced](c, EffectDescDamageReduced)
		if ok {
			red.Increase(redAmount)
		} else {
			c.AddEffect(NewEffectDamageReduced(redAmount))
		}

		if greenTokensNumber(c) < 5 {
			increaseGreenTokens(c)
		}
	},
	IsAppropriate: func(c, opp *game.Character, turnState game.TurnState) bool {
		return greenTokensNumber(c) < 5
	},
}

// Reduce your opponent's defense to Green by 1. Gain a Black Token.
var SkillWeaken = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Weaken",
		IsUltimate: false,
		Colour:     game.ColourBlack,
	},
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		opp.AddEffect(EffectDefenceReduced{})

		if blackTokensNumber(c) < 5 {
			increaseBlackTokens(c)
		}
	},
	IsAppropriate: func(c, opp *game.Character, turnState game.TurnState) bool {
		return blackTokensNumber(c) < 5
	},
}

// Next turn, you'll use two skills but not your ultimate. Gain a Green Token.
var SkillSpeed = game.SkillData{
	Desc:       game.SkillDescription{},
	Cooldown:   0,
	UnlockTurn: 0,
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		c.AddEffect(NewEffectSpedUp(turnState))

		if greenTokensNumber(c) < 5 {
			increaseGreenTokens(c)
		}
	},
	IsAppropriate: func(c, opp *game.Character, turnState game.TurnState) bool {
		eff, ok := game.CharacterEffect[EffectSpedUp](c, EffectDescSpedUp)
		if ok && eff.TurnsLeft(turnState) > 1 {
			return false
		}

		return greenTokensNumber(c) < 5 || blackTokensNumber(c) < 5
	},
}

// For each of your tokens, deal 6 Green&Black damage.
var SkillStab = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Stab",
		IsUltimate: true,
		Colour:     game.ColourBlack,
	},
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		mul := 6

		green := greenTokensNumber(c)
		c.Damage(opp, green*mul, game.ColourGreen)

		black := blackTokensNumber(c)
		c.Damage(opp, black*mul, game.ColourBlack)
	},
	IsAppropriate: func(c, opp *game.Character, turnState game.TurnState) bool {
		return greenTokensNumber(c) > 0 || blackTokensNumber(c) > 0
	},
}

func greenTokensNumber(c *game.Character) int {
	tokens, ok := game.CharacterEffect[EffectGreenTokens](c, EffectDescGreenTokens)
	if !ok {
		return 0
	}

	return tokens.Amount()
}

func increaseGreenTokens(c *game.Character) {
	tokens, ok := game.CharacterEffect[EffectGreenTokens](c, EffectDescGreenTokens)
	if ok {
		tokens.Increase(1)
		return
	}

	c.AddEffect(NewEffectGreenTokens(1))
}

func blackTokensNumber(c *game.Character) int {
	tokens, ok := game.CharacterEffect[EffectBlackTokens](c, EffectDescBlackTokens)
	if !ok {
		return 0
	}

	return tokens.Amount()
}

func increaseBlackTokens(c *game.Character) {
	tokens, ok := game.CharacterEffect[EffectBlackTokens](c, EffectDescBlackTokens)
	if ok {
		tokens.Increase(1)
		return
	}

	c.AddEffect(NewEffectBlackTokens(1))
}
