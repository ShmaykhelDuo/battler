package notification

import (
	"context"
	"fmt"

	model "github.com/ShmaykhelDuo/battler/internal/model/notification"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db"
	"github.com/google/uuid"
)

type NotificationRepository interface {
	PendingNotifications(ctx context.Context, userID uuid.UUID) ([]model.Notification, error)
	MarkNotificationReceived(ctx context.Context, id uuid.UUID) error
}

type TransactionManager interface {
	Transact(ctx context.Context, isolation db.TxIsolation, f func(context.Context) error) error
}

type Service struct {
	repo NotificationRepository
	tm   TransactionManager
}

func NewService(repo NotificationRepository, tm TransactionManager) *Service {
	return &Service{
		repo: repo,
		tm:   tm,
	}
}

func (s *Service) PendingNotifications(ctx context.Context, userID uuid.UUID) ([]model.Notification, error) {
	var res []model.Notification
	err := s.tm.Transact(ctx, db.TxIsolationReadCommitted, func(ctx context.Context) error {
		var err error
		res, err = s.repo.PendingNotifications(ctx, userID)
		if err != nil {
			return fmt.Errorf("get pending notifications: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) MarkNotificationReceived(ctx context.Context, id uuid.UUID) error {
	err := s.repo.MarkNotificationReceived(ctx, id)
	if err != nil {
		return fmt.Errorf("mark notification received: %w", err)
	}

	return nil
}
