package money

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ShmaykhelDuo/battler/internal/model/api"
	"github.com/ShmaykhelDuo/battler/internal/model/errs"
	model "github.com/ShmaykhelDuo/battler/internal/model/money"
	"github.com/ShmaykhelDuo/battler/internal/model/notification"
	"github.com/ShmaykhelDuo/battler/internal/pkg/auth"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db"
	"github.com/google/uuid"
)

type BalanceRepository interface {
	Balance(ctx context.Context, userID uuid.UUID) (map[model.Currency]int64, error)
	CurrencyBalance(ctx context.Context, userID uuid.UUID, currency model.Currency) (int64, error)
	SetBalance(ctx context.Context, userID uuid.UUID, currency model.Currency, amount int64) error
}

type ConversionRepository interface {
	CreateConversion(ctx context.Context, userID uuid.UUID, c model.CurrencyConversion) error
	HasActiveConversionByUser(ctx context.Context, userID uuid.UUID) (bool, error)
	ActiveConversionByUser(ctx context.Context, userID uuid.UUID) (model.CurrencyConversion, error)
	ClaimConversion(ctx context.Context, id uuid.UUID) error
}

type TransactionManager interface {
	Transact(ctx context.Context, isolation db.TxIsolation, f func(context.Context) error) error
}

type NotificationRepository interface {
	CreateNotification(ctx context.Context, userID uuid.UUID, n notification.Notification) error
}

type Service struct {
	br BalanceRepository
	cr ConversionRepository
	tm TransactionManager
	nr NotificationRepository
}

func NewService(br BalanceRepository, cr ConversionRepository, tm TransactionManager, nr NotificationRepository) *Service {
	return &Service{
		br: br,
		cr: cr,
		tm: tm,
		nr: nr,
	}
}

func (s *Service) Balance(ctx context.Context) (map[model.Currency]int64, error) {
	session, err := auth.Session(ctx)
	if err != nil {
		return nil, api.Error{Kind: api.KindUnauthenticated}
	}

	balance, err := s.br.Balance(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("get balance: %w", err)
	}

	return balance, nil
}

func (s *Service) Convert(ctx context.Context, sourceCurrency model.Currency, amount int64) (model.CurrencyConversion, error) {
	session, err := auth.Session(ctx)
	if err != nil {
		return model.CurrencyConversion{}, api.Error{Kind: api.KindUnauthenticated}
	}

	if amount <= 0 {
		return model.CurrencyConversion{}, api.Error{
			Kind:    api.KindInvalidArgument,
			Message: "amount must be positive",
		}
	}

	spec, ok := model.CurrencyConversionSpecs[sourceCurrency]
	if !ok {
		return model.CurrencyConversion{}, api.Error{
			Kind:    api.KindInvalidArgument,
			Message: "specified currency is not convertable",
		}
	}

	targetAmount := amount / spec.Rate.Denom().Int64() * spec.Rate.Num().Int64()
	if targetAmount <= 0 {
		return model.CurrencyConversion{}, api.Error{
			Kind:    api.KindInvalidArgument,
			Message: "conversion must produce at least one unit of target currency",
		}
	}

	sourceAmount := targetAmount / spec.Rate.Num().Int64() * spec.Rate.Denom().Int64()

	duration := time.Duration(targetAmount) * spec.TimePerUnit

	var conv model.CurrencyConversion
	err = s.tm.Transact(ctx, db.TxIsolationRepeatableRead, func(ctx context.Context) error {
		hasConv, err := s.cr.HasActiveConversionByUser(ctx, session.UserID)
		if err != nil {
			return fmt.Errorf("get has active conversion: %w", err)
		}
		if hasConv {
			return api.Error{
				Kind:    api.KindAlreadyExists,
				Message: "there is already a conversion in progress",
			}
		}

		balance, err := s.br.CurrencyBalance(ctx, session.UserID, sourceCurrency)
		if err != nil {
			return fmt.Errorf("get currency balance: %w", err)
		}
		if balance < sourceAmount {
			return api.Error{
				Kind:    api.KindInvalidRequest,
				Message: "not enough source currency",
			}
		}

		balance -= sourceAmount
		err = s.br.SetBalance(ctx, session.UserID, sourceCurrency, balance)
		if err != nil {
			return fmt.Errorf("set balance: %w", err)
		}

		now := time.Now()

		conv = model.CurrencyConversion{
			ID:             uuid.Must(uuid.NewV7()),
			StartTime:      now,
			FinishTime:     now.Add(duration),
			TargetCurrency: spec.NextCurrency,
			TargetAmount:   targetAmount,
		}
		err = s.cr.CreateConversion(ctx, session.UserID, conv)
		if err != nil {
			return fmt.Errorf("create conversion: %w", err)
		}

		payload, err := json.Marshal(model.CurrencyConversionNotification{
			ConversionID:   conv.ID,
			TargetCurrency: spec.NextCurrency,
			TargetAmount:   targetAmount,
		})
		if err != nil {
			return fmt.Errorf("marshal notification: %w", err)
		}

		notification := notification.Notification{
			ID:         uuid.Must(uuid.NewV7()),
			Type:       notification.TypeCurrencyConversionFinished,
			Payload:    payload,
			CreateTime: conv.FinishTime,
		}

		err = s.nr.CreateNotification(ctx, session.UserID, notification)
		if err != nil {
			return fmt.Errorf("create notification: %w", err)
		}

		return nil
	})
	if err != nil {
		return model.CurrencyConversion{}, err
	}

	return conv, nil
}

