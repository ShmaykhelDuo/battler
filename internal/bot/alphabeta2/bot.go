package alphabeta2

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type Bot struct {
	depth int
	skill int
}

func NewBot(depth int) *Bot {
	return &Bot{
		depth: depth,
	}
}

func (b *Bot) SendState(ctx context.Context, state match.GameState) error {
	if !state.PlayerTurn {
		return nil
	}

	res, err := Minimax(ctx, state.Clone(), b.depth)
	if err != nil {
		return fmt.Errorf("alphabeta: %w", err)
	}

	b.skill = res.Skill

	return nil
}

func (b *Bot) SendError(ctx context.Context, err error) error {
	return fmt.Errorf("received error from game: %w", err)
}

func (b *Bot) SendEnd(ctx context.Context) error {
	return nil
}

func (b *Bot) RequestSkill(ctx context.Context) (int, error) {
	if b.skill == -1 {
		return 0, errors.New("no skills are cached")
	}

	res := b.skill
	b.skill = -1
	return res, nil
}

func (b *Bot) GivenUp() <-chan any {
	return nil
}
