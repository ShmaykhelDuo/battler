package game_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/gametest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type skillUnlockTurnModifierEffect struct {
	delta int
}

// Desc returns the effect's description.
func (e skillUnlockTurnModifierEffect) Desc() game.EffectDescription {
	return game.EffectDescription{}
}

// Clone returns a clone of the effect.
func (e skillUnlockTurnModifierEffect) Clone() game.Effect {
	return e
}

// ModifySkillUnlockTurn returns the modified turn number when skill is to be unlocked.
func (e skillUnlockTurnModifierEffect) ModifySkillUnlockTurn(s *game.Skill, unlockTurn int) int {
	return unlockTurn + e.delta
}

type skillAvailabilityFilterEffect struct {
	isAvailable bool
}

// Desc returns the effect's description.
func (e skillAvailabilityFilterEffect) Desc() game.EffectDescription {
	return game.EffectDescription{}
}

// Clone returns a clone of the effect.
func (e skillAvailabilityFilterEffect) Clone() game.Effect {
	return e
}

// IsSkillAvailable reports whether the skill can be used.
func (e skillAvailabilityFilterEffect) IsSkillAvailable(s *game.Skill) bool {
	return e.isAvailable
}

func TestNewSkill(t *testing.T) {
	t.Parallel()

	charData := game.CharacterData{}
	c := game.NewCharacter(charData)

	skillData := game.SkillData{
		Desc: game.SkillDescription{
			Name:       "Simple",
			IsUltimate: false,
			Colour:     game.ColourWhite,
		},
		Cooldown:   2,
		UnlockTurn: 5,
		Use: func(c *game.Character, opp *game.Character, turnState game.TurnState) {
		},
	}
	s := game.NewSkill(c, skillData)

	wantDesc := game.SkillDescription{
		Name:       "Simple",
		IsUltimate: false,
		Colour:     game.ColourWhite,
	}
	assert.Equal(t, wantDesc, s.Desc(), "description")

	wantCooldown := 2
	assert.Equal(t, wantCooldown, s.Cooldown(), "cooldown")

	wantUnlockTurn := 5
	assert.Equal(t, wantUnlockTurn, s.UnlockTurn(), "unlock turn")
}

func TestSkill_UnlockTurn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		data       game.SkillData
		effs       []game.Effect
		unlockTurn int
	}{
		{
			name: "Absent",
			data: game.SkillData{
				UnlockTurn: 0,
			},
			effs:       []game.Effect{},
			unlockTurn: 0,
		},
		{
			name: "Present",
			data: game.SkillData{
				UnlockTurn: 2,
			},
			effs:       []game.Effect{},
			unlockTurn: 2,
		},
		{
			name: "AbsentDelay",
			data: game.SkillData{
				UnlockTurn: 0,
			},
			effs: []game.Effect{
				skillUnlockTurnModifierEffect{3},
			},
			unlockTurn: 0,
		},
		{
			name: "PresentDelayInBounds",
			data: game.SkillData{
				UnlockTurn: 2,
			},
			effs: []game.Effect{
				skillUnlockTurnModifierEffect{3},
			},
			unlockTurn: 5,
		},
		{
			name: "PresentDelayOutOfBounds",
			data: game.SkillData{
				UnlockTurn: 8,
			},
			effs: []game.Effect{
				skillUnlockTurnModifierEffect{3},
			},
			unlockTurn: 10,
		},
		{
			name: "AbsentSpeedUp",
			data: game.SkillData{
				UnlockTurn: 0,
			},
			effs: []game.Effect{
				skillUnlockTurnModifierEffect{-3},
			},
			unlockTurn: 0,
		},
		{
			name: "PresentSpeedUpInBounds",
			data: game.SkillData{
				UnlockTurn: 5,
			},
			effs: []game.Effect{
				skillUnlockTurnModifierEffect{-3},
			},
			unlockTurn: 2,
		},
		{
			name: "PresentSpeedUpToOutOfBounds",
			data: game.SkillData{
				UnlockTurn: 2,
			},
			effs: []game.Effect{
				skillUnlockTurnModifierEffect{-3},
			},
			unlockTurn: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			charData := game.CharacterData{}
			c := game.NewCharacter(charData)

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			s := game.NewSkill(c, tt.data)

			assert.Equal(t, tt.unlockTurn, s.UnlockTurn())
		})
	}
}

