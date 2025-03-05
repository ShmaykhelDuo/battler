package conversion

import (
	"context"
	"errors"
	"time"

	"github.com/ShmaykhelDuo/battler/internal/model/errs"
	model "github.com/ShmaykhelDuo/battler/internal/model/money"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CurrencyConversion struct {
	ID             uuid.UUID
	StartTime      time.Time `db:"started_at"`
	FinishTime     time.Time `db:"finishes_at"`
	TargetCurrency int       `db:"target_currency_id"`
	TargetAmount   int64     `db:"amount"`
}

type PostgresRepository struct {
	db *postgres.DB
}

func NewPostgresRepository(db *postgres.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateConversion(ctx context.Context, userID uuid.UUID, c model.CurrencyConversion) error {
	sql := "INSERT INTO currency_conversions (id, user_id, started_at, finishes_at, target_currency_id, amount) VALUES ($1, $2, $3, $4, $5, $6);"

	_, err := r.db.Exec(ctx, sql, c.ID, userID, c.StartTime, c.FinishTime, int(c.TargetCurrency), c.TargetAmount)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) HasActiveConversionByUser(ctx context.Context, userID uuid.UUID) (bool, error) {
	sql := "SELECT count(id) FROM currency_conversions WHERE user_id = $1 AND is_claimed = false;"

	var count int
	err := r.db.Get(ctx, &count, sql, userID)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *PostgresRepository) ActiveConversionByUser(ctx context.Context, userID uuid.UUID) (model.CurrencyConversion, error) {
	sql := "SELECT id, started_at, finishes_at, target_currency_id, amount FROM currency_conversions WHERE user_id = $1 AND is_claimed = false;"

	var conv CurrencyConversion
	err := r.db.Get(ctx, &conv, sql, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.CurrencyConversion{}, errs.ErrNotFound
		}
		return model.CurrencyConversion{}, err
	}

	return model.CurrencyConversion{
		ID:             conv.ID,
		StartTime:      conv.StartTime,
		FinishTime:     conv.FinishTime,
		TargetCurrency: model.Currency(conv.TargetCurrency),
		TargetAmount:   conv.TargetAmount,
	}, nil
}

func (r *PostgresRepository) ClaimConversion(ctx context.Context, id uuid.UUID) error {
	sql := "UPDATE currency_conversions SET is_claimed = true WHERE id = $1;"

	tag, err := r.db.Exec(ctx, sql, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}
