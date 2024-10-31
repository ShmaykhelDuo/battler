package game_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/gametest"
	"github.com/stretchr/testify/assert"
)

func TestNewCharacter(t *testing.T) {
	t.Parallel()

	data := &game.CharacterData{
		Desc: game.CharacterDescription{
			Name:   "Simple",
			Number: 0,
		},
		DefaultHP: 100,
		Defences: map[game.Colour]int{
			game.ColourWhite: 1,
			game.ColourBlack: -1,
		},
	}
	c := game.NewCharacter(data)

	wantDesc := game.CharacterDescription{
		Name:   "Simple",
		Number: 0,
	}
	assert.Equal(t, wantDesc, c.Desc(), "description")

	wantHP := 100
	assert.Equal(t, wantHP, c.HP(), "hp")
	assert.Equal(t, wantHP, c.MaxHP(), "max hp")

	wantDefences := map[game.Colour]int{
		game.ColourWhite: 1,
		game.ColourBlack: -1,
	}
	assert.Equal(t, wantDefences, c.Defences(), "defences")

	assert.Empty(t, c.Effects(), "effects")

	assert.NotContains(t, c.Skills(), (*game.Skill)(nil), "skills")
}

type defenceModifierEffect struct {
	colour game.Colour
	delta  int
}

// Desc returns the effect's description.
func (e defenceModifierEffect) Desc() game.EffectDescription {
	return game.EffectDescription{}
}

// Clone returns a clone of the effect.
func (e defenceModifierEffect) Clone() game.Effect {
	return e
}

// ModifyDefences returns the modified defences.
func (e defenceModifierEffect) ModifyDefences(def map[game.Colour]int) {
	def[e.colour] += e.delta
}

func TestCharacter_Defences(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		data     *game.CharacterData
		effs     []game.Effect
		defences map[game.Colour]int
	}{
		{
			name: "Basic",
			data: &game.CharacterData{
				Defences: map[game.Colour]int{
					game.ColourWhite: 1,
					game.ColourBlack: -1,
				},
			},
			defences: map[game.Colour]int{
				game.ColourWhite: 1,
				game.ColourBlack: -1,
			},
		},
		{
			name: "EffectModified",
			data: &game.CharacterData{
				Defences: map[game.Colour]int{
					game.ColourWhite: 1,
					game.ColourBlack: -1,
				},
			},
			effs: []game.Effect{
				defenceModifierEffect{
					colour: game.ColourBlack,
					delta:  1,
				},
			},
			defences: map[game.Colour]int{
				game.ColourWhite: 1,
				game.ColourBlack: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(tt.data)

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			assert.Equal(t, tt.defences, c.Defences())
		})
	}
}

type descriptionEffect struct {
	desc game.EffectDescription
}

// Desc returns the effect's description.
func (e descriptionEffect) Desc() game.EffectDescription {
	return e.desc
}

func TestCharacter_LastUsedSkill(t *testing.T) {
	t.Parallel()

	data := &game.CharacterData{
		SkillData: [4]*game.SkillData{
			{Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {}},
			{Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {}},
			{Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {}},
			{Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {}},
		},
	}
	c := game.NewCharacter(data)
	opp := game.NewCharacter(data)

	assert.Nil(t, c.LastUsedSkill(), "before any skills used")

	for i, s := range c.Skills() {
		turnState := game.TurnState{
			TurnNum: i + 1,
		}

		s.Use(c, opp, turnState)

		assert.Same(t, s, c.LastUsedSkill(), "after skill #%d", i+1)
	}
}

type controlEffect struct {
	takenControl bool
}

// Desc returns the effect's description.
func (e controlEffect) Desc() game.EffectDescription {
	return game.EffectDescription{}
}

// Clone returns a clone of the effect.
func (e controlEffect) Clone() game.Effect {
	return e
}

// HasTakenControl reports whether the opponent has taken control over the character.
func (e controlEffect) HasTakenControl() bool {
	return e.takenControl
}

func TestCharacter_IsControlledByOpp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		effs              []game.Effect
		isControlledByOpp bool
	}{
		{
			name:              "Basic",
			isControlledByOpp: false,
		},
		{
			name: "EffectTakenControl",
			effs: []game.Effect{
				controlEffect{takenControl: true},
			},
			isControlledByOpp: true,
		},
		{
			name: "EffectNotTakenControl",
			effs: []game.Effect{
				controlEffect{takenControl: false},
			},
			isControlledByOpp: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(&game.CharacterData{})

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			assert.Equal(t, tt.isControlledByOpp, c.IsControlledByOpp())
		})
	}
}

