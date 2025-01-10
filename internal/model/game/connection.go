package game

import (
	"context"
	"errors"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
	"github.com/ShmaykhelDuo/battler/internal/model/api"
	"github.com/google/uuid"
)

type Connection struct {
	userID    uuid.UUID
	state     match.GameState
	stateChan chan match.GameState
	errorChan chan error
	endChan   chan any
	skillChan chan int
}

func NewConnection(userID uuid.UUID) *Connection {
	return &Connection{
		userID:    userID,
		stateChan: make(chan match.GameState),
		errorChan: make(chan error),
		endChan:   make(chan any),
		skillChan: make(chan int),
	}
}

func (c *Connection) UserID() uuid.UUID {
	return c.userID
}

func (c *Connection) State() <-chan match.GameState {
	return c.stateChan
}

func (c *Connection) SendState(ctx context.Context, state match.GameState) error {
	select {
	case c.stateChan <- state:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Connection) Error() <-chan error {
	return c.errorChan
}

func (c *Connection) SendError(ctx context.Context, err error) error {
	err = c.handleError(err)

	select {
	case c.errorChan <- err:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Connection) handleError(err error) error {
	if errors.Is(err, game.ErrSkillNotAvailable) {
		return api.Error{
			Kind:    api.KindInvalidArgument,
			Message: "skill is not available",
		}
	}

	return err
}

func (c *Connection) End() <-chan any {
	return c.endChan
}

func (c *Connection) SendEnd(ctx context.Context) error {
	select {
	case c.endChan <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Connection) Skill() chan<- int {
	return c.skillChan
}

func (c *Connection) RequestSkill(ctx context.Context) (int, error) {
	select {
	case skill, ok := <-c.skillChan:
		if !ok {
			return 0, ErrChanClosed
		}
		return skill, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}
