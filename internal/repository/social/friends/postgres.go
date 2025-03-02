package friends

import (
	"context"
	"errors"

	"github.com/ShmaykhelDuo/battler/internal/model/errs"
	"github.com/ShmaykhelDuo/battler/internal/model/social"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *PostgresRepository) Friends(ctx context.Context, userID uuid.UUID) ([]social.Profile, error) {
	sql := `
		SELECT a.friend_id AS user_id, u.username FROM users_friends a
		JOIN users_friends b
			ON a.user_id = b.friend_id
			AND a.friend_id = b.user_id
		JOIN users u ON a.friend_id = u.id
		WHERE a.user_id = $1;
	`

	var res []Profile
	err := r.db.Select(ctx, &res, sql, userID)
	if err != nil {
		return nil, err
	}

	return dtoToProfiles(res), nil
}

func (r *PostgresRepository) IncomingFriendRequests(ctx context.Context, userID uuid.UUID) ([]social.Profile, error) {
	sql := `
		SELECT a.user_id, u.username FROM users_friends a
		LEFT JOIN users_friends b
			ON a.user_id = b.friend_id
			AND a.friend_id = b.user_id
		JOIN users u ON a.user_id = u.id
		WHERE a.friend_id = $1
			AND b.user_id IS NULL;
	`

	var res []Profile
	err := r.db.Select(ctx, &res, sql, userID)
	if err != nil {
		return nil, err
	}

	return dtoToProfiles(res), nil
}

func (r *PostgresRepository) OutgoingFriendRequests(ctx context.Context, userID uuid.UUID) ([]social.Profile, error) {
	sql := `
		SELECT a.friend_id AS user_id, u.username FROM users_friends a
		LEFT JOIN users_friends b
			ON a.user_id = b.friend_id
			AND a.friend_id = b.user_id
		JOIN users u ON a.friend_id = u.id
		WHERE a.user_id = $1
			AND b.user_id IS NULL;
	`

	var res []Profile
	err := r.db.Select(ctx, &res, sql, userID)
	if err != nil {
		return nil, err
	}

	return dtoToProfiles(res), nil
}

func (r *PostgresRepository) CreateFriendLink(ctx context.Context, userID uuid.UUID, otherID uuid.UUID) error {
	sql := "INSERT INTO users_friends (user_id, friend_id) VALUES ($1, $2);"

	_, err := r.db.Exec(ctx, sql, userID, otherID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == postgres.ErrCodeUniqueViolation {
			return errs.ErrAlreadyExists
		}

		return err
	}

	return nil
}

func (r *PostgresRepository) RemoveFriendLink(ctx context.Context, userID uuid.UUID, otherID uuid.UUID) error {
	sql := "DELETE FROM users_friends WHERE user_id = $1 AND friend_id = $2;"

	_, err := r.db.Exec(ctx, sql, userID, otherID)
	if err != nil {
		return err
	}

	return nil
}

func dtoToProfiles(dtos []Profile) []social.Profile {
	res := make([]social.Profile, len(dtos))
	for i, dto := range dtos {
		res[i] = social.Profile{
			ID:       dto.ID,
			Username: dto.Username,
		}
	}
	return res
}
