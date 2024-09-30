package game_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type skillUnlockTurnModifierEffect struct {
	delta int
}

// ModifySkillUnlockTurn returns the modified turn number when skill is to be unlocked.
func (e skillUnlockTurnModifierEffect) ModifySkillUnlockTurn(unlockTurn int) int {
	return unlockTurn + e.delta
}

type skillAvailabilityFilterEffect struct {
	isAvailable bool
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
		Use: func(c *game.Character, opp *game.Character, gameCtx game.Context) {
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
	prevUseCtx  game.Context
	effs        []game.Effect
	gameCtx     game.Context
	isAvailable bool
}{
	{
		name: "Basic",
		data: game.SkillData{},
		gameCtx: game.Context{
			TurnNum: 0,
		},
		isAvailable: true,
	},
	{
		name: "NotUnlocked",
		data: game.SkillData{
			UnlockTurn: 2,
		},
		gameCtx: game.Context{
			TurnNum: 1,
		},
		isAvailable: false,
	},
	{
		name: "JustUnlocked",
		data: game.SkillData{
			UnlockTurn: 2,
		},
		gameCtx: game.Context{
			TurnNum: 2,
		},
		isAvailable: true,
	},
	{
		name: "CooldownNotPassed",
		data: game.SkillData{
			Cooldown: 2,
			Use:      func(c, opp *game.Character, gameCtx game.Context) {},
		},
		wasUsed: true,
		prevUseCtx: game.Context{
			TurnNum: 2,
		},
		gameCtx: game.Context{
			TurnNum: 4,
		},
		isAvailable: false,
	},
	{
		name: "CooldownJustPassed",
		data: game.SkillData{
			Cooldown: 2,
			Use:      func(c, opp *game.Character, gameCtx game.Context) {},
		},
		wasUsed: true,
		prevUseCtx: game.Context{
			TurnNum: 2,
		},
		gameCtx: game.Context{
			TurnNum: 5,
		},
		isAvailable: true,
	},
	{
		name: "ConditionNotFulfilled",
		data: game.SkillData{
			IsAvailable: func(c *game.Character, opp *game.Character, gameCtx game.Context) bool {
				return false
			},
		},
		gameCtx: game.Context{
			TurnNum: 0,
		},
		isAvailable: false,
	},
	{
		name: "ConditionFulfilled",
		data: game.SkillData{
			IsAvailable: func(c *game.Character, opp *game.Character, gameCtx game.Context) bool {
				return true
			},
		},
		gameCtx: game.Context{
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
		gameCtx: game.Context{
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
		gameCtx: game.Context{
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
		gameCtx: game.Context{
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
		gameCtx: game.Context{
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

			assert.Equal(t, tt.isAvailable, s.IsAvailable(opp, tt.gameCtx))
		})
	}
}

func TestSkill_Use(t *testing.T) {
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
			var gotGameCtx game.Context
			data := tt.data
			data.Use = func(c *game.Character, opp *game.Character, gameCtx game.Context) {
				gotC = c
				gotOpp = opp
				gotGameCtx = gameCtx
			}
			s := game.NewSkill(c, data)

			if tt.wasUsed {
				err := s.Use(opp, tt.prevUseCtx)
				require.NoError(t, err)
			}

			err := s.Use(opp, tt.gameCtx)

			if tt.isAvailable {
				require.NoError(t, err)
				assert.Same(t, c, gotC, "character")
				assert.Same(t, opp, gotOpp, "opponent")
				assert.Equal(t, tt.gameCtx, gotGameCtx, "game context")
			} else {
				assert.ErrorIs(t, err, game.ErrSkillNotAvailable)
			}
		})
	}
}
