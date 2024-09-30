package game

import "maps"

// CharacterDescription is a list of constant features of a character.
type CharacterDescription struct {
	Name   string // character's name
	Number int    // character's number
}

// SkillCount is a number of skills provided by a character.
const SkillCount = 4

// CharacterData is a list of features of a character.
type CharacterData struct {
	Desc      CharacterDescription  // character's description
	DefaultHP int                   // character's HP on the beginning of the match
	Defences  map[Colour]int        // map of the character's defences by colour
	SkillData [SkillCount]SkillData // array of character's skills' descriptions
}

// DefenceModifier modifies a character's defences.
type DefenceModifier interface {
	// ModifyDefences returns the modified defences.
	ModifyDefences(def map[Colour]int)
}

// EffectFilter filters effects allowed to be applied to a character.
type EffectFilter interface {
	// IsEffectAllowed reports whether the effect can be applied to a character.
	IsEffectAllowed(e Effect) bool
}

// DealtDamageModifier modifies the amount of damage for a character to deal.
type DealtDamageModifier interface {
	// ModifyDealtDamage returns the modified amount of damage based on provided amount of damage and damage colour.
	ModifyDealtDamage(dmg int, colour Colour) int
}

// TakenDamageModifier modifies the amount of damage for a character to take.
type TakenDamageModifier interface {
	// ModifyTakenDamage returns the modified amount of damage based on provided amount of damage and damage colour.
	ModifyTakenDamage(dmg int, colour Colour) int
}

// HealFilter filters healing of a character.
type HealFilter interface {
	// IsHealAllowed reports whether the healing is allowed based on provided amount of healing.
	IsHealAllowed(heal int) bool
}

// TurnEndHandler handles the end-of-turn action.
type TurnEndHandler interface {
	// OnTurnEnd executes the end-of-turn action.
	OnTurnEnd(c, opp *Character, gameCtx Context)
}

// Character is a representation of a character in a match.
type Character struct {
	desc          CharacterDescription
	hp            int
	maxHP         int
	defences      map[Colour]int
	skills        [SkillCount]*Skill
	effects       []Effect
	lastUsedSkill *Skill
}

// NewCharacter returns a new character composed using provided data.
func NewCharacter(data CharacterData) *Character {
	c := &Character{
		desc:     data.Desc,
		hp:       data.DefaultHP,
		maxHP:    data.DefaultHP,
		defences: data.Defences,
	}

	for i := range SkillCount {
		c.skills[i] = NewSkill(c, data.SkillData[i])
	}

	return c
}

// Desc returns the character's description.
func (c *Character) Desc() CharacterDescription {
	return c.desc
}

// HP returns the character's current HP.
func (c *Character) HP() int {
	return c.hp
}

// MaxHP returns the character's current maximum HP.
func (c *Character) MaxHP() int {
	return c.maxHP
}

// Defences returns a map of the character's defences for each colour.
// Defence is a modifier which is subtracted from damage of specific colour.
func (c *Character) Defences() map[Colour]int {
	defs := maps.Clone(c.defences)

	for _, e := range c.effects {
		mod, ok := e.(DefenceModifier)
		if ok {
			mod.ModifyDefences(defs)
		}
	}

	return defs
}

// Effects returns a slice of effects applied to the character.
func (c *Character) Effects() []Effect {
	return c.effects
}

// Skills returns an array of skills provided by the character.
func (c *Character) Skills() [SkillCount]*Skill {
	return c.skills
}

// LastUsedSkill returns the skill used last.
func (c *Character) LastUsedSkill() *Skill {
	return c.lastUsedSkill
}

// SetMaxHP sets the character's maximum HP.
// If the new value is less than the character's HP, HP is decreased to match maximum HP.
func (c *Character) SetMaxHP(maxHP int) {
	c.maxHP = maxHP

	if c.hp > maxHP {
		c.hp = maxHP
	}
}

// AddEffect applies an effect to the character.
// The effect can be blocked by any of effects already applied to the character.
func (c *Character) AddEffect(eff Effect) {
	for _, e := range c.effects {
		filter, ok := e.(EffectFilter)
		if ok && !filter.IsEffectAllowed(eff) {
			return
		}
	}

	c.effects = append(c.effects, eff)
}

// Damage decreases the opponent's HP.
// It returns the actual amount of damage taken by the opponent.
// The actual amount of damage is affected by the character's and opponent's effects and opponent's defences.
func (c *Character) Damage(opp *Character, dmg int, colour Colour) int {
	for _, e := range c.effects {
		mod, ok := e.(DealtDamageModifier)
		if ok {
			dmg = mod.ModifyDealtDamage(dmg, colour)
		}
	}

	for _, e := range opp.effects {
		mod, ok := e.(TakenDamageModifier)
		if ok {
			dmg = mod.ModifyTakenDamage(dmg, colour)
		}
	}

	dmg -= opp.defences[colour]

	if dmg < 0 {
		dmg = 0
	}

	if dmg > opp.hp {
		dmg = opp.hp
	}

	opp.hp -= dmg
	return dmg
}

// Kill immediately kills the character.
func (c *Character) Kill() {
	c.hp = 0
}

// Heal increases the character's HP.
// It returns the actual amount of healing applied to the character.
// Healing can be blocked by any of the effects applied to the character.
func (c *Character) Heal(heal int) int {
	for _, e := range c.effects {
		filter, ok := e.(HealFilter)
		if ok && !filter.IsHealAllowed(heal) {
			return 0
		}
	}

	if heal < 0 {
		heal = 0
	}

	if c.hp+heal > c.maxHP {
		heal = c.maxHP - c.hp
	}

	c.hp += heal
	return heal
}

// OnTurnEnd triggers all the end-of-turn actions provided by effects applied to the character.
func (c *Character) OnTurnEnd(opp *Character, gameCtx Context) {
	for _, e := range c.effects {
		h, ok := e.(TurnEndHandler)
		if ok {
			h.OnTurnEnd(c, opp, gameCtx)
		}
	}
}
