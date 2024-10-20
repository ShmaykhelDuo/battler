package main

import (
	"errors"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/bot/minimax"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/bot/ml"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/match"
)

type Bot struct {
	depth     int
	cached    []int
	lastState match.GameState
	dataOut   chan<- []int
}

func NewBot(depth int, dataOut chan<- []int) *Bot {
	return &Bot{depth: depth, dataOut: dataOut}
}

func (b *Bot) SendState(state match.GameState) error {
	if !state.PlayerTurn {
		// fmt.Printf("Skip\n")
		return nil
	}

	b.lastState = state

	if len(b.cached) > 0 {
		// fmt.Printf("SelfHasCache: %v\n", b.cached)
		return nil
	}

	clonedC, clonedOpp := game.Clone(state.Character, state.Opponent)

	skills := clonedC.SkillsPerTurn()
	if state.AsOpp {
		skills = clonedOpp.SkillsPerTurn()
	}
	_, strategy := minimax.MiniMax(clonedC, clonedOpp, state.TurnState, skills, b.depth, state.AsOpp)
	// fmt.Printf("Self Strategy: %v, skills: %d\n", strategy, skills)
	b.dataOut <- append(ml.NewState(state).ToSlice(), strategy[0])
	b.cached = strategy[:skills]
	// fmt.Printf("%v", b.cached)
	// log.Printf("Got Strat: %v\n", b.cached)
	return nil
}

func (b *Bot) SendError() error {
	// log.Printf("\n\n\n\nError!!!!\n\n\n\n")
	// for i, s := range b.lastState.Character.Skills() {
	// 	log.Printf("Skill #%d: %v", i, s.IsAvailable(b.lastState.Opponent, b.lastState.Context))
	// }

	return nil
}

func (b *Bot) SendEnd() error {
	return nil
}

func (b *Bot) RequestSkill() (int, error) {
	// log.Printf("Has Cached: %v\n", b.cached)
	if len(b.cached) < 1 {
		return 0, errors.New("Heck!")
	}

	res := b.cached[0]
	b.cached = b.cached[1:]
	return res, nil
}
