package speed

import (
	"github.com/ShmaykhelDuo/battler/backend/internal/game"
)

// Your opponent's next attack will deal 5 less damage. Gain a Green Token.
var SkillRun = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Run",
		IsUltimate: false,
		Colour:     game.ColourGreen,
	},
	Use: func(c *game.Character, opp *game.Character, gameCtx game.Context) {
		redAmount := 5

		e := c.Effect(EffectDescDamageReduced)
		red, ok := e.(*EffectDamageReduced)
		if ok {
			red.Increase(redAmount)
		} else {
			c.AddEffect(NewEffectDamageReduced(redAmount))
		}

		if greenTokensNumber(c) < 5 {
			increaseGreenTokens(c)
		}
	},
}

// Reduce your opponent's defense to Green by 1. Gain a Black Token.
var SkillWeaken = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Weaken",
		IsUltimate: false,
		Colour:     game.ColourBlack,
	},
	Use: func(c *game.Character, opp *game.Character, gameCtx game.Context) {
		opp.AddEffect(EffectDefenceReduced{})

		if blackTokensNumber(c) < 5 {
			increaseBlackTokens(c)
		}
	},
}

// Next turn, you'll use two skills but not your ultimate. Gain a Green Token.
var SkillSpeed = game.SkillData{
	Desc:       game.SkillDescription{},
	Cooldown:   0,
	UnlockTurn: 0,
	Use: func(c *game.Character, opp *game.Character, gameCtx game.Context) {
		c.AddEffect(EffectSpedUp{})

		if greenTokensNumber(c) < 5 {
			increaseGreenTokens(c)
		}
	},
}

// For each of your tokens, deal 6 Green&Black damage.
var SkillStab = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Stab",
		IsUltimate: true,
		Colour:     game.ColourBlack,
	},
	Use: func(c *game.Character, opp *game.Character, gameCtx game.Context) {
		mul := 6

		green := greenTokensNumber(c)
		c.Damage(opp, green*mul, game.ColourGreen)

		black := blackTokensNumber(c)
		c.Damage(opp, black*mul, game.ColourBlack)
	},
}

func tokens(c *game.Character, desc game.EffectDescription) *EffectTokens {
	eff := c.Effect(desc)
	source, _ := eff.(*EffectTokens)
	return source
}

func greenTokensNumber(c *game.Character) int {
	source := tokens(c, EffectDescGreenTokens)
	if source == nil {
		return 0
	}

	return source.Number()
}

func increaseGreenTokens(c *game.Character) {
	source := tokens(c, EffectDescGreenTokens)
	if source != nil {
		source.Increase()
		return
	}

	c.AddEffect(NewEffectGreenTokens(1))
}

func blackTokensNumber(c *game.Character) int {
	source := tokens(c, EffectDescBlackTokens)
	if source == nil {
		return 0
	}

	return source.Number()
}

func increaseBlackTokens(c *game.Character) {
	source := tokens(c, EffectDescBlackTokens)
	if source != nil {
		source.Increase()
		return
	}

	c.AddEffect(NewEffectBlackTokens(1))
}
