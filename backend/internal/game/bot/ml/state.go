package ml

import (
	"slices"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/euphoria"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/milana"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/ruby"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/speed"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/storyteller"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/structure"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/z89"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/match"
)

type DefenceState struct {
	Red    int
	Orange int
	Yellow int
	Green  int
	Cyan   int
	Blue   int
	Violet int
	Pink   int
	Gray   int
	Brown  int
	Black  int
	White  int
}

func NewDefenceState(def map[game.Colour]int) DefenceState {
	return DefenceState{
		Red:    def[game.ColourRed],
		Orange: def[game.ColourOrange],
		Yellow: def[game.ColourYellow],
		Green:  def[game.ColourGreen],
		Cyan:   def[game.ColourCyan],
		Blue:   def[game.ColourBlue],
		Violet: def[game.ColourViolet],
		Pink:   def[game.ColourPink],
		Gray:   def[game.ColourGray],
		Brown:  def[game.ColourBrown],
		Black:  def[game.ColourBlack],
		White:  def[game.ColourWhite],
	}
}

func (s DefenceState) AppendSlice(res []int) []int {
	return append(
		res,
		s.Red,
		s.Orange,
		s.Yellow,
		s.Green,
		s.Cyan,
		s.Blue,
		s.Violet,
		s.Pink,
		s.Gray,
		s.Brown,
		s.Black,
		s.White,
	)
}

type SkillState struct {
	Colour      int
	Cooldown    int
	UnlockTurn  int
	IsAvailable int
}

func NewSkillState(s *game.Skill, opp *game.Character, turnState game.TurnState) SkillState {
	res := SkillState{
		Colour:     int(s.Desc().Colour),
		Cooldown:   s.Cooldown(),
		UnlockTurn: s.UnlockTurn(),
	}

	if s.IsAvailable(opp, turnState) {
		res.IsAvailable = 1
	}

	return res
}

func (s SkillState) AppendSlice(res []int) []int {
	return append(
		res,
		s.Colour,
		s.Cooldown,
		s.UnlockTurn,
		s.IsAvailable,
	)
}

type EffectsState struct {
	StorytellerCannotUseColour   int
	StorytellerControlled        int
	Z89UltimateSlowAmount        int
	EuphoriaEuphoricSourceAmount int
	EuphoriaUltimateEarlyAmount  int
	EuphoriaEuphoricHeal         int
	RubyDoubleDamageTurnsLeft    int
	RubyCannotHealTurnsLeft      int
	SpeedGreenTokensNumber       int
	SpeedBlackTokensNumber       int
	SpeedDamageReduced           int
	SpeedDefenceReduced          int
	SpeedSpedUp                  int
	MilanaStolenHPAmount         int
	MilanaMintMistTurnsLeft      int
	StructureIBoostAmount        int
	StructureSLayersThreshold    int
	StructureLastChance          int
}

func NewEffectsState(c *game.Character, turnState game.TurnState) EffectsState {
	res := EffectsState{}

	for _, e := range c.Effects() {
		switch eff := e.(type) {
		case storyteller.EffectCannotUse:
			res.StorytellerCannotUseColour = int(eff.Colour())
		case storyteller.EffectControlled:
			res.StorytellerControlled = 1
		case *z89.EffectUltimateSlow:
			res.Z89UltimateSlowAmount = eff.Amount()
		case euphoria.EffectEuphoricSource:
			res.EuphoriaEuphoricSourceAmount = eff.Amount()
		case *euphoria.EffectUltimateEarly:
			res.EuphoriaUltimateEarlyAmount = eff.Amount()
		case euphoria.EffectEuphoricHeal:
			res.EuphoriaEuphoricHeal = 1
		case ruby.EffectDoubleDamage:
			res.RubyDoubleDamageTurnsLeft = eff.TurnsLeft(turnState)
		case ruby.EffectCannotHeal:
			res.RubyCannotHealTurnsLeft = eff.TurnsLeft(turnState)
		case speed.EffectGreenTokens:
			res.SpeedGreenTokensNumber = eff.Amount()
		case speed.EffectBlackTokens:
			res.SpeedBlackTokensNumber = eff.Amount()
		case *speed.EffectDamageReduced:
			res.SpeedDamageReduced = 1
		case milana.EffectStolenHP:
			res.MilanaStolenHPAmount = eff.Amount()
		case milana.EffectMintMist:
			res.MilanaMintMistTurnsLeft = eff.TurnsLeft(turnState)
		case *structure.EffectIBoost:
			res.StructureIBoostAmount = eff.Amount()
		case structure.EffectSLayers:
			res.StructureSLayersThreshold = eff.Threshold()
		case structure.EffectLastChance:
			res.StructureLastChance = 1
		}
	}

	return res
}

func (s EffectsState) AppendSlice(res []int) []int {
	return append(
		res,
		s.StorytellerCannotUseColour,
		s.StorytellerControlled,
		s.Z89UltimateSlowAmount,
		s.EuphoriaEuphoricSourceAmount,
		s.EuphoriaUltimateEarlyAmount,
		s.EuphoriaEuphoricHeal,
		s.RubyDoubleDamageTurnsLeft,
		s.RubyCannotHealTurnsLeft,
		s.SpeedGreenTokensNumber,
		s.SpeedBlackTokensNumber,
		s.SpeedDamageReduced,
		s.SpeedDefenceReduced,
		s.SpeedSpedUp,
		s.MilanaStolenHPAmount,
		s.MilanaMintMistTurnsLeft,
		s.StructureIBoostAmount,
		s.StructureSLayersThreshold,
		s.StructureLastChance,
	)
}

type CharState struct {
	Number        int
	HP            int
	MaxHP         int
	Defences      DefenceState
	Skills        [4]SkillState
	LastUsedSkill int
	Effects       EffectsState
}

func NewCharState(c *game.Character, opp *game.Character, turnState game.TurnState) CharState {
	res := CharState{
		Number:   c.Desc().Number,
		HP:       c.HP(),
		MaxHP:    c.MaxHP(),
		Defences: NewDefenceState(c.Defences()),
		Effects:  NewEffectsState(c, turnState),
	}

	s := c.Skills()

	for i := range 4 {
		res.Skills[i] = NewSkillState(s[i], opp, turnState)
	}

	res.LastUsedSkill = slices.Index(s[:], c.LastUsedSkill())

	return res
}

func (s CharState) AppendSlice(res []int) []int {
	res = append(res, s.Number, s.HP, s.MaxHP)
	res = s.Defences.AppendSlice(res)
	for _, ss := range s.Skills {
		res = ss.AppendSlice(res)
	}
	res = append(res, s.LastUsedSkill)
	res = s.Effects.AppendSlice(res)
	return res
}

type State struct {
	Char       CharState
	Opp        CharState
	TurnNum    int
	GoingFirst int
}

func NewState(in match.GameState) State {
	s := State{
		Char:    NewCharState(in.Character, in.Opponent, in.TurnState),
		Opp:     NewCharState(in.Opponent, in.Character, in.TurnState),
		TurnNum: in.TurnState.TurnNum,
	}

	if in.TurnState.IsGoingFirst {
		s.GoingFirst = 1
	}

	return s
}

func (s State) ToSlice() (res []int) {
	res = s.Char.AppendSlice(res)
	res = s.Opp.AppendSlice(res)
	res = append(res, s.TurnNum, s.GoingFirst)
	return res
}
