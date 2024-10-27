package bot

import (
	"context"
	"math/rand/v2"

	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type RandomBot struct {
	available []int
}

func (b *RandomBot) SendState(ctx context.Context, state match.GameState) error {
	var available []int

	for i, s := range state.Character.Skills() {
		if s.IsAvailable(state.Opponent, state.TurnState) {
			available = append(available, i)
		}
	}

	b.available = available
	return nil
}

func (b *RandomBot) SendError(ctx context.Context) error {
	return nil
}

func (b *RandomBot) SendEnd(ctx context.Context) error {
	return nil
}

func (b RandomBot) RequestSkill(ctx context.Context) (int, error) {
	i := rand.IntN(len(b.available))
	return b.available[i], nil
}