func (s *Service) ConversionStatus(ctx context.Context) (model.CurrencyConversion, error) {
	session, err := auth.Session(ctx)
	if err != nil {
		return model.CurrencyConversion{}, api.Error{Kind: api.KindUnauthenticated}
	}

	conv, err := s.cr.ActiveConversionByUser(ctx, session.UserID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return model.CurrencyConversion{}, api.Error{
				Kind:    api.KindNotFound,
				Message: "no active conversions found",
			}
		}
		return model.CurrencyConversion{}, fmt.Errorf("get active conversion: %w", err)
	}

	return conv, nil
}

func (s *Service) ClaimConversion(ctx context.Context) (map[model.Currency]int64, error) {
	session, err := auth.Session(ctx)
	if err != nil {
		return nil, api.Error{Kind: api.KindUnauthenticated}
	}

	var balances map[model.Currency]int64
	err = s.tm.Transact(ctx, db.TxIsolationRepeatableRead, func(ctx context.Context) error {
		conv, err := s.cr.ActiveConversionByUser(ctx, session.UserID)
		if err != nil {
			if errors.Is(err, errs.ErrNotFound) {
				return api.Error{
					Kind:    api.KindNotFound,
					Message: "no active conversions found",
				}
			}
			return fmt.Errorf("get active conversion: %w", err)
		}

		if conv.FinishTime.After(time.Now()) {
			return api.Error{
				Kind:    api.KindInvalidRequest,
				Message: "conversion has not finished yet",
			}
		}

		balance, err := s.br.CurrencyBalance(ctx, session.UserID, conv.TargetCurrency)
		if err != nil {
			return fmt.Errorf("get currency balance: %w", err)
		}

		balance += conv.TargetAmount
		err = s.br.SetBalance(ctx, session.UserID, conv.TargetCurrency, balance)
		if err != nil {
			return fmt.Errorf("set currency balance: %w", err)
		}

		err = s.cr.ClaimConversion(ctx, conv.ID)
		if err != nil {
			return fmt.Errorf("claim conversion: %w", err)
		}

		balances, err = s.br.Balance(ctx, session.UserID)
		if err != nil {
			return fmt.Errorf("get balance: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return balances, nil
}
