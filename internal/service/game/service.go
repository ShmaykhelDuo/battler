package game

import (
	"context"
	"fmt"

	"github.com/ShmaykhelDuo/battler/internal/model/api"
	model "github.com/ShmaykhelDuo/battler/internal/model/game"
	"github.com/ShmaykhelDuo/battler/internal/pkg/auth"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db"
	"github.com/google/uuid"
)

type AvailableCharacterRepository interface {
	AvailableCharacters(ctx context.Context, userID uuid.UUID) ([]model.Character, error)
	AvailableCharactersCount(ctx context.Context, userID uuid.UUID) (int, error)
	AddCharacters(ctx context.Context, userID uuid.UUID, chars []model.Character) error
}

type CharacterPicker interface {
	RandomCharacters(n int) []int
}

type TransactionManager interface {
	Transact(ctx context.Context, isolation db.TxIsolation, f func(context.Context) error) error
}

type Service struct {
	cr AvailableCharacterRepository
	cp CharacterPicker
	tm TransactionManager
}

func NewService(cr AvailableCharacterRepository, cp CharacterPicker, tm TransactionManager) *Service {
	return &Service{
		cr: cr,
		cp: cp,
		tm: tm,
	}
}

func (s *Service) AvailableCharacters(ctx context.Context) ([]model.Character, error) {
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

func (s *Service) UnlockInitialCharacters(ctx context.Context) error {
	session, err := auth.Session(ctx)
	if err != nil {
		return api.Error{Kind: api.KindUnauthenticated}
	}

	charNums := s.cp.RandomCharacters(2)
	chars := make([]model.Character, len(charNums))
	for i, n := range charNums {
		chars[i] = model.Character{
			Number: n,
		}
	}

	err = s.tm.Transact(ctx, db.TxIsolationSerializable, func(ctx context.Context) error {
		count, err := s.cr.AvailableCharactersCount(ctx, session.UserID)
		if err != nil {
			return fmt.Errorf("get available characters count: %w", err)
		}

		if count > 0 {
			return api.Error{
				Kind:    api.KindAlreadyExists,
				Message: "initial characters were already unlocked",
			}
		}

		err = s.cr.AddCharacters(ctx, session.UserID, chars)
		if err != nil {
			return fmt.Errorf("add characters: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
