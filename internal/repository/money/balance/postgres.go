package balance

import (
	"context"
	"errors"

	model "github.com/ShmaykhelDuo/battler/internal/model/money"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Balance struct {
	CurrencyID int
	Amount     int64
}

type PostgresRepository struct {
	db *postgres.DB
}

func NewPostgresRepository(db *postgres.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Balance(ctx context.Context, userID uuid.UUID) (map[model.Currency]int64, error) {
	sql := "SELECT currency_id, amount FROM users_balances WHERE user_id = $1;"

	var dtos []Balance
	err := r.db.Select(ctx, &dtos, sql, userID)
	if err != nil {
		return nil, err
	}

	res := make(map[model.Currency]int64, len(dtos))
	for _, dto := range dtos {
		res[model.Currency(dto.CurrencyID)] = dto.Amount
	}

	return res, nil
}

func (r *PostgresRepository) CurrencyBalance(ctx context.Context, userID uuid.UUID, currency model.Currency) (int64, error) {
	sql := "SELECT amount FROM users_balances WHERE user_id = $1 AND currency_id = $2;"

	var amount int64
	err := r.db.Get(ctx, &amount, sql, userID, int(currency))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return amount, nil
}

func (r *PostgresRepository) SetBalance(ctx context.Context, userID uuid.UUID, currency model.Currency, amount int64) error {
	sql := `
		INSERT INTO users_balances (user_id, currency_id, amount) VALUES ($1, $2, $3)
		ON CONFLICT (user_id, currency_id) DO UPDATE SET amount = $3;
	`

	_, err := r.db.Exec(ctx, sql, userID, int(currency), amount)
	if err != nil {
		return err
	}

	return nil
}
