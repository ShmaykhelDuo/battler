package minimax

import (
	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/match"
)

type Bot struct {
	cached []int
}

func (b *Bot) Skill(state match.GameState) int {
	clonedC, clonedOpp := game.Clone(state.Character, state.Opponent)

	if b.cached == nil {
		skills := clonedC.SkillsPerTurn()
		_, strategy := MiniMax(clonedC, clonedOpp, state.Context, skills, 5, false)
		b.cached = strategy[:skills]
	}

	res := b.cached[0]
	b.cached = b.cached[1:]
	return res
}
