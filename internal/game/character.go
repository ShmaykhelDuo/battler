package game

import (
	"maps"
)

// CharacterDescription is a list of constant features of a character.
type CharacterDescription struct {
	Name   string // character's name
	Number int    // character's number
}

// SkillCount is a number of skills provided by a character.
const SkillCount = 4

// CharacterData is a list of features of a character.
type CharacterData struct {
	SkillData [SkillCount]*SkillData // array of character's skills' descriptions
	Desc      CharacterDescription   // character's description
	DefaultHP int                    // character's HP on the beginning of the match
	Defences  map[Colour]int         // map of the character's defences by colour
}

// DefenceModifier modifies a character's defences.
type DefenceModifier interface {
	// ModifyDefences returns the modified defences.
	ModifyDefences(def map[Colour]int)
}

// ControlHandler handles taking control of the character.
type ControlHandler interface {
	// HasTakenControl reports whether the opponent has taken control over the character.
	HasTakenControl() bool
}

// SkillsPerTurnHandler handles changing number of skills per turn.
type SkillsPerTurnHandler interface {
	// SkillsPerTurn returns a number of tines available for the character to use skills this turn.
	SkillsPerTurn() int
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
	OnTurnEnd(c, opp *Character, turnState TurnState)
}

// Expirable represents an effect which can be expired.
type Expirable interface {
	// HasExpired reports whether the effect has expired.
	HasExpired(turnState TurnState) bool
}

// Character is a representation of a character in a match.
type Character struct {
	skills        [SkillCount]*Skill
	data          *CharacterData
	hp            int
	maxHP         int
	effects       map[EffectDescription]Effect
	lastUsedSkill *Skill
}

// NewCharacter returns a new character composed using provided data.
func NewCharacter(data *CharacterData) *Character {
	c := &Character{
		data:    data,
		hp:      data.DefaultHP,
		maxHP:   data.DefaultHP,
		effects: make(map[EffectDescription]Effect),
	}

	for i := range SkillCount {
		c.skills[i] = NewSkill(data.SkillData[i])
	}

	return c
}

// Desc returns the character's description.
func (c *Character) Desc() CharacterDescription {
	return c.data.Desc
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
	defs := maps.Clone(c.data.Defences)

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
	effs := make([]Effect, 0, len(c.effects))
	for _, e := range c.effects {
		effs = append(effs, e)
	}

	return effs
}

// Skills returns an array of skills provided by the character.
func (c *Character) Skills() [SkillCount]*Skill {
	return c.skills
}

// LastUsedSkill returns the skill used last.
func (c *Character) LastUsedSkill() *Skill {
	return c.lastUsedSkill
}

// IsControlledByOpp reports whether the opponent is in control of the character.
func (c *Character) IsControlledByOpp() bool {
	for _, e := range c.effects {
		control, ok := e.(ControlHandler)
		if ok && control.HasTakenControl() {
			return true
		}
	}

	return false
}

// SkillsPerTurn return a number of times a skill can be used this turn.
func (c *Character) SkillsPerTurn() int {
	for _, e := range c.effects {
		h, ok := e.(SkillsPerTurnHandler)
		if ok {
			return h.SkillsPerTurn()
		}
	}

	return 1
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

	desc := eff.Desc()
	c.effects[desc] = eff
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

	dmg -= opp.Defences()[colour]

	if dmg < 0 {
		dmg = 0
	}

	opp.hp -= dmg
	return dmg
}

// Kill immediately kills the character.
func (c *Character) Kill() {
	c.hp = 0
}

func (c *Character) CanHeal(heal int) bool {
	for _, e := range c.effects {
		filter, ok := e.(HealFilter)
		if ok && !filter.IsHealAllowed(heal) {
			return false
		}
	}

	return true
}

// Heal increases the character's HP.
// It returns the actual amount of healing applied to the character.
// Healing can be blocked by any of the effects applied to the character.
func (c *Character) Heal(heal int) int {
	if !c.CanHeal(heal) {
		return 0
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
func (c *Character) OnTurnEnd(opp *Character, turnState TurnState) {
	for _, e := range c.effects {
		h, ok := e.(TurnEndHandler)
		if ok {
			h.OnTurnEnd(c, opp, turnState)
		}
	}

	c.removeExpiredEffects(turnState)
}

func (c *Character) removeExpiredEffects(turnState TurnState) {
	for desc, e := range c.effects {
		exp, ok := e.(Expirable)
		if ok && exp.HasExpired(turnState) {
			delete(c.effects, desc)
		}
	}
}

// Clones returns a deep clone of the character.
func (c *Character) Clone() *Character {
	cloned := &Character{}
	*cloned = *c

	for i, s := range c.skills {
		clonedSkill := &Skill{}
		*clonedSkill = *s
		cloned.skills[i] = clonedSkill

		if c.lastUsedSkill == s {
			cloned.lastUsedSkill = clonedSkill
		}
	}

	cloned.effects = make(map[EffectDescription]Effect, len(c.effects))
	for i, e := range c.effects {
		cloned.effects[i] = e.Clone()
	}
	return cloned
}

// Clones returns a deep clone of the character.
func (c *Character) CloneWithSkill(i int) *Character {
	cloned := &Character{}
	*cloned = *c

	s := c.skills[i]
	clonedSkill := &Skill{}
	*clonedSkill = *s
	cloned.skills[i] = clonedSkill

	if c.lastUsedSkill == s {
		cloned.lastUsedSkill = clonedSkill
	}

	cloned.effects = make(map[EffectDescription]Effect, len(c.effects))
	for i, e := range c.effects {
		cloned.effects[i] = e.Clone()
	}
	return cloned
}

// Clones returns a deep clone of the character.
func (c *Character) CloneWithoutSkills() *Character {
	cloned := &Character{}
	*cloned = *c

	cloned.effects = make(map[EffectDescription]Effect, len(c.effects))
	for i, e := range c.effects {
		cloned.effects[i] = e.Clone()
	}
	return cloned
}

// Effect returns an applied effect with matching description and whether it is found.
// If no such effect is found, zero value of type is returned and found is false.
func CharacterEffect[T Effect](c *Character, desc EffectDescription) (eff T, found bool) {
	e, ok := c.effects[desc]
	if !ok {
		return eff, false
	}

	eff, found = e.(T)
	return
}
