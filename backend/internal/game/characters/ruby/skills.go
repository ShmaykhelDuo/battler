package ruby

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

// Double all of your damage.
// Lasts 2 turns.
var SkillDance = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Dance",
		IsUltimate: false,
		Colour:     game.ColourYellow,
	},
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		c.AddEffect(NewEffectDoubleDamage(turnState))
	},
}

// Deal 24 - 2 * your turn number Red damage.
var SkillRage = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Rage",
		IsUltimate: false,
		Colour:     game.ColourRed,
	},
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		dmg := 24 - 2*turnState.TurnNum
		c.Damage(opp, dmg, game.ColourRed)
	},
}

// Every player can not heal until the end of their next turn.
// While this is active for you, .Execute becomes stronger.
// Cooldown 1.
var SkillStop = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Stop",
		IsUltimate: false,
		Colour:     game.ColourCyan,
	},
	Cooldown: 1,
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		c.AddEffect(NewEffectCannotHeal(turnState, false))
		opp.AddEffect(NewEffectCannotHeal(turnState, true))
	},
}

// If your opponent's at less than 10% of their max HP, defeat them instantly.
// While Stop effect is active, the threshold goes to 20% of opponent's max hp.
var SkillExecute = game.SkillData{
	Desc: game.SkillDescription{
		Name:       ".Execute",
		IsUltimate: true,
		Colour:     game.ColourRed,
	},
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		threshold := 0.1

		if hasEffectCannotHeal(c) {
			threshold = 0.2
		}

		if float64(opp.HP()) < threshold*float64(opp.MaxHP()) {
			opp.Kill()
		}
	},
}

func hasEffectCannotHeal(c *game.Character) bool {
	_, found := game.CharacterEffect[EffectCannotHeal](c, EffectDescCannotHeal)
	return found
}
