package euphoria_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/euphoria"
	"github.com/stretchr/testify/assert"
)

func TestNewEffectEuphoricSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		amount int
	}{
		{
			name:   "Amount1",
			amount: 1,
		},
		{
			name:   "Amount2",
			amount: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			eff := euphoria.NewEffectEuphoricSource(tt.amount)

			assert.Equal(t, tt.amount, eff.Amount())
		})
	}
}

func TestNewEffectUltimateEarly(t *testing.T) {
	t.Parallel()

	eff := euphoria.NewEffectUltimateEarly()

	assert.Equal(t, 1, eff.Amount())
}

func TestEffectUltimateEarly_Increase(t *testing.T) {
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

			e := euphoria.NewEffectUltimateEarly()

			for range tt.increaseCount {
				e.Increase()
			}

			assert.Equal(t, tt.amount, e.Amount())
		})
	}
}

func TestEffectUltimateEarly_ModifySkillUnlockTurn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		increaseCount int
		skill         game.SkillData
		unlockTurn    int
	}{
		{
			name:          "UltimateInitial",
			increaseCount: 0,
			skill: game.SkillData{
				Desc: game.SkillDescription{
					IsUltimate: true,
				},
				UnlockTurn: 5,
			},
			unlockTurn: 4,
		},
		{
			name:          "UltimateTwice",
			increaseCount: 2,
			skill: game.SkillData{
				Desc: game.SkillDescription{
					IsUltimate: true,
				},
				UnlockTurn: 4,
			},
			unlockTurn: 1,
		},
		{
			name:          "NotUltimateInitial",
			increaseCount: 0,
			skill: game.SkillData{
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
			skill: game.SkillData{
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

			data := game.CharacterData{}
			opp := game.NewCharacter(data)

			e := euphoria.NewEffectUltimateEarly()
			opp.AddEffect(e)

			for range tt.increaseCount {
				e.Increase()
			}

			s := game.NewSkill(opp, tt.skill)

			assert.Equal(t, tt.unlockTurn, s.UnlockTurn())
		})
	}
}

func TestEffectEuphoricHeal_OnTurnEnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		oppData      game.CharacterData
		prevDmg      int
		oppPrevDmg   int
		effs         []game.Effect
		oppEffs      []game.Effect
		turnState    game.TurnState
		hp           int
		oppHP        int
		sourceAmount int
	}{
		{
			name: "Basic",
			oppData: game.CharacterData{
				DefaultHP: 100,
			},
			prevDmg:    25,
			oppPrevDmg: 25,
			effs: []game.Effect{
				euphoria.NewEffectEuphoricSource(20),
			},
			hp:           112,
			oppHP:        95,
			sourceAmount: 11,
		},
		{
			name: "Ending",
			oppData: game.CharacterData{
				DefaultHP: 100,
			},
			prevDmg:    25,
			oppPrevDmg: 25,
			effs: []game.Effect{
				euphoria.NewEffectEuphoricSource(5),
			},
			hp:    97,
			oppHP: 80,
		},
		{
			name: "None",
			oppData: game.CharacterData{
				DefaultHP: 100,
			},
			prevDmg:    25,
			oppPrevDmg: 25,
			hp:         92,
			oppHP:      75,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := game.NewCharacter(euphoria.CharacterEuphoria)
			opp := game.NewCharacter(tt.oppData)

			opp.Damage(c, tt.prevDmg, game.ColourNone)
			c.Damage(opp, tt.oppPrevDmg, game.ColourNone)

			for _, e := range tt.effs {
				c.AddEffect(e)
			}

			eff := euphoria.EffectEuphoricHeal{}
			eff.OnTurnEnd(c, opp, tt.turnState)

			assert.Equal(t, tt.hp, c.HP(), "HP")
			assert.Equal(t, tt.oppHP, opp.HP(), "opponent's HP")
			assert.Equal(t, tt.sourceAmount, euphoricSourceAmount(c), "euphoric source amount")
		})
	}
}
