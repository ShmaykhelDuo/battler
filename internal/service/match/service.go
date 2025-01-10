package match

import (
	"context"
	"fmt"

	"github.com/ShmaykhelDuo/battler/internal/game/match"
	"github.com/ShmaykhelDuo/battler/internal/model/api"
	model "github.com/ShmaykhelDuo/battler/internal/model/game"
	"github.com/google/uuid"
)

type ConnectionRepository interface {
	CreateConnection(ctx context.Context, conn *model.Connection) error
	Connection(ctx context.Context, userID uuid.UUID) (*model.Connection, error)
}

type AvailableCharacterRepository interface {
	AreAllAvailable(ctx context.Context, userID uuid.UUID, numbers []int) (bool, error)
}

type Matchmaker interface {
	CreateMatch(ctx context.Context, conn match.Player, main, secondary int) error
}

type Service struct {
	cr  ConnectionRepository
	acr AvailableCharacterRepository
	mm  Matchmaker
}

func NewService(cr ConnectionRepository, acr AvailableCharacterRepository, mm Matchmaker) *Service {
	return &Service{
		cr:  cr,
		acr: acr,
		mm:  mm,
	}
}

func (s *Service) ConnectToNewMatch(ctx context.Context, userID uuid.UUID, main, secondary int) (*model.Connection, error) {
	avail, err := s.acr.AreAllAvailable(ctx, userID, []int{main, secondary})
	if err != nil {
		return nil, fmt.Errorf("available characters: %w", err)
	}

	if !avail {
		return nil, api.Error{
			Kind:    api.KindInvalidArgument,
			Message: "character is not available",
		}
	}

	if main == secondary {
		return nil, api.Error{
			Kind:    api.KindInvalidArgument,
			Message: "main and secondary characters must be different",
		}
	}

	conn := model.NewConnection(userID)
	err = s.cr.CreateConnection(ctx, conn)
	if err != nil {
		return nil, fmt.Errorf("create connection: %w", err)
	}

	err = s.mm.CreateMatch(ctx, conn, main, secondary)
	if err != nil {
		return nil, fmt.Errorf("matchmaker: %w", err)
	}

	return conn, nil
}

func (s *Service) ReconnectToMatch(ctx context.Context, userID uuid.UUID) (*model.Connection, error) {
	conn, err := s.cr.Connection(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get connection: %w", err)
	}

	return conn, nil
}
