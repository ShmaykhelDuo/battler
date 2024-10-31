package z89_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/z89"
	"github.com/stretchr/testify/assert"
)

func TestNewEffectUltimateSlow(t *testing.T) {
	t.Parallel()

	e := z89.NewEffectUltimateSlow()

	assert.Equal(t, 1, e.Amount())
}

func TestEffectUltimateSlow_Increase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		increaseCount int
		amount        int
	}{
		{
			name:          "Once",
			increaseCount: 1,
			amount:        2,
		},
		{
			name:          "Thrice",
			increaseCount: 3,
			amount:        4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := z89.NewEffectUltimateSlow()

			for range tt.increaseCount {
				e.Increase()
			}

			assert.Equal(t, tt.amount, e.Amount())
		})
	}
}

func TestEffectUltimateSlow_ModifySkillUnlockTurn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		increaseCount int
		skill         *game.SkillData
		unlockTurn    int
	}{
		{
			name:          "UltimateInitial",
			increaseCount: 0,
			skill: &game.SkillData{
				Desc: game.SkillDescription{
					IsUltimate: true,
				},
				UnlockTurn: 5,
			},
			unlockTurn: 6,
		},
		{
			name:          "UltimateTwice",
			increaseCount: 2,
			skill: &game.SkillData{
				Desc: game.SkillDescription{
					IsUltimate: true,
				},
				UnlockTurn: 3,
			},
			unlockTurn: 6,
		},
		{
			name:          "NotUltimateInitial",
			increaseCount: 0,
			skill: &game.SkillData{
				Desc: game.SkillDescription{
					IsUltimate: false,
				},
				UnlockTurn: 5,
			},
			unlockTurn: 5,
		},
		{
			name:          "NotUltimateTwice",
			increaseCount: 2,
			skill: &game.SkillData{
				Desc: game.SkillDescription{
					IsUltimate: false,
				},
				UnlockTurn: 3,
			},
			unlockTurn: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data := &game.CharacterData{}
			opp := game.NewCharacter(data)

			e := z89.NewEffectUltimateSlow()
			opp.AddEffect(e)

			for range tt.increaseCount {
				e.Increase()
			}

			s := game.NewSkill(tt.skill)

			assert.Equal(t, tt.unlockTurn, s.UnlockTurn(opp))
		})
	}
}
