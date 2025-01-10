package ws

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ShmaykhelDuo/battler/internal/app/game/ws/gamestate"
	"github.com/ShmaykhelDuo/battler/internal/model/api"
	model "github.com/ShmaykhelDuo/battler/internal/model/game"
	service "github.com/ShmaykhelDuo/battler/internal/service/match"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

type Conn struct {
	userID uuid.UUID
	ws     *websocket.Conn
	s      *service.Service
	conn   *model.Connection
}

func NewConn(userID uuid.UUID, ws *websocket.Conn, s *service.Service) *Conn {
	return &Conn{
		userID: userID,
		ws:     ws,
		s:      s,
	}
}

func (c *Conn) Handle(ctx context.Context) error {
	defer c.ws.Close()

	err := c.initiate(ctx)
	if err != nil {
		return fmt.Errorf("initiate: %w", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return c.handleSend(ctx)
	})

	eg.Go(func() error {
		return c.handleReceive(ctx)
	})

	return eg.Wait()
}

func (c *Conn) initiate(ctx context.Context) error {
	msg, err := c.receiveMessage()
	if err != nil {
		return err
	}

	switch m := msg.(type) {
	case MessageMatchRequest:
		c.conn, err = c.s.ConnectToNewMatch(ctx, c.userID, m.MainCharacter, m.SecondaryCharacter)
		if err != nil {
			return c.handleError(fmt.Errorf("new match: %w", err))
		}
	case MessageMatchReconnect:
		c.conn, err = c.s.ReconnectToMatch(ctx, c.userID)
		if err != nil {
			return c.handleError(fmt.Errorf("match reconnect: %w", err))
		}
	default:
		return c.handleError(api.Error{
			Kind:    api.KindInvalidRequest,
			Message: "illegal initial request type",
		})
	}

	return nil
}

func (c *Conn) handleSend(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case state, ok := <-c.conn.State():
			if !ok {
				return model.ErrChanClosed
			}
			msg := MessageGameState{
				State: gamestate.NewGameState(state),
			}
			err := c.sendMessage(msg)
			if err != nil {
				return err
			}
		case retErr, ok := <-c.conn.Error():
			if !ok {
				return model.ErrChanClosed
			}
			err := c.handleError(retErr)
			if err != nil {
				return err
			}
		case <-c.conn.End():
			msg := MessageGameEnd{}
			err := c.sendMessage(msg)
			if err != nil {
				return err
			}
		}
	}
}

func (c *Conn) handleReceive(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m, err := c.receiveMessage()
			if err != nil {
				return err
			}
			switch msg := m.(type) {
			case MessageMove:
				select {
				case <-ctx.Done():
					return ctx.Err()
				case c.conn.Skill() <- msg.Move:
				}
			default:
				err := c.handleError(api.Error{
					Kind:    api.KindInvalidRequest,
					Message: "illegal request type",
				})
				if err != nil {
					return err
				}
			}
		}
	}
}

func (c *Conn) handleError(err error) error {
	var apiError api.Error
	if !errors.As(err, &apiError) {
		slog.Error("unhandled error", "err", err)

		apiError = api.Error{
			Kind: api.KindInternal,
		}
	}

	msg := MessageError{
		ID:      apiError.Kind.ID,
		Kind:    apiError.Kind.Description,
		Message: apiError.Message,
	}
	return c.sendMessage(msg)
}

func (c *Conn) sendMessage(payload any) error {
	msg, err := NewMessage(payload)
	if err != nil {
		return err
	}

	err = c.ws.WriteJSON(msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *Conn) receiveMessage() (any, error) {
	var msg Message
	err := c.ws.ReadJSON(&msg)
	if err != nil {
		return nil, err
	}

	res, err := msg.UnmarshalPayload()
	if err != nil {
		return nil, err
	}

	return res, nil
}
