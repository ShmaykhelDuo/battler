package notification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	model "github.com/ShmaykhelDuo/battler/internal/model/notification"
	service "github.com/ShmaykhelDuo/battler/internal/service/notification"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

type ReceivedConfirmation struct {
	NotificationID uuid.UUID `json:"notification_id"`
}

type Notification struct {
	ID         uuid.UUID       `json:"id"`
	Type       model.Type      `json:"type"`
	Payload    json.RawMessage `json:"payload"`
	CreateTime time.Time       `json:"created_at"`
}

type Conn struct {
	userID        uuid.UUID
	ws            *websocket.Conn
	s             *service.Service
	checkInterval time.Duration
}

func NewConn(userID uuid.UUID, ws *websocket.Conn, s *service.Service) *Conn {
	return &Conn{
		userID:        userID,
		ws:            ws,
		s:             s,
		checkInterval: 5 * time.Second,
	}
}

func (c *Conn) Handle(ctx context.Context) error {
	defer c.ws.Close()

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		err := c.handleSend(ctx)
		if err != nil {
			return fmt.Errorf("handle send: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		err := c.handleReceive(ctx)
		if err != nil {
			return fmt.Errorf("handle send: %w", err)
		}

		return nil
	})

	err := eg.Wait()
	if err != nil {
		var closeErr *websocket.CloseError
		if errors.As(err, &closeErr) && closeErr.Code == websocket.CloseNormalClosure {
			return nil
		}
		return err
	}

	return nil
}

func (c *Conn) handleSend(ctx context.Context) error {
	for {
		err := c.sendNotifications(ctx)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(c.checkInterval):
		}
	}
}

func (c *Conn) handleReceive(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var conf ReceivedConfirmation
			err := c.ws.ReadJSON(&conf)
			if err != nil {
				return fmt.Errorf("read json: %w", err)
			}

			err = c.s.MarkNotificationReceived(ctx, conf.NotificationID)
			if err != nil {
				return fmt.Errorf("mark notification received: %w", err)
			}
		}
	}
}

func (c *Conn) sendNotifications(ctx context.Context) error {
	notifications, err := c.s.PendingNotifications(ctx, c.userID)
	if err != nil {
		return fmt.Errorf("get pending notifications: %w", err)
	}

	for _, n := range notifications {
		dto := Notification{
			ID:         n.ID,
			Type:       n.Type,
			Payload:    n.Payload,
			CreateTime: n.CreateTime,
		}

		err := c.ws.WriteJSON(dto)
		if err != nil {
			return fmt.Errorf("write json: %w", err)
		}
	}

	return nil
}
