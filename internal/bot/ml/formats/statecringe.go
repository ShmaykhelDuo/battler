package formats

import (
	"fmt"
	"slices"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/euphoria"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/milana"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/ruby"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/speed"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/storyteller"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/structure"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/z89"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

func boolToInt64(a bool) int64 {
	if a {
		return 1
	}

	return 0
}

func UpdateMapDefenceCringe(def map[game.Colour]int, m map[string]any, prefix string) {
	for c, s := range colourMap {
		m[prefix+"_defence_"+s] = float32(def[c])
	}
}

func UpdateMapSkillCringe(s *game.Skill, c, opp *game.Character, turnState game.TurnState, m map[string]any, prefix string) {
	m[prefix+"_colour"] = colourMap[s.Desc().Colour]
	m[prefix+"_cooldown"] = float32(s.Cooldown())
	m[prefix+"_unlock_turn"] = float32(s.UnlockTurn(c))
	m[prefix+"_is_available"] = boolToInt64(s.IsAvailable(c, opp, turnState))
}

func UpdateMapEffectsCringe(c *game.Character, turnState game.TurnState, m map[string]any, prefix string) {
	s := effectsState{}

	for _, e := range c.Effects() {
		switch eff := e.(type) {
		case storyteller.EffectCannotUse:
			s.StorytellerCannotUseColour = eff.Colour()
		case storyteller.EffectControlled:
			s.StorytellerControlled = true
		case *z89.EffectUltimateSlow:
			s.Z89UltimateSlowAmount = eff.Amount()
		case euphoria.EffectEuphoricSource:
			s.EuphoriaEuphoricSourceAmount = eff.Amount()
		case *euphoria.EffectUltimateEarly:
			s.EuphoriaUltimateEarlyAmount = eff.Amount()
		case euphoria.EffectEuphoricHeal:
			s.EuphoriaEuphoricHeal = true
		case ruby.EffectDoubleDamage:
			s.RubyDoubleDamageTurnsLeft = eff.TurnsLeft(turnState)
		case ruby.EffectCannotHeal:
			s.RubyCannotHealTurnsLeft = eff.TurnsLeft(turnState)
		case speed.EffectGreenTokens:
			s.SpeedGreenTokensNumber = eff.Amount()
		case speed.EffectBlackTokens:
			s.SpeedBlackTokensNumber = eff.Amount()
		case *speed.EffectDamageReduced:
			s.SpeedDamageReduced = true
		case speed.EffectDefenceReduced:
			s.SpeedDefenceReduced = true
		case speed.EffectSpedUp:
			s.SpeedSpedUp = true
		case milana.EffectStolenHP:
			s.MilanaStolenHPAmount = eff.Amount()
		case milana.EffectMintMist:
			s.MilanaMintMistTurnsLeft = eff.TurnsLeft(turnState)
		case *structure.EffectIBoost:
			s.StructureIBoostAmount = eff.Amount()
		case structure.EffectSLayers:
			s.StructureSLayersThreshold = eff.Threshold()
		case structure.EffectLastChance:
			s.StructureLastChance = true
		}
	}

	m[prefix+"_storyteller_cannotuse_colour"] = colourMap[s.StorytellerCannotUseColour]
	m[prefix+"_storyteller_controlled"] = boolToInt64(s.StorytellerControlled)
	m[prefix+"_z89_ultimateslow_amount"] = float32(s.Z89UltimateSlowAmount)
	m[prefix+"_euphoria_euphoricsource_amount"] = float32(s.EuphoriaEuphoricSourceAmount)
	m[prefix+"_euphoria_ultimateearly_amount"] = float32(s.EuphoriaUltimateEarlyAmount)
	m[prefix+"_euphoria_euphoricheal"] = boolToInt64(s.EuphoriaEuphoricHeal)
	m[prefix+"_ruby_doubledamage_turnsleft"] = float32(s.RubyDoubleDamageTurnsLeft)
	m[prefix+"_ruby_cannotheal_turnsleft"] = float32(s.RubyCannotHealTurnsLeft)
	m[prefix+"_speed_greentokens_number"] = float32(s.SpeedGreenTokensNumber)
	m[prefix+"_speed_blacktokens_number"] = float32(s.SpeedBlackTokensNumber)
	m[prefix+"_speed_damagereduced"] = boolToInt64(s.SpeedDamageReduced)
	m[prefix+"_speed_defencereduced"] = boolToInt64(s.SpeedDefenceReduced)
	m[prefix+"_speed_spedup"] = boolToInt64(s.SpeedSpedUp)
	m[prefix+"_milana_stolenhp_amount"] = float32(s.MilanaStolenHPAmount)
	m[prefix+"_milana_mintmist_turnsleft"] = float32(s.MilanaMintMistTurnsLeft)
	m[prefix+"_structure_iboost_amount"] = float32(s.StructureIBoostAmount)
	m[prefix+"_structure_slayers_threshold"] = float32(s.StructureSLayersThreshold)
	m[prefix+"_structure_lastchance"] = boolToInt64(s.StructureLastChance)
}

func UpdateMapCharCringe(c *game.Character, opp *game.Character, turnState game.TurnState, m map[string]any, prefix string) {
	m[prefix+"_number"] = int64(c.Desc().Number)
	m[prefix+"_hp"] = float32(c.HP())
	m[prefix+"_maxhp"] = float32(c.MaxHP())
	UpdateMapDefenceCringe(c.Defences(), m, prefix+"_defence")
	UpdateMapEffectsCringe(c, turnState, m, prefix+"_effect")

	s := c.Skills()

	for i := range 4 {
		p := fmt.Sprintf("%s_skill%d", prefix, i)
		UpdateMapSkillCringe(s[i], c, opp, turnState, m, p)
	}

	m[prefix+"_lastusedskill"] = int64(slices.Index(s[:], c.LastUsedSkill()))
}

func GetMapStateCringe(in match.GameState) map[string]any {
	res := make(map[string]any)

	UpdateMapCharCringe(in.Character, in.Opponent, in.TurnState, res, "char")
	UpdateMapCharCringe(in.Opponent, in.Character, in.TurnState, res, "opp")

	res["turnnum"] = float32(in.TurnState.TurnNum)
	res["goingfirst"] = boolToInt64(in.TurnState.IsGoingFirst)
	res["asopp"] = float32(boolToInt64(in.AsOpp))

	return res
}
