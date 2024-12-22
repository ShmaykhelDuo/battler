package bot

import (
	"context"

	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type RevAdapter struct {
	last  match.GameState
	skill chan int
	state chan match.GameState
	err   chan error
}

func NewAdapter() *RevAdapter {
	return &RevAdapter{
		skill: make(chan int),
		state: make(chan match.GameState),
		err:   make(chan error),
	}
}

func (a *RevAdapter) GetStateInit() match.GameState {
	return <-a.state
}

func (a *RevAdapter) GetState(skill int) (res match.GameState, err error) {
	a.skill <- skill

	select {
	case res = <-a.state:
		return res, nil
	case err = <-a.err:
		return match.GameState{}, err
	}
}

func (a *RevAdapter) SendState(ctx context.Context, state match.GameState) error {
	a.last = state

	if !state.PlayerTurn {
		return nil
	}

	select {
	case a.state <- state:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (a *RevAdapter) SendError(ctx context.Context, err error) error {
	a.err <- err
	return nil
}

func (a *RevAdapter) SendEnd(ctx context.Context) error {
	select {
	case a.state <- a.last:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (a *RevAdapter) RequestSkill(ctx context.Context) (int, error) {
	select {
	case skill := <-a.skill:
		return skill, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}
