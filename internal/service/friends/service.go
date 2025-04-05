package friends

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ShmaykhelDuo/battler/internal/model/api"
	"github.com/ShmaykhelDuo/battler/internal/model/notification"
	"github.com/ShmaykhelDuo/battler/internal/model/social"
	"github.com/ShmaykhelDuo/battler/internal/pkg/auth"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db"
	"github.com/google/uuid"
)

type FriendRepository interface {
	Friends(ctx context.Context, userID uuid.UUID) ([]social.Profile, error)
	IncomingFriendRequests(ctx context.Context, userID uuid.UUID) ([]social.Profile, error)
	OutgoingFriendRequests(ctx context.Context, userID uuid.UUID) ([]social.Profile, error)
	CreateFriendLink(ctx context.Context, userID uuid.UUID, otherID uuid.UUID) error
	RemoveFriendLink(ctx context.Context, userID uuid.UUID, otherID uuid.UUID) error
	FriendLinkExists(ctx context.Context, userID uuid.UUID, otherID uuid.UUID) (bool, error)
}

type ProfileRepository interface {
	Profile(ctx context.Context, id uuid.UUID) (social.Profile, error)
}

type TransactionManager interface {
	Transact(ctx context.Context, isolation db.TxIsolation, f func(context.Context) error) error
}

type NotificationRepository interface {
	CreateNotification(ctx context.Context, userID uuid.UUID, n notification.Notification) error
}

type Service struct {
	fr FriendRepository
	pr ProfileRepository
	tm TransactionManager
	nr NotificationRepository
}

func NewService(fr FriendRepository, pr ProfileRepository, tm TransactionManager, nr NotificationRepository) *Service {
	return &Service{
		fr: fr,
		tm: tm,
		nr: nr,
		pr: pr,
	}
}

func (s *Service) Friends(ctx context.Context) ([]social.Profile, error) {
	session, err := auth.Session(ctx)
	if err != nil {
		return nil, api.Error{Kind: api.KindUnauthenticated}
	}

	return s.fr.Friends(ctx, session.UserID)
}

func (s *Service) IncomingFriendRequests(ctx context.Context) ([]social.Profile, error) {
	session, err := auth.Session(ctx)
	if err != nil {
		return nil, api.Error{Kind: api.KindUnauthenticated}
	}

	return s.fr.IncomingFriendRequests(ctx, session.UserID)
}

func (s *Service) OutgoingFriendRequests(ctx context.Context) ([]social.Profile, error) {
	session, err := auth.Session(ctx)
	if err != nil {
		return nil, api.Error{Kind: api.KindUnauthenticated}
	}

	return s.fr.OutgoingFriendRequests(ctx, session.UserID)
}

func (s *Service) CreateFriendLink(ctx context.Context, id uuid.UUID) error {
	session, err := auth.Session(ctx)
	if err != nil {
		return api.Error{Kind: api.KindUnauthenticated}
	}

	err = s.tm.Transact(ctx, db.TxIsolationRepeatableRead, func(ctx context.Context) error {
		err := s.fr.CreateFriendLink(ctx, session.UserID, id)
		if err != nil {
			return fmt.Errorf("create friend link: %w", err)
		}

		profile, err := s.pr.Profile(ctx, session.UserID)
		if err != nil {
			return fmt.Errorf("get profile: %w", err)
		}

		hasIncoming, err := s.fr.FriendLinkExists(ctx, id, session.UserID)
		if err != nil {
			return fmt.Errorf("get friend link exists: %w", err)
		}

		var notifType notification.Type
		var payload json.RawMessage

		if hasIncoming {
			notifType = notification.TypeFriendRequestAccepted
			payload, err = json.Marshal(social.FriendRequestAcceptedNotification{
				ID:       profile.ID,
				Username: profile.Username,
			})
			if err != nil {
				return fmt.Errorf("marshal notification: %w", err)
			}
		} else {
			notifType = notification.TypeNewFriendRequest
			payload, err = json.Marshal(social.NewFriendRequestNotification{
				ID:       profile.ID,
				Username: profile.Username,
			})
			if err != nil {
				return fmt.Errorf("marshal notification: %w", err)
			}
		}

		n := notification.Notification{
			ID:         uuid.Must(uuid.NewV7()),
			Type:       notifType,
			Payload:    payload,
			CreateTime: time.Now(),
		}
		err = s.nr.CreateNotification(ctx, id, n)
		if err != nil {
			return fmt.Errorf("create notification: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveFriendLink(ctx context.Context, id uuid.UUID) error {
	session, err := auth.Session(ctx)
	if err != nil {
		return api.Error{Kind: api.KindUnauthenticated}
	}

	return s.fr.RemoveFriendLink(ctx, session.UserID, id)
}
