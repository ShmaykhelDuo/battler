package minimax

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type Bot struct {
	runner Runner
	depth  int
	cached []int
}

func NewBot(runner Runner, depth int) *Bot {
	return &Bot{
		runner: runner,
		depth:  depth,
	}
}

func (b *Bot) SendState(ctx context.Context, state match.GameState) error {
	if !state.PlayerTurn {
		return nil
	}

	if len(b.cached) > 0 {
		return nil
	}

	res, err := b.runner.MiniMax(ctx, state.Clone(), b.depth)
	if err != nil {
		return fmt.Errorf("minimax: %w", err)
	}

	b.cached = res.Strategy[state.TurnState]

	return nil
}

func (b *Bot) SendError(ctx context.Context, err error) error {
	return fmt.Errorf("received error from game: %w", err)
}

func (b *Bot) SendEnd(ctx context.Context) error {
	return nil
}

func (b *Bot) RequestSkill(ctx context.Context) (int, error) {
	if len(b.cached) < 1 {
		return 0, errors.New("no skills are cached")
	}

	res := b.cached[0]
	b.cached = b.cached[1:]
	return res, nil
}

func (b *Bot) GivenUp() <-chan any {
	return nil
}
