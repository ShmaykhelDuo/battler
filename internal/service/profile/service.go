package profile

import (
	"context"
	"fmt"

	"github.com/ShmaykhelDuo/battler/internal/model/api"
	"github.com/ShmaykhelDuo/battler/internal/model/social"
	"github.com/ShmaykhelDuo/battler/internal/pkg/auth"
	"github.com/google/uuid"
)

type ProfileRepository interface {
	Profile(ctx context.Context, id uuid.UUID) (social.Profile, error)
}

type Service struct {
	pr ProfileRepository
}

func NewService(pr ProfileRepository) *Service {
	return &Service{
		pr: pr,
	}
}

func (s *Service) Profile(ctx context.Context) (social.Profile, error) {
	session, err := auth.Session(ctx)
	if err != nil {
		return social.Profile{}, api.Error{Kind: api.KindUnauthenticated}
	}

	profile, err := s.pr.Profile(ctx, session.UserID)
	if err != nil {
		return social.Profile{}, fmt.Errorf("get profile: %w", err)
	}

	return profile, nil
}
