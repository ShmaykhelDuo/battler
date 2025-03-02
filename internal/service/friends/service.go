package friends

import (
	"context"

	"github.com/ShmaykhelDuo/battler/internal/model/api"
	"github.com/ShmaykhelDuo/battler/internal/model/social"
	"github.com/ShmaykhelDuo/battler/internal/pkg/auth"
	"github.com/google/uuid"
)

type FriendRepository interface {
	Friends(ctx context.Context, userID uuid.UUID) ([]social.Profile, error)
	IncomingFriendRequests(ctx context.Context, userID uuid.UUID) ([]social.Profile, error)
	OutgoingFriendRequests(ctx context.Context, userID uuid.UUID) ([]social.Profile, error)
	CreateFriendLink(ctx context.Context, userID uuid.UUID, otherID uuid.UUID) error
	RemoveFriendLink(ctx context.Context, userID uuid.UUID, otherID uuid.UUID) error
}

type Service struct {
	repo FriendRepository
}

func NewService(repo FriendRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Friends(ctx context.Context) ([]social.Profile, error) {
	session, err := auth.Session(ctx)
	if err != nil {
		return nil, api.Error{Kind: api.KindUnauthenticated}
	}

	return s.repo.Friends(ctx, session.UserID)
}

func (s *Service) IncomingFriendRequests(ctx context.Context) ([]social.Profile, error) {
	session, err := auth.Session(ctx)
	if err != nil {
		return nil, api.Error{Kind: api.KindUnauthenticated}
	}

	return s.repo.IncomingFriendRequests(ctx, session.UserID)
}

func (s *Service) OutgoingFriendRequests(ctx context.Context) ([]social.Profile, error) {
	session, err := auth.Session(ctx)
	if err != nil {
		return nil, api.Error{Kind: api.KindUnauthenticated}
	}

	return s.repo.OutgoingFriendRequests(ctx, session.UserID)
}

func (s *Service) CreateFriendLink(ctx context.Context, id uuid.UUID) error {
	session, err := auth.Session(ctx)
	if err != nil {
		return api.Error{Kind: api.KindUnauthenticated}
	}

	return s.repo.CreateFriendLink(ctx, session.UserID, id)
}

func (s *Service) RemoveFriendLink(ctx context.Context, id uuid.UUID) error {
	session, err := auth.Session(ctx)
	if err != nil {
		return api.Error{Kind: api.KindUnauthenticated}
	}

	return s.repo.RemoveFriendLink(ctx, session.UserID, id)
}
