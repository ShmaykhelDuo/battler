package match

import (
	"context"
	"fmt"

	"github.com/ShmaykhelDuo/battler/internal/game/match"
	"github.com/ShmaykhelDuo/battler/internal/model/api"
	model "github.com/ShmaykhelDuo/battler/internal/model/game"
	"github.com/ShmaykhelDuo/battler/internal/model/money"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db"
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
	MakeMatch(ctx context.Context, conn match.Player, main, secondary int) error
}

type BalanceRepository interface {
	CurrencyBalance(ctx context.Context, userID uuid.UUID, currency money.Currency) (int64, error)
	SetBalance(ctx context.Context, userID uuid.UUID, currency money.Currency, amount int64) error
}

type TransactionManager interface {
	Transact(ctx context.Context, isolation db.TxIsolation, f func(context.Context) error) error
}

type Service struct {
	cr  ConnectionRepository
	acr AvailableCharacterRepository
	mm  Matchmaker
	br  BalanceRepository
	tm  TransactionManager
}

func NewService(cr ConnectionRepository, acr AvailableCharacterRepository, mm Matchmaker, br BalanceRepository, tm TransactionManager) *Service {
	return &Service{
		cr:  cr,
		acr: acr,
		mm:  mm,
		br:  br,
		tm:  tm,
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

	conn := model.NewConnection(userID, s.matchEnd)
	err = s.cr.CreateConnection(ctx, conn)
	if err != nil {
		return nil, fmt.Errorf("create connection: %w", err)
	}

	err = s.mm.MakeMatch(ctx, conn, main, secondary)
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

func (s *Service) matchEnd(ctx context.Context, userID uuid.UUID) error {
	err := s.tm.Transact(ctx, db.TxIsolationDefault, func(ctx context.Context) error {
		balance, err := s.br.CurrencyBalance(ctx, userID, money.CurrencyWhiteDust)
		if err != nil {
			return fmt.Errorf("get currency balance: %w", err)
		}

		balance += 10
		err = s.br.SetBalance(ctx, userID, money.CurrencyWhiteDust, balance)
		if err != nil {
			return fmt.Errorf("set balance: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
