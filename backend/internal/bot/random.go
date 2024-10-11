package bot

import (
	"math/rand/v2"

	"github.com/ShmaykhelDuo/battler/backend/internal/game/match"
)

type RandomBot struct {
}

func (b RandomBot) Skill(state match.GameState) int {
	var available []int

	for i, s := range state.Character.Skills() {
		if s.IsAvailable(state.Opponent, state.Context) {
			available = append(available, i)
		}
	}

	i := rand.IntN(len(available))
	return available[i]
}
