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

type ProfileStatistics struct {
	ID         uuid.UUID `db:"user_id"`
	Username   string    `db:"username"`
	MatchCount int       `db:"match_count"`
	WinCount   int       `db:"win_count"`
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

func (r *PostgresRepository) ProfileStatistics(ctx context.Context, id uuid.UUID) (social.ProfileStatistics, error) {
	sql := `
		SELECT
			u.id AS user_id,
			u.username,
			count(mp.match_id) AS match_count,
			count(mp.match_id) FILTER (WHERE mp.result = 1) AS win_count
		FROM users u
		LEFT JOIN match_participants mp
			ON u.id = mp.user_id
		WHERE u.id = $1
		GROUP BY u.id;
	`

	var dto ProfileStatistics
	err := r.db.Get(ctx, &dto, sql, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return social.ProfileStatistics{}, errs.ErrNotFound
		}

		return social.ProfileStatistics{}, err
	}

	res := social.ProfileStatistics{
		ID:         dto.ID,
		Username:   dto.Username,
		MatchCount: dto.MatchCount,
		WinCount:   dto.WinCount,
	}
	return res, nil
}