type skillsPerTurnEffect struct {
	number int
}

// Desc returns the effect's description.
func (e skillsPerTurnEffect) Desc() game.EffectDescription {
	return game.EffectDescription{}
}

// Clone returns a clone of the effect.
func (e skillsPerTurnEffect) Clone() game.Effect {
	return e
}

// SkillsPerTurn returns a number of tines available for the character to use skills this turn.
func (e skillsPerTurnEffect) SkillsPerTurn() int {
	return e.number
}

func TestCharacter_SkillsPerTurn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		effs          []game.Effect
		skillsPerTurn int
	}{
		{
			name:          "Basic",
			skillsPerTurn: 1,
		},
		{
			name: "Modified",
			effs: []game.Effect{
				skillsPerTurnEffect{number: 2},
			},
			skillsPerTurn: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(&game.CharacterData{})

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			assert.Equal(t, tt.skillsPerTurn, c.SkillsPerTurn())
		})
	}
}

func TestCharacter_SetMaxHP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    *game.CharacterData
		prevDmg int
		maxHP   int
		hp      int
	}{
		{
			name: "Basic",
			data: &game.CharacterData{
				DefaultHP: 100,
			},
			prevDmg: 25,
			maxHP:   80,
			hp:      75,
		},
		{
			name: "BelowHP",
			data: &game.CharacterData{
				DefaultHP: 100,
			},
			prevDmg: 15,
			maxHP:   80,
			hp:      80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(tt.data)
			opp := game.NewCharacter(tt.data)

			opp.Damage(c, tt.prevDmg, game.ColourNone)

			c.SetMaxHP(tt.maxHP)

			assert.Equal(t, tt.maxHP, c.MaxHP(), "maximum HP")
			assert.Equal(t, tt.hp, c.HP(), "HP")
		})
	}
}

type testEffect struct{}

// Desc returns the effect's description.
func (e testEffect) Desc() game.EffectDescription {
	return game.EffectDescription{}
}

// Clone returns a clone of the effect.
func (e testEffect) Clone() game.Effect {
	return e
}

type effectFilterEffect struct {
	isAllowed bool
}

// Desc returns the effect's description.
func (e effectFilterEffect) Desc() game.EffectDescription {
	return game.EffectDescription{}
}

// Clone returns a clone of the effect.
func (e effectFilterEffect) Clone() game.Effect {
	return e
}

// IsEffectAllowed reports whether the effect can be applied to a character.
func (e effectFilterEffect) IsEffectAllowed(eff game.Effect) bool {
	return e.isAllowed
}

func TestCharacter_AddEffect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		effs    []game.Effect
		eff     game.Effect
		isAdded bool
	}{
		{
			name:    "Basic",
			eff:     testEffect{},
			isAdded: true,
		},
		{
			name: "Forbidden",
			effs: []game.Effect{
				effectFilterEffect{isAllowed: false},
			},
			eff:     testEffect{},
			isAdded: false,
		},
		{
			name: "NotForbidden",
			effs: []game.Effect{
				effectFilterEffect{isAllowed: true},
			},
			eff:     testEffect{},
			isAdded: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data := &game.CharacterData{}
			c := game.NewCharacter(data)

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			c.AddEffect(tt.eff)

			if tt.isAdded {
				assert.Contains(t, c.Effects(), tt.eff)
			} else {
				assert.NotContains(t, c.Effects(), tt.eff)
			}
		})
	}
}

type dealtDamageModifierEffect struct {
	delta int
}

// Desc returns the effect's description.
func (e dealtDamageModifierEffect) Desc() game.EffectDescription {
	return game.EffectDescription{}
}

// Clone returns a clone of the effect.
func (e dealtDamageModifierEffect) Clone() game.Effect {
	return e
}

// ModifyDealtDamage returns the modified amount of damage based on provided amount of damage and damage colour.
func (e dealtDamageModifierEffect) ModifyDealtDamage(dmg int, colour game.Colour) int {
	return dmg + e.delta
}

type takenDamageModifierEffect struct {
	delta int
}

// Desc returns the effect's description.
func (e takenDamageModifierEffect) Desc() game.EffectDescription {
	return game.EffectDescription{}
}

// Clone returns a clone of the effect.
func (e takenDamageModifierEffect) Clone() game.Effect {
	return e
}

// ModifyTakenDamage returns the modified amount of damage based on provided amount of damage and damage colour.
func (e takenDamageModifierEffect) ModifyTakenDamage(dmg int, colour game.Colour) int {
	return dmg + e.delta
}

