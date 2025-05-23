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

var colourMap = map[game.Colour]string{
	game.ColourRed:    "red",
	game.ColourOrange: "orange",
	game.ColourYellow: "yellow",
	game.ColourGreen:  "green",
	game.ColourCyan:   "cyan",
	game.ColourBlue:   "blue",
	game.ColourViolet: "violet",
	game.ColourPink:   "pink",
	game.ColourGray:   "gray",
	game.ColourBrown:  "brown",
	game.ColourBlack:  "black",
	game.ColourWhite:  "white",
}

func UpdateMapDefence(def map[game.Colour]int, m map[string]any, prefix string) {
	for c, s := range colourMap {
		m[prefix+"_defence_"+s] = int64(def[c])
	}
}

func UpdateMapSkill(s *game.Skill, c, opp *game.Character, turnState game.TurnState, m map[string]any, prefix string) {
	m[prefix+"_colour"] = colourMap[s.Desc().Colour]
	m[prefix+"_cooldown"] = int64(s.Cooldown())
	m[prefix+"_unlock_turn"] = int64(s.UnlockTurn(c))
	m[prefix+"_is_available"] = s.IsAvailable(c, opp, turnState)
}

type effectsState struct {
	StorytellerCannotUseColour   game.Colour
	StorytellerControlled        bool
	Z89UltimateSlowAmount        int
	EuphoriaEuphoricSourceAmount int
	EuphoriaUltimateEarlyAmount  int
	EuphoriaEuphoricHeal         bool
	RubyDoubleDamageTurnsLeft    int
	RubyCannotHealTurnsLeft      int
	SpeedGreenTokensNumber       int
	SpeedBlackTokensNumber       int
	SpeedDamageReduced           bool
	SpeedDefenceReduced          bool
	SpeedSpedUp                  bool
	MilanaStolenHPAmount         int
	MilanaMintMistTurnsLeft      int
	StructureIBoostAmount        int
	StructureSLayersThreshold    int
	StructureLastChance          bool
}

func UpdateMapEffects(c *game.Character, turnState game.TurnState, m map[string]any, prefix string) {
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
	m[prefix+"_storyteller_controlled"] = s.StorytellerControlled
	m[prefix+"_z89_ultimateslow_amount"] = int64(s.Z89UltimateSlowAmount)
	m[prefix+"_euphoria_euphoricsource_amount"] = int64(s.EuphoriaEuphoricSourceAmount)
	m[prefix+"_euphoria_ultimateearly_amount"] = int64(s.EuphoriaUltimateEarlyAmount)
	m[prefix+"_euphoria_euphoricheal"] = s.EuphoriaEuphoricHeal
	m[prefix+"_ruby_doubledamage_turnsleft"] = int64(s.RubyDoubleDamageTurnsLeft)
	m[prefix+"_ruby_cannotheal_turnsleft"] = int64(s.RubyCannotHealTurnsLeft)
	m[prefix+"_speed_greentokens_number"] = int64(s.SpeedGreenTokensNumber)
	m[prefix+"_speed_blacktokens_number"] = int64(s.SpeedBlackTokensNumber)
	m[prefix+"_speed_damagereduced"] = s.SpeedDamageReduced
	m[prefix+"_speed_defencereduced"] = s.SpeedDefenceReduced
	m[prefix+"_speed_spedup"] = s.SpeedSpedUp
	m[prefix+"_milana_stolenhp_amount"] = int64(s.MilanaStolenHPAmount)
	m[prefix+"_milana_mintmist_turnsleft"] = int64(s.MilanaMintMistTurnsLeft)
	m[prefix+"_structure_iboost_amount"] = int64(s.StructureIBoostAmount)
	m[prefix+"_structure_slayers_threshold"] = int64(s.StructureSLayersThreshold)
	m[prefix+"_structure_lastchance"] = s.StructureLastChance
}

func UpdateMapChar(c *game.Character, opp *game.Character, turnState game.TurnState, m map[string]any, prefix string) {
	m[prefix+"_number"] = int64(c.Desc().Number)
	m[prefix+"_hp"] = int64(c.HP())
	m[prefix+"_maxhp"] = int64(c.MaxHP())
	UpdateMapDefence(c.Defences(), m, prefix+"_defence")
	UpdateMapEffects(c, turnState, m, prefix+"_effect")

	s := c.Skills()

	for i := range 4 {
		p := fmt.Sprintf("%s_skill%d", prefix, i)
		UpdateMapSkill(s[i], c, opp, turnState, m, p)
	}

	m[prefix+"_lastusedskill"] = int64(slices.Index(s[:], c.LastUsedSkill()))
}

func GetMapState(in match.GameState) map[string]any {
	res := make(map[string]any)

	UpdateMapChar(in.Character, in.Opponent, in.TurnState, res, "char")
	UpdateMapChar(in.Opponent, in.Character, in.TurnState, res, "opp")

	res["turnnum"] = int64(in.TurnState.TurnNum)
	res["goingfirst"] = in.TurnState.IsGoingFirst
	res["asopp"] = in.AsOpp

	return res
}
