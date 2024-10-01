package game_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/stretchr/testify/assert"
)

func TestNewCharacter(t *testing.T) {
	t.Parallel()

	data := game.CharacterData{
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
}

// Desc returns the effect's description.
func (e defenceModifierEffect) Desc() game.EffectDescription {
	return game.EffectDescription{}
}

// ModifyDefences returns the modified defences.
func (e defenceModifierEffect) ModifyDefences(def map[game.Colour]int) {
	def[game.ColourBlack] += 1
}

func TestCharacter_Defences(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		data     game.CharacterData
		effs     []game.Effect
		defences map[game.Colour]int
	}{
		{
			name: "Basic",
			data: game.CharacterData{
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
			data: game.CharacterData{
				Defences: map[game.Colour]int{
					game.ColourWhite: 1,
					game.ColourBlack: -1,
				},
			},
			effs: []game.Effect{
				defenceModifierEffect{},
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

func TestCharacter_Effect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		effs []game.Effect
		desc game.EffectDescription
		eff  game.Effect
	}{
		{
			name: "Empty",
			desc: game.EffectDescription{Name: "TestName1"},
			eff:  nil,
		},
		{
			name: "Found",
			effs: []game.Effect{
				descriptionEffect{
					desc: game.EffectDescription{Name: "TestName1"},
				},
				descriptionEffect{
					desc: game.EffectDescription{Name: "TestName2"},
				},
				descriptionEffect{
					desc: game.EffectDescription{Name: "TestName3"},
				},
			},
			desc: game.EffectDescription{Name: "TestName2"},
			eff: descriptionEffect{
				desc: game.EffectDescription{Name: "TestName2"},
			},
		},
		{
			name: "NotFound",
			effs: []game.Effect{
				descriptionEffect{
					desc: game.EffectDescription{Name: "TestName1"},
				},
				descriptionEffect{
					desc: game.EffectDescription{Name: "TestName3"},
				},
			},
			desc: game.EffectDescription{Name: "TestName2"},
			eff:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data := game.CharacterData{}
			c := game.NewCharacter(data)

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			eff := c.Effect(tt.desc)
			assert.Equal(t, tt.eff, eff)
		})
	}
}

func TestCharacter_LastUsedSkill(t *testing.T) {
	t.Parallel()

	data := game.CharacterData{
		SkillData: [4]game.SkillData{
			{Use: func(c *game.Character, opp *game.Character, gameCtx game.Context) {}},
			{Use: func(c *game.Character, opp *game.Character, gameCtx game.Context) {}},
			{Use: func(c *game.Character, opp *game.Character, gameCtx game.Context) {}},
			{Use: func(c *game.Character, opp *game.Character, gameCtx game.Context) {}},
		},
	}
	c := game.NewCharacter(data)
	opp := game.NewCharacter(data)

	assert.Nil(t, c.LastUsedSkill(), "before any skills used")

	for i, s := range c.Skills() {
		gameCtx := game.Context{
			TurnNum: i + 1,
		}

		s.Use(opp, gameCtx)

		assert.Same(t, s, c.LastUsedSkill(), "after skill #%d", i+1)
	}
}

func TestCharacter_SetMaxHP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    game.CharacterData
		prevDmg int
		maxHP   int
		hp      int
	}{
		{
			name: "Basic",
			data: game.CharacterData{
				DefaultHP: 100,
			},
			prevDmg: 25,
			maxHP:   80,
			hp:      75,
		},
		{
			name: "BelowHP",
			data: game.CharacterData{
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

type effectFilterEffect struct {
	isAllowed bool
}

// Desc returns the effect's description.
func (e effectFilterEffect) Desc() game.EffectDescription {
	return game.EffectDescription{}
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

			data := game.CharacterData{}
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

// ModifyTakenDamage returns the modified amount of damage based on provided amount of damage and damage colour.
func (e takenDamageModifierEffect) ModifyTakenDamage(dmg int, colour game.Colour) int {
	return dmg + e.delta
}

func TestCharacter_Damage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    game.CharacterData
		effs    []game.Effect
		oppEffs []game.Effect
		dmg     int
		colour  game.Colour
		effDmg  int
		hp      int
	}{
		{
			name: "Basic",
			data: game.CharacterData{
				DefaultHP: 100,
			},
			dmg:    40,
			effDmg: 40,
			hp:     60,
		},
		{
			name: "Kill",
			data: game.CharacterData{
				DefaultHP: 100,
			},
			dmg:    120,
			effDmg: 100,
			hp:     0,
		},
		{
			name: "Negative",
			data: game.CharacterData{
				DefaultHP: 100,
			},
			dmg:    -10,
			effDmg: 0,
			hp:     100,
		},
		{
			name: "PositiveDefence",
			data: game.CharacterData{
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
			data: game.CharacterData{
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
			name: "AttEffectModification",
			data: game.CharacterData{
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
			data: game.CharacterData{
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

	data := game.CharacterData{
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

// IsHealAllowed reports whether the healing is allowed based on provided amount of healing.
func (e healFilterEffect) IsHealAllowed(heal int) bool {
	return e.isAllowed
}

func TestCharacter_Heal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    game.CharacterData
		effs    []game.Effect
		prevDmg int
		heal    int
		effHeal int
		hp      int
	}{
		{
			name: "Basic",
			data: game.CharacterData{
				DefaultHP: 100,
			},
			prevDmg: 40,
			heal:    25,
			effHeal: 25,
			hp:      85,
		},
		{
			name: "Full",
			data: game.CharacterData{
				DefaultHP: 100,
			},
			prevDmg: 40,
			heal:    55,
			effHeal: 40,
			hp:      100,
		},
		{
			name: "Negative",
			data: game.CharacterData{
				DefaultHP: 100,
			},
			prevDmg: 40,
			heal:    -10,
			effHeal: 0,
			hp:      60,
		},
		{
			name: "Forbidden",
			data: game.CharacterData{
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
			data: game.CharacterData{
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
	gotGameCtx   game.Context
}

// Desc returns the effect's description.
func (e *turnEndHandlerEffect) Desc() game.EffectDescription {
	return game.EffectDescription{}
}

// OnTurnEnd executes the end-of-turn action.
func (e *turnEndHandlerEffect) OnTurnEnd(c, opp *game.Character, gameCtx game.Context) {
	e.gotC = c
	e.gotOpp = opp
	e.gotGameCtx = gameCtx
}

func TestCharacter_OnTurnEnd(t *testing.T) {
	t.Parallel()

	data := game.CharacterData{}
	c := game.NewCharacter(data)
	opp := game.NewCharacter(data)

	eff := &turnEndHandlerEffect{}
	c.AddEffect(eff)

	gameCtx := game.Context{
		TurnNum: 4,
	}

	c.OnTurnEnd(opp, gameCtx)

	assert.Same(t, c, eff.gotC, "character")
	assert.Same(t, opp, eff.gotOpp, "opponent")
	assert.Equal(t, gameCtx, eff.gotGameCtx, "game context")
}
