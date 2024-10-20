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

type DefenceStateV2 struct {
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

func NewDefenceStateV2(def map[game.Colour]int) DefenceStateV2 {
	return DefenceStateV2{
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

func (s DefenceStateV2) AppendSlice(res []int) []int {
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

type ColourStateV2 struct {
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

func (s ColourStateV2) AppendSlice(res []int) []int {
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

func NewColourStateV2(c game.Colour) ColourStateV2 {
	res := ColourStateV2{}
	switch c {
	case game.ColourRed:
		res.Red = 1
	case game.ColourOrange:
		res.Orange = 1
	case game.ColourYellow:
		res.Yellow = 1
	case game.ColourGreen:
		res.Green = 1
	case game.ColourCyan:
		res.Cyan = 1
	case game.ColourBlue:
		res.Blue = 1
	case game.ColourViolet:
		res.Violet = 1
	case game.ColourPink:
		res.Pink = 1
	case game.ColourGray:
		res.Gray = 1
	case game.ColourBrown:
		res.Brown = 1
	case game.ColourBlack:
		res.Black = 1
	case game.ColourWhite:
		res.White = 1
	}
	return res
}

type SkillStateV2 struct {
	Colour      ColourStateV2
	Cooldown    int
	UnlockTurn  int
	IsAvailable int
}

func NewSkillStateV2(s *game.Skill, opp *game.Character, turnState game.TurnState) SkillStateV2 {
	res := SkillStateV2{
		Colour:     NewColourStateV2(s.Desc().Colour),
		Cooldown:   s.Cooldown(),
		UnlockTurn: s.UnlockTurn(),
	}

	if s.IsAvailable(opp, turnState) {
		res.IsAvailable = 1
	}

	return res
}

func (s SkillStateV2) AppendSlice(res []int) []int {
	res = s.Colour.AppendSlice(res)
	return append(
		res,
		s.Cooldown,
		s.UnlockTurn,
		s.IsAvailable,
	)
}

type EffectsStateV2 struct {
	StorytellerCannotUseColour   ColourStateV2
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

func NewEffectsStateV2(c *game.Character, turnState game.TurnState) EffectsStateV2 {
	res := EffectsStateV2{}

	for _, e := range c.Effects() {
		switch eff := e.(type) {
		case storyteller.EffectCannotUse:
			res.StorytellerCannotUseColour = NewColourStateV2(eff.Colour())
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

func (s EffectsStateV2) AppendSlice(res []int) []int {
	res = s.StorytellerCannotUseColour.AppendSlice(res)
	return append(
		res,
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

type GirlStateV2 struct {
	Storyteller int
	Z89         int
	Euphoria    int
	Ruby        int
	Speed       int
	Milana      int
	Structure   int
}

func NewGirlStateV2(c *game.Character) GirlStateV2 {
	res := GirlStateV2{}

	switch c.Desc().Number {
	case 1:
		res.Storyteller = 1
	case 8:
		res.Z89 = 1
	case 9:
		res.Euphoria = 1
	case 10:
		res.Ruby = 1
	case 33:
		res.Speed = 1
	case 51:
		res.Milana = 1
	case 119:
		res.Structure = 1
	}

	return res
}

func (s GirlStateV2) AppendSlice(res []int) []int {
	return append(
		res,
		s.Storyteller,
		s.Z89,
		s.Euphoria,
		s.Ruby,
		s.Speed,
		s.Milana,
		s.Structure,
	)
}

type CharStateV2 struct {
	Number        int
	Girl          GirlStateV2
	HP            int
	MaxHP         int
	Defences      DefenceStateV2
	Skills        [4]SkillStateV2
	LastUsedSkill int
	Effects       EffectsStateV2
}

func NewCharStateV2(c *game.Character, opp *game.Character, turnState game.TurnState) CharStateV2 {
	res := CharStateV2{
		Number:   c.Desc().Number,
		Girl:     NewGirlStateV2(c),
		HP:       c.HP(),
		MaxHP:    c.MaxHP(),
		Defences: NewDefenceStateV2(c.Defences()),
		Effects:  NewEffectsStateV2(c, turnState),
	}

	s := c.Skills()

	for i := range 4 {
		res.Skills[i] = NewSkillStateV2(s[i], opp, turnState)
	}

	res.LastUsedSkill = slices.Index(s[:], c.LastUsedSkill())

	return res
}

func (s CharStateV2) AppendSlice(res []int) []int {
	res = append(res, s.Number)
	res = s.Girl.AppendSlice(res)
	res = append(res, s.HP, s.MaxHP)
	res = s.Defences.AppendSlice(res)
	for _, ss := range s.Skills {
		res = ss.AppendSlice(res)
	}
	res = append(res, s.LastUsedSkill)
	res = s.Effects.AppendSlice(res)
	return res
}

type StateV2 struct {
	Char       CharStateV2
	Opp        CharStateV2
	TurnNum    int
	GoingFirst int
}

func NewStateV2(in match.GameState) StateV2 {
	s := StateV2{
		Char:    NewCharStateV2(in.Character, in.Opponent, in.TurnState),
		Opp:     NewCharStateV2(in.Opponent, in.Character, in.TurnState),
		TurnNum: in.TurnState.TurnNum,
	}

	if in.TurnState.IsGoingFirst {
		s.GoingFirst = 1
	}

	return s
}

func (s StateV2) ToSlice() (res []int) {
	res = s.Char.AppendSlice(res)
	res = s.Opp.AppendSlice(res)
	res = append(res, s.TurnNum, s.GoingFirst)
	return res
}
