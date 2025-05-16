package game

import (
	"context"
	"errors"
	"time"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
	"github.com/ShmaykhelDuo/battler/internal/model/api"
	"github.com/google/uuid"
)

type Connection struct {
	userID     uuid.UUID
	state      match.GameState
	stateChan  chan match.GameState
	errorChan  chan error
	endChan    chan MatchPlayerEndResult
	skillChan  chan int
	giveUpChan chan any
}

func NewConnection(userID uuid.UUID) *Connection {
	return &Connection{
		userID:     userID,
		stateChan:  make(chan match.GameState),
		errorChan:  make(chan error),
		endChan:    make(chan MatchPlayerEndResult),
		skillChan:  make(chan int),
		giveUpChan: make(chan any),
	}
}

func (c *Connection) UserID() uuid.UUID {
	return c.userID
}

func (c *Connection) State() <-chan match.GameState {
	return c.stateChan
}

func (c *Connection) SendState(ctx context.Context, state match.GameState) error {
	c.state = state
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

func (c *Connection) End() <-chan MatchPlayerEndResult {
	return c.endChan
}

func (c *Connection) SendEnd(ctx context.Context) error {
	return nil
}

func (c *Connection) SendEndResult(ctx context.Context, res MatchPlayerEndResult) error {
	select {
	case c.endChan <- res:
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
	case <-time.After(2 * time.Minute):
		select {
		case c.giveUpChan <- struct{}{}:
		case <-ctx.Done():
			return 0, ctx.Err()
		}
		<-ctx.Done()
		return 0, ctx.Err()
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func (c *Connection) SendGivenUp() chan<- any {
	return c.giveUpChan
}

func (c *Connection) GivenUp() <-chan any {
	return c.giveUpChan
}

func (c *Connection) LastState() match.GameState {
	return c.state
}
