package bot

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/ShmaykhelDuo/battler/internal/bot/ml/model"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type Bot struct {
	model  *model.Model
	action int
	state  match.GameState
}

func NewBot(model *model.Model) *Bot {
	return &Bot{model: model}
}

func (b *Bot) SendState(ctx context.Context, state match.GameState) error {
	if !state.PlayerTurn {
		return nil
	}

	b.state = state

	res, err := b.model.Predict(state)
	if err != nil {
		return fmt.Errorf("model: %w", err)
	}

	if res < 0 || res > 3 {
		log.Printf("invalid action: %d", res)
		res = getRandomAction(state)
	}

	b.action = res
	return nil
}

func (b *Bot) SendError(ctx context.Context, err error) error {
	log.Printf("got error: %v", err)
	b.action = getRandomAction(b.state)
	return nil
}

func (b *Bot) SendEnd(ctx context.Context) error {
	return nil
}

func (b *Bot) RequestSkill(ctx context.Context) (int, error) {
	return b.action, nil
}

func getRandomAction(state match.GameState) int {
	var available []int

	for i, s := range state.Character.Skills() {
		if s.IsAvailable(state.Character, state.Opponent, state.TurnState) {
			available = append(available, i)
		}
	}
	i := rand.IntN(len(available))
	return available[i]
}
