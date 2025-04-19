package game

import (
	"context"
	"fmt"

	"github.com/ShmaykhelDuo/battler/internal/model/api"
	model "github.com/ShmaykhelDuo/battler/internal/model/game"
	"github.com/ShmaykhelDuo/battler/internal/pkg/auth"
	"github.com/google/uuid"
)

type AvailableCharacterRepository interface {
	AvailableCharacters(ctx context.Context, userID uuid.UUID) ([]model.AvailableCharacter, error)
}

type Service struct {
	cr AvailableCharacterRepository
}

func NewService(cr AvailableCharacterRepository) *Service {
	return &Service{
		cr: cr,
	}
}

func (s *Service) AvailableCharacters(ctx context.Context) ([]model.AvailableCharacter, error) {
	session, err := auth.Session(ctx)
	if err != nil {
		return nil, api.Error{Kind: api.KindUnauthenticated}
	}

	chars, err := s.cr.AvailableCharacters(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("get available characters: %w", err)
	}

	return chars, nil
}
