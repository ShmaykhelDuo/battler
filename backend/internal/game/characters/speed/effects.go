package speed

import (
	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/common"
)

var EffectDescGreenTokens = game.EffectDescription{
	Name: "Green Tokens",
	Type: game.EffectTypeNumeric,
}

type EffectGreenTokens struct {
	*common.Collectible
}

func NewEffectGreenTokens(number int) EffectGreenTokens {
	return EffectGreenTokens{
		Collectible: common.NewCollectible(number),
	}
}

// Desc returns the effect's description.
func (e EffectGreenTokens) Desc() game.EffectDescription {
	return EffectDescGreenTokens
}

// Clone returns a clone of the effect.
func (e EffectGreenTokens) Clone() game.Effect {
	return NewEffectGreenTokens(e.Amount())
}

var EffectDescBlackTokens = game.EffectDescription{
	Name: "Black Tokens",
	Type: game.EffectTypeNumeric,
}

type EffectBlackTokens struct {
	*common.Collectible
}

func NewEffectBlackTokens(number int) EffectBlackTokens {
	return EffectBlackTokens{
		Collectible: common.NewCollectible(number),
	}
}

// Desc returns the effect's description.
func (e EffectBlackTokens) Desc() game.EffectDescription {
	return EffectDescBlackTokens
}

// Clone returns a clone of the effect.
func (e EffectBlackTokens) Clone() game.Effect {
	return NewEffectBlackTokens(e.Amount())
}

var EffectDescDamageReduced = game.EffectDescription{
	Name: "Damage Reduced",
	Type: game.EffectTypeBuff,
}

// Your opponent's next attack will deal this much less damage.
type EffectDamageReduced struct {
	amount int
	used   bool
}

func NewEffectDamageReduced(amount int) *EffectDamageReduced {
	return &EffectDamageReduced{amount: amount}
}

// Desc returns the effect's description.
func (e *EffectDamageReduced) Desc() game.EffectDescription {
	return EffectDescDamageReduced
}

// Clone returns a clone of the effect.
func (e *EffectDamageReduced) Clone() game.Effect {
	return &EffectDamageReduced{
		amount: e.amount,
		used:   e.used,
	}
}

func (e *EffectDamageReduced) Amount() int {
	return e.amount
}

func (e *EffectDamageReduced) Increase(amount int) {
	e.amount += amount
}

// ModifyTakenDamage returns the modified amount of damage based on provided amount of damage and damage colour.
func (e *EffectDamageReduced) ModifyTakenDamage(dmg int, colour game.Colour) int {
	e.used = true

	return dmg - e.amount
}

// HasExpired reports whether the effect has expired.
func (e *EffectDamageReduced) HasExpired(turnState game.TurnState) bool {
	return e.used
}

var EffectDescDefenceReduced = game.EffectDescription{
	Name: "Defence Reduced",
	Type: game.EffectTypeDebuff,
}

type EffectDefenceReduced struct {
}

// Desc returns the effect's description.
func (e EffectDefenceReduced) Desc() game.EffectDescription {
	return EffectDescDefenceReduced
}

// Clone returns a clone of the effect.
func (e EffectDefenceReduced) Clone() game.Effect {
	return e
}

// ModifyDefences returns the modified defences.
func (e EffectDefenceReduced) ModifyDefences(def map[game.Colour]int) {
	def[game.ColourGreen]--
}

var EffectDescSpedUp = game.EffectDescription{
	Name: "Sped Up",
	Type: game.EffectTypeState,
}

// This turn, you can use two skills but not your ultimate.
type EffectSpedUp struct {
	common.DurationExpirable
}

func NewEffectSpedUp(turnState game.TurnState) EffectSpedUp {
	return EffectSpedUp{
		DurationExpirable: common.NewDurationExpirable(turnState.AddTurns(1, false)),
	}
}

// Desc returns the effect's description.
func (e EffectSpedUp) Desc() game.EffectDescription {
	return EffectDescSpedUp
}

// Clone returns a clone of the effect.
func (e EffectSpedUp) Clone() game.Effect {
	return e
}

// SkillsPerTurn returns a number of tines available for the character to use skills this turn.
func (e EffectSpedUp) SkillsPerTurn() int {
	return 2
}

// IsSkillAvailable reports whether the skill can be used.
func (e EffectSpedUp) IsSkillAvailable(s *game.Skill) bool {
	return !s.Desc().IsUltimate
}
