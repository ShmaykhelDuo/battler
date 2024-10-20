package milana

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

// Deal 12 Green damage and add that to Stolen HP.
// With Mint Mist, deal 20 Green damage instead.
var SkillRoyalMove = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Royal Move",
		IsUltimate: false,
		Colour:     game.ColourGreen,
	},
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		dmg := 12
		_, hasMintMist := game.CharacterEffect[EffectMintMist](c, EffectDescMintMist)
		if hasMintMist {
			dmg = 20
		}

		stolen := c.Damage(opp, dmg, game.ColourGreen)

		effStolen, ok := game.CharacterEffect[EffectStolenHP](c, EffectDescStolenHP)
		if ok {
			effStolen.Increase(stolen)
		} else {
			c.AddEffect(NewEffectStolenHP(stolen))
		}
	},
}

// Spend some Stolen HP to heal yourself for up to 20.
// With Mint Mist, heal up to 30.
var SkillComposure = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Composure",
		IsUltimate: false,
		Colour:     game.ColourWhite,
	},
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		effStolen, ok := game.CharacterEffect[EffectStolenHP](c, EffectDescStolenHP)
		if !ok {
			return
		}

		stolen := effStolen.Amount()

		cost := 6
		heal := 20

		_, hasMintMist := game.CharacterEffect[EffectMintMist](c, EffectDescMintMist)
		if hasMintMist {
			cost = 10
			heal = 30
		}

		if stolen < cost {
			cost = c.Heal(stolen)
		} else {
			c.Heal(heal)
		}

		effStolen.Decrease(cost)
	},
	IsAppropriate: func(c, opp *game.Character, turnState game.TurnState) bool {
		_, ok := game.CharacterEffect[EffectStolenHP](c, EffectDescStolenHP)
		return ok
	},
}

// You become invisible, your opponent can't debuff you.
// Your Royal Move and Composure become stronger.
// Lasts 2 turns.
// Cooldown: 2.
var SkillMintMist = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Mint Mist",
		IsUltimate: false,
		Colour:     game.ColourWhite,
	},
	Cooldown: 2,
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		c.AddEffect(NewEffectMintMist(turnState))
	},
}

// Spend all of your Stolen HP to deal as much Cyan damage.
// Unlocks on turn 8.
var SkillPride = game.SkillData{
	Desc: game.SkillDescription{
		Name:       "Pride",
		IsUltimate: true,
		Colour:     game.ColourCyan,
	},
	UnlockTurn: 8,
	Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		effStolen, ok := game.CharacterEffect[EffectStolenHP](c, EffectDescStolenHP)
		if !ok {
			return
		}

		dmg := effStolen.Amount()
		c.Damage(opp, dmg, game.ColourCyan)
		effStolen.Decrease(dmg)
	},
	IsAppropriate: func(c, opp *game.Character, turnState game.TurnState) bool {
		_, ok := game.CharacterEffect[EffectStolenHP](c, EffectDescStolenHP)
		return ok
	},
}