var skillAvailabilityTests = []struct {
	name        string
	data        game.SkillData
	wasUsed     bool
	prevUseCtx  game.TurnState
	effs        []game.Effect
	turnState   game.TurnState
	isAvailable bool
}{
	{
		name: "Basic",
		data: game.SkillData{},
		turnState: game.TurnState{
			TurnNum: 0,
		},
		isAvailable: true,
	},
	{
		name: "NotUnlocked",
		data: game.SkillData{
			UnlockTurn: 2,
		},
		turnState: game.TurnState{
			TurnNum: 1,
		},
		isAvailable: false,
	},
	{
		name: "JustUnlocked",
		data: game.SkillData{
			UnlockTurn: 2,
		},
		turnState: game.TurnState{
			TurnNum: 2,
		},
		isAvailable: true,
	},
	{
		name: "CooldownNotPassed",
		data: game.SkillData{
			Cooldown: 2,
			Use:      func(c, opp *game.Character, turnState game.TurnState) {},
		},
		wasUsed: true,
		prevUseCtx: game.TurnState{
			TurnNum: 2,
		},
		turnState: game.TurnState{
			TurnNum: 4,
		},
		isAvailable: false,
	},
	{
		name: "CooldownJustPassed",
		data: game.SkillData{
			Cooldown: 2,
			Use:      func(c, opp *game.Character, turnState game.TurnState) {},
		},
		wasUsed: true,
		prevUseCtx: game.TurnState{
			TurnNum: 2,
		},
		turnState: game.TurnState{
			TurnNum: 5,
		},
		isAvailable: true,
	},
	{
		name: "ConditionNotFulfilled",
		data: game.SkillData{
			IsAvailable: func(c *game.Character, opp *game.Character, turnState game.TurnState) bool {
				return false
			},
		},
		turnState: game.TurnState{
			TurnNum: 0,
		},
		isAvailable: false,
	},
	{
		name: "ConditionFulfilled",
		data: game.SkillData{
			IsAvailable: func(c *game.Character, opp *game.Character, turnState game.TurnState) bool {
				return true
			},
		},
		turnState: game.TurnState{
			TurnNum: 0,
		},
		isAvailable: true,
	},
	{
		name: "BlockedByEffect",
		data: game.SkillData{},
		effs: []game.Effect{
			skillAvailabilityFilterEffect{isAvailable: false},
		},
		turnState: game.TurnState{
			TurnNum: 0,
		},
		isAvailable: false,
	},
	{
		name: "NotBlockedByEffect",
		data: game.SkillData{},
		effs: []game.Effect{
			skillAvailabilityFilterEffect{isAvailable: true},
		},
		turnState: game.TurnState{
			TurnNum: 0,
		},
		isAvailable: true,
	},
	{
		name: "NotUnlockedWithTurnModification",
		data: game.SkillData{
			UnlockTurn: 2,
		},
		effs: []game.Effect{
			skillUnlockTurnModifierEffect{delta: 2},
		},
		turnState: game.TurnState{
			TurnNum: 3,
		},
		isAvailable: false,
	},
	{
		name: "JustUnlockedWithTurnModification",
		data: game.SkillData{
			UnlockTurn: 2,
		},
		effs: []game.Effect{
			skillUnlockTurnModifierEffect{delta: 2},
		},
		turnState: game.TurnState{
			TurnNum: 4,
		},
		isAvailable: true,
	},
}

func TestSkill_IsAvailable(t *testing.T) {
	t.Parallel()

	tests := skillAvailabilityTests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			charData := game.CharacterData{}
			c := game.NewCharacter(charData)
			opp := game.NewCharacter((charData))

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			s := game.NewSkill(c, tt.data)

			if tt.wasUsed {
				err := s.Use(opp, tt.prevUseCtx)
				require.NoError(t, err)
			}

			assert.Equal(t, tt.isAvailable, s.IsAvailable(opp, tt.turnState))
		})
	}
}

func TestSkill_Use(t *testing.T) {
	t.Parallel()

	t.Run("RunsWhenAvailable", func(t *testing.T) {
		t.Parallel()

		tests := skillAvailabilityTests
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				charData := game.CharacterData{}
				c := game.NewCharacter(charData)
				opp := game.NewCharacter((charData))

				for _, e := range tt.effs {
					c.AddEffect(e)
				}

				var gotC, gotOpp *game.Character
				var gotturnState game.TurnState
				data := tt.data
				data.Use = func(c *game.Character, opp *game.Character, turnState game.TurnState) {
					gotC = c
					gotOpp = opp
					gotturnState = turnState
				}
				s := game.NewSkill(c, data)

				if tt.wasUsed {
					err := s.Use(opp, tt.prevUseCtx)
					require.NoError(t, err)
				}

				err := s.Use(opp, tt.turnState)

				if tt.isAvailable {
					require.NoError(t, err)
					assert.Same(t, c, gotC, "character")
					assert.Same(t, opp, gotOpp, "opponent")
					assert.Equal(t, tt.turnState, gotturnState, "game context")
				} else {
					assert.ErrorIs(t, err, game.ErrSkillNotAvailable)
				}
			})
		}
	})

	t.Run("RemovesExpiredEffects", func(t *testing.T) {
		t.Parallel()

		charData := game.CharacterData{}
		c := game.NewCharacter(charData)
		opp := game.NewCharacter((charData))

		eff := gametest.NewEffectExpirable(false)
		c.AddEffect(eff)

		oppEff := gametest.NewEffectExpirable(false)
		opp.AddEffect(oppEff)

		data := game.SkillData{
			Use: func(c, opp *game.Character, turnState game.TurnState) {
				eff.Expire()
				oppEff.Expire()
			},
		}
		s := game.NewSkill(c, data)

		s.Use(opp, game.TurnState{})

		_, found := game.CharacterEffect[*gametest.EffectExpirable](c, gametest.EffectDescExpirable)
		assert.False(t, found, "effect after expiry")

		_, oppFound := game.CharacterEffect[*gametest.EffectExpirable](opp, gametest.EffectDescExpirable)
		assert.False(t, oppFound, "opponent's effect after expiry")
	})
}