func TestCharacter_Damage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    *game.CharacterData
		effs    []game.Effect
		oppEffs []game.Effect
		dmg     int
		colour  game.Colour
		effDmg  int
		hp      int
	}{
		{
			name: "Basic",
			data: &game.CharacterData{
				DefaultHP: 100,
			},
			dmg:    40,
			effDmg: 40,
			hp:     60,
		},
		{
			name: "Kill",
			data: &game.CharacterData{
				DefaultHP: 100,
			},
			dmg:    120,
			effDmg: 120,
			hp:     -20,
		},
		{
			name: "Negative",
			data: &game.CharacterData{
				DefaultHP: 100,
			},
			dmg:    -10,
			effDmg: 0,
			hp:     100,
		},
		{
			name: "PositiveDefence",
			data: &game.CharacterData{
				DefaultHP: 100,
				Defences: map[game.Colour]int{
					game.ColourViolet: 3,
				},
			},
			dmg:    40,
			colour: game.ColourViolet,
			effDmg: 37,
			hp:     63,
		},
		{
			name: "NegativeDefence",
			data: &game.CharacterData{
				DefaultHP: 100,
				Defences: map[game.Colour]int{
					game.ColourViolet: -3,
				},
			},
			dmg:    40,
			colour: game.ColourViolet,
			effDmg: 43,
			hp:     57,
		},
		{
			name: "ModifiedDefence",
			data: &game.CharacterData{
				DefaultHP: 100,
				Defences: map[game.Colour]int{
					game.ColourViolet: 2,
				},
			},
			oppEffs: []game.Effect{
				defenceModifierEffect{
					colour: game.ColourViolet,
					delta:  -1,
				},
			},
			dmg:    40,
			colour: game.ColourViolet,
			effDmg: 39,
			hp:     61,
		},
		{
			name: "AttEffectModification",
			data: &game.CharacterData{
				DefaultHP: 100,
			},
			effs: []game.Effect{
				dealtDamageModifierEffect{delta: 5},
			},
			dmg:    40,
			effDmg: 45,
			hp:     55,
		},
		{
			name: "OppEffectModification",
			data: &game.CharacterData{
				DefaultHP: 100,
			},
			oppEffs: []game.Effect{
				takenDamageModifierEffect{delta: 5},
			},
			dmg:    40,
			effDmg: 45,
			hp:     55,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(tt.data)
			opp := game.NewCharacter(tt.data)

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			for _, e := range tt.oppEffs {
				opp.AddEffect(e)
			}

			effDmg := c.Damage(opp, tt.dmg, tt.colour)
			assert.Equal(t, tt.effDmg, effDmg, "effective damage")

			assert.Equal(t, tt.hp, opp.HP(), "hp after damage")
		})
	}
}

func TestChararcter_Kill(t *testing.T) {
	t.Parallel()

	data := &game.CharacterData{
		DefaultHP: 100,
	}
	c := game.NewCharacter(data)

	c.Kill()

	assert.Equal(t, 0, c.HP())
}

type healFilterEffect struct {
	isAllowed bool
}

// Desc returns the effect's description.
func (e healFilterEffect) Desc() game.EffectDescription {
	return game.EffectDescription{}
}

// Clone returns a clone of the effect.
func (e healFilterEffect) Clone() game.Effect {
	return e
}

// IsHealAllowed reports whether the healing is allowed based on provided amount of healing.
func (e healFilterEffect) IsHealAllowed(heal int) bool {
	return e.isAllowed
}

func TestCharacter_Heal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    *game.CharacterData
		effs    []game.Effect
		prevDmg int
		heal    int
		effHeal int
		hp      int
	}{
		{
			name: "Basic",
			data: &game.CharacterData{
				DefaultHP: 100,
			},
			prevDmg: 40,
			heal:    25,
			effHeal: 25,
			hp:      85,
		},
		{
			name: "Full",
			data: &game.CharacterData{
				DefaultHP: 100,
			},
			prevDmg: 40,
			heal:    55,
			effHeal: 40,
			hp:      100,
		},
		{
			name: "Negative",
			data: &game.CharacterData{
				DefaultHP: 100,
			},
			prevDmg: 40,
			heal:    -10,
			effHeal: 0,
			hp:      60,
		},
		{
			name: "Forbidden",
			data: &game.CharacterData{
				DefaultHP: 100,
			},
			effs: []game.Effect{
				healFilterEffect{isAllowed: false},
			},
			prevDmg: 40,
			heal:    25,
			effHeal: 0,
			hp:      60,
		},
		{
			name: "NotForbidden",
			data: &game.CharacterData{
				DefaultHP: 100,
			},
			effs: []game.Effect{
				healFilterEffect{isAllowed: true},
			},
			prevDmg: 40,
			heal:    25,
			effHeal: 25,
			hp:      85,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(tt.data)
			opp := game.NewCharacter(tt.data)

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			opp.Damage(c, tt.prevDmg, game.ColourNone)

			effHeal := c.Heal(tt.heal)
			assert.Equal(t, tt.effHeal, effHeal, "effective heal")

			assert.Equal(t, tt.hp, c.HP(), "hp after heal")
		})
	}
}

