package game

import "errors"

// SkillDescription is a list of constant features of a skill.
type SkillDescription struct {
	Name       string // skill's name
	Colour     Colour // skill's colour
	IsUltimate bool   // whether the skill is the ultimate of the character
}

// SkillData is a list of features of a skill.
type SkillData struct {
	Desc          SkillDescription                                  // skill's description
	Cooldown      int                                               // skill's cooldown, 0 if absent
	UnlockTurn    int                                               // the turn number when skill is to be unlocked, 0 if always unlocked
	Use           func(c, opp *Character, turnState TurnState)      // action to be executed on use of skill
	IsAvailable   func(c, opp *Character, turnState TurnState) bool // optional filter, reports whether the skill is available
	IsAppropriate func(c, opp *Character, turnState TurnState) bool
}

// SkillAvailabilityFilter filters skills to be used by a character.
type SkillAvailabilityFilter interface {
	// IsSkillAvailable reports whether the skill can be used.
	IsSkillAvailable(s *Skill) bool
}

// SkillUnlockTurnModifier modified the turn number when skill is to be unlocked.
type SkillUnlockTurnModifier interface {
	// ModifySkillUnlockTurn returns the modified turn number when skill is to be unlocked.
	ModifySkillUnlockTurn(s *Skill, unlockTurn int) int
}

// ErrSkillNotAvailable is an error which is returned when skill is not available.
var ErrSkillNotAvailable = errors.New("skill is not available")

// Skill is a representation of a skill of a character in a match.
type Skill struct {
	desc            SkillDescription
	cooldown        int
	unlockTurn      int
	prevUsedTurn    int
	useFunc         func(c, opp *Character, turnState TurnState)
	availableFunc   func(c, opp *Character, turnState TurnState) bool
	appropriateFunc func(c, opp *Character, turnState TurnState) bool
	c               *Character
}

// NewSkill returns a new skill composed using provided character and data.
func NewSkill(c *Character, data SkillData) *Skill {
	return &Skill{
		desc:            data.Desc,
		cooldown:        data.Cooldown,
		unlockTurn:      data.UnlockTurn,
		useFunc:         data.Use,
		availableFunc:   data.IsAvailable,
		appropriateFunc: data.IsAppropriate,
		c:               c,
	}
}

// Desc returns the skill's description.
func (s *Skill) Desc() SkillDescription {
	return s.desc
}

// Cooldown returns the skill's cooldown.
// Returns 0 when no cooldown is applicable to the skill.
func (s *Skill) Cooldown() int {
	return s.cooldown
}

// UnlockTurn returns the turn number when the skill is to be unlocked.
// Returns 0 when skill is unlocked from the beginning of the game.
// The value can be modified by the character's effects.
func (s *Skill) UnlockTurn() int {
	turn := s.unlockTurn

	if turn == 0 {
		return 0
	}

	for _, e := range s.c.Effects() {
		mod, ok := e.(SkillUnlockTurnModifier)
		if ok {
			turn = mod.ModifySkillUnlockTurn(s, turn)
		}
	}

	if turn < MinTurnNumber {
		return MinTurnNumber
	}

	if turn > MaxTurnNumber {
		return MaxTurnNumber
	}

	return turn
}

// IsAvailable reports whether the skill is available.
// Unlock turn number, cooldown and effects are taken into account.
func (s *Skill) IsAvailable(opp *Character, turnState TurnState) bool {
	if turnState.TurnNum < s.UnlockTurn() {
		return false
	}

	if s.prevUsedTurn != 0 && turnState.TurnNum <= s.prevUsedTurn+s.cooldown {
		return false
	}

	if s.availableFunc != nil && !s.availableFunc(s.c, opp, turnState) {
		return false
	}

	for _, e := range s.c.Effects() {
		filter, ok := e.(SkillAvailabilityFilter)
		if ok && !filter.IsSkillAvailable(s) {
			return false
		}
	}

	return true
}

// Use executes the skill's action.
// Returns [ErrSkillNotAvailable] when the skill is not available for use.
func (s *Skill) Use(opp *Character, turnState TurnState) error {
	if !s.IsAvailable(opp, turnState) {
		return ErrSkillNotAvailable
	}

	s.useFunc(s.c, opp, turnState)

	s.c.removeExpiredEffects(turnState)
	opp.removeExpiredEffects(turnState)

	s.prevUsedTurn = turnState.TurnNum
	s.c.lastUsedSkill = s

	return nil
}

func (s *Skill) IsAppropriate(opp *Character, turnState TurnState) bool {
	if !s.IsAvailable(opp, turnState) {
		return false
	}

	return s.appropriateFunc == nil || s.appropriateFunc(s.c, opp, turnState)
}
