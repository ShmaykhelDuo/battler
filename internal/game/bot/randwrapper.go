package bot

import (
	"math/rand/v2"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type RandomWrapperBot struct {
	bot       match.Player
	p         float64
	lastCtx   game.TurnState
	isRand    bool
	available []int
}

func NewRandomWrapperBot(bot match.Player, p float64) *RandomWrapperBot {
	return &RandomWrapperBot{
		bot: bot,
		p:   p,
	}
}

func (b *RandomWrapperBot) SendState(state match.GameState) error {
	if state.TurnState != b.lastCtx {
		b.isRand = rand.Float64() < b.p
		b.lastCtx = state.TurnState
	}

	if !b.isRand {
		return b.bot.SendState(state)
	}

	var available []int

	for i, s := range state.Character.Skills() {
		if s.IsAvailable(state.Opponent, state.TurnState) {
			available = append(available, i)
		}
	}

	b.available = available
	return nil
}

func (b *RandomWrapperBot) SendError() error {
	if b.isRand {
		return nil
	}
	return b.bot.SendError()
}

func (b *RandomWrapperBot) SendEnd() error {
	return b.bot.SendEnd()
}

func (b RandomWrapperBot) RequestSkill() (int, error) {
	if b.isRand {
		i := rand.IntN(len(b.available))
		return b.available[i], nil
	}

	return b.bot.RequestSkill()
}
