package user

import (
	"context"
	"errors"

	model "github.com/ShmaykhelDuo/battler/internal/model/auth"
	"github.com/ShmaykhelDuo/battler/internal/model/errs"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash []byte
}

type PostgresRepository struct {
	db *postgres.DB
}

func NewPostgresRepository(db *postgres.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, u model.User) error {
	sql := "INSERT INTO users (id, username, password_hash) VALUES ($1, $2, $3);"

	_, err := r.db.Exec(ctx, sql, u.ID, u.Username, u.PasswordHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == postgres.ErrCodeUniqueViolation {
			return errs.ErrAlreadyExists
		}

		return err
	}

	return nil
}

func (r *PostgresRepository) UserByUsername(ctx context.Context, username string) (model.User, error) {
	sql := "SELECT id, username, password_hash FROM users WHERE username = $1;"

	var dto User
	err := r.db.Get(ctx, &dto, sql, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, errs.ErrNotFound
		}

		return model.User{}, err
	}

	u := model.User{
		ID:           dto.ID,
		Username:     dto.Username,
		PasswordHash: dto.PasswordHash,
	}
	return u, nil
}
