package match

import (
	"context"

	"github.com/ShmaykhelDuo/battler/internal/game/match"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db/postgres"
	"github.com/google/uuid"
)

type PostgresRepository struct {
	db *postgres.DB
}

func NewPostgresRepository(db *postgres.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateMatch(ctx context.Context, id uuid.UUID) error {
	sql := "INSERT INTO matches (id, created_at) VALUES ($1, now());"
	_, err := r.db.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) CreateMatchParticipant(ctx context.Context, userID uuid.UUID, matchID uuid.UUID, characterNum int, res match.ResultPlayer) error {
	sql := "INSERT INTO match_participants (user_id, match_id, character_number, result) VALUES ($1, $2, $3, $4);"

	_, err := r.db.Exec(ctx, sql, userID, matchID, characterNum, int(res.Status))
	if err != nil {
		return err
	}

	return nil
}