type turnEndHandlerEffect struct {
	gotC, gotOpp *game.Character
	gotturnState game.TurnState
}

// Desc returns the effect's description.
func (e *turnEndHandlerEffect) Desc() game.EffectDescription {
	return game.EffectDescription{}
}

// Clone returns a clone of the effect.
func (e *turnEndHandlerEffect) Clone() game.Effect {
	return e
}

// OnTurnEnd executes the end-of-turn action.
func (e *turnEndHandlerEffect) OnTurnEnd(c, opp *game.Character, turnState game.TurnState) {
	e.gotC = c
	e.gotOpp = opp
	e.gotturnState = turnState
}

func TestCharacter_OnTurnEnd(t *testing.T) {
	t.Parallel()

	t.Run("EffectOnTurnEndExecutes", func(t *testing.T) {
		t.Parallel()

		data := &game.CharacterData{}
		c := game.NewCharacter(data)
		opp := game.NewCharacter(data)

		eff := &turnEndHandlerEffect{}
		c.AddEffect(eff)

		turnState := game.TurnState{
			TurnNum: 4,
		}

		c.OnTurnEnd(opp, turnState)

		assert.Same(t, c, eff.gotC, "character")
		assert.Same(t, opp, eff.gotOpp, "opponent")
		assert.Equal(t, turnState, eff.gotturnState, "game context")
	})

	t.Run("RemovesExpiredEffects", func(t *testing.T) {
		t.Parallel()

		data := &game.CharacterData{}
		c := game.NewCharacter(data)
		opp := game.NewCharacter(data)

		eff := gametest.NewEffectExpirable(true)
		c.AddEffect(eff)

		turnState := game.TurnState{
			TurnNum: 4,
		}

		c.OnTurnEnd(opp, turnState)

		_, found := game.CharacterEffect[*gametest.EffectExpirable](c, gametest.EffectDescExpirable)
		assert.False(t, found, "effect after expiry")
	})
}

func TestCharacter_Clone(t *testing.T) {
	t.Parallel()

	data := &game.CharacterData{}
	c := game.NewCharacter(data)

	clone := c.Clone()

	assert.Equal(t, c, clone, "clone is equal")
	assert.NotSame(t, c, clone, "different pointers")
}

var desc1 = game.EffectDescription{
	Name: "1",
}

type effectType1 struct {
}

// Desc returns the effect's description.
func (e effectType1) Desc() game.EffectDescription {
	return desc1
}

// Clone returns a clone of the effect.
func (e effectType1) Clone() game.Effect {
	return e
}

var desc2 = game.EffectDescription{
	Name: "2",
}

type effectType2 struct {
}

// Desc returns the effect's description.
func (e effectType2) Desc() game.EffectDescription {
	return desc2
}

// Clone returns a clone of the effect.
func (e effectType2) Clone() game.Effect {
	return e
}

var desc3 = game.EffectDescription{
	Name: "3",
}

type effectType3 struct {
}

// Desc returns the effect's description.
func (e effectType3) Desc() game.EffectDescription {
	return desc3
}

// Clone returns a clone of the effect.
func (e effectType3) Clone() game.Effect {
	return e
}

func TestCharacterEffect(t *testing.T) {
	t.Parallel()

	c := game.NewCharacter(&game.CharacterData{})

	eff1 := &effectType1{}
	eff2 := &effectType2{}

	c.AddEffect(eff1)
	c.AddEffect(eff2)

	got1, ok1 := game.CharacterEffect[*effectType1](c, desc1)
	assert.True(t, ok1, "ok 1")
	assert.Same(t, eff1, got1, "same 1")

	got2, ok2 := game.CharacterEffect[*effectType2](c, desc2)
	assert.True(t, ok2, "ok 2")
	assert.Same(t, eff2, got2, "same 2")

	got3, ok3 := game.CharacterEffect[*effectType3](c, desc3)
	assert.False(t, ok3, "ok 3")
	assert.Zero(t, got3, "zero 3")
}
