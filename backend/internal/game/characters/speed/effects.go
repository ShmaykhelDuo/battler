package speed

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

var EffectDescGreenTokens = game.EffectDescription{
	Name: "Green Tokens",
	Type: game.EffectTypeNumeric,
}

type EffectGreenTokens struct {
	*Tokens
}

func NewEffectGreenTokens(number int) EffectGreenTokens {
	return EffectGreenTokens{
		Tokens: &Tokens{
			number: number,
		},
	}
}

// Desc returns the effect's description.
func (e EffectGreenTokens) Desc() game.EffectDescription {
	return EffectDescGreenTokens
}

var EffectDescBlackTokens = game.EffectDescription{
	Name: "Black Tokens",
	Type: game.EffectTypeNumeric,
}

type EffectBlackTokens struct {
	*Tokens
}

func NewEffectBlackTokens(number int) EffectBlackTokens {
	return EffectBlackTokens{
		Tokens: &Tokens{
			number: number,
		},
	}
}

// Desc returns the effect's description.
func (e EffectBlackTokens) Desc() game.EffectDescription {
	return EffectDescBlackTokens
}

// Tokens add damage to your Stab.
type Tokens struct {
	number int
}

// Number returns the number of tokens.
func (e *Tokens) Number() int {
	return e.number
}

// Increase increases the number of tokens by 1.
func (e *Tokens) Increase() {
	e.number++
}

var EffectDescDamageReduced = game.EffectDescription{
	Name: "Damage Reduced",
	Type: game.EffectTypeBuff,
}

// Your opponent's next attack will deal this much less damage.
type EffectDamageReduced struct {
	amount int
}

func NewEffectDamageReduced(amount int) *EffectDamageReduced {
	return &EffectDamageReduced{amount: amount}
}

// Desc returns the effect's description.
func (e *EffectDamageReduced) Desc() game.EffectDescription {
	return EffectDescDamageReduced
}

func (e *EffectDamageReduced) Amount() int {
	return e.amount
}

func (e *EffectDamageReduced) Increase(amount int) {
	e.amount += amount
}

// ModifyTakenDamage returns the modified amount of damage based on provided amount of damage and damage colour.
func (e *EffectDamageReduced) ModifyTakenDamage(dmg int, colour game.Colour) int {
	return dmg - e.amount
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
}

// Desc returns the effect's description.
func (e EffectSpedUp) Desc() game.EffectDescription {
	return EffectDescSpedUp
}

// IsSkillAvailable reports whether the skill can be used.
func (e EffectSpedUp) IsSkillAvailable(s *game.Skill) bool {
	return !s.Desc().IsUltimate
}
