package profile

import (
	"context"
	"errors"

	"github.com/ShmaykhelDuo/battler/internal/model/errs"
	"github.com/ShmaykhelDuo/battler/internal/model/social"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Profile struct {
	ID       uuid.UUID `db:"user_id"`
	Username string    `db:"username"`
}

type PostgresRepository struct {
	db *postgres.DB
}

func NewPostgresRepository(db *postgres.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Profile(ctx context.Context, id uuid.UUID) (social.Profile, error) {
	sql := "SELECT id AS user_id, username FROM users WHERE id = $1;"

	var dto Profile
	err := r.db.Get(ctx, &dto, sql, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return social.Profile{}, errs.ErrNotFound
		}

		return social.Profile{}, err
	}

	res := social.Profile{
		ID:       dto.ID,
		Username: dto.Username,
	}
	return res, nil
}
