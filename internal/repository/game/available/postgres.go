package available

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ShmaykhelDuo/battler/internal/model/errs"
	model "github.com/ShmaykhelDuo/battler/internal/model/game"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type Character struct {
	Number          int `db:"number"`
	Level           int `db:"level"`
	LevelExperience int `db:"level_experience"`
	MatchCount      int `db:"match_count"`
	WinCount        int `db:"win_count"`
}

type CharacterLevelExperience struct {
	Level           int `db:"level"`
	LevelExperience int `db:"level_experience"`
}

type PostgresRepository struct {
	db *postgres.DB
}

func NewPostgresRepository(db *postgres.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) AvailableCharacters(ctx context.Context, userID uuid.UUID) ([]model.AvailableCharacter, error) {
	sql := `
		SELECT
			ac.number, ac.level, ac.level_experience,
			count(mp.match_id) AS match_count,
			count(mp.match_id) FILTER (WHERE mp.result = 1) AS win_count
		FROM available_characters ac
		LEFT JOIN match_participants mp
			ON ac.user_id = mp.user_id
			AND ac.number = mp.character_number
		WHERE ac.user_id = $1
		GROUP BY ac.user_id, ac.number
		ORDER BY ac.number;
	`

	var dto []Character
	err := r.db.Select(ctx, &dto, sql, userID)
	if err != nil {
		return nil, err
	}

	chars := make([]model.AvailableCharacter, len(dto))
	for i, c := range dto {
		chars[i] = model.AvailableCharacter{
			Number:          c.Number,
			Level:           c.Level,
			LevelExperience: c.LevelExperience,
			MatchCount:      c.MatchCount,
			WinCount:        c.WinCount,
		}
	}
	return chars, nil
}

func (r *PostgresRepository) AvailableCharactersCount(ctx context.Context, userID uuid.UUID) (int, error) {
	sql := "SELECT count(*) FROM available_characters WHERE user_id = $1;"

	var count int
	err := r.db.Get(ctx, &count, sql, userID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresRepository) AddCharacters(ctx context.Context, userID uuid.UUID, numbers []int) error {
	valuesSQL := make([]string, len(numbers))
	args := make([]any, 2*len(numbers))
	for i, num := range numbers {
		valuesSQL[i] = fmt.Sprintf("($%d, $%d)", 2*i+1, 2*i+2)
		args[2*i] = userID
		args[2*i+1] = num
	}

	sql := fmt.Sprintf("INSERT INTO available_characters (user_id, number) VALUES %s;", strings.Join(valuesSQL, ", "))

	_, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == postgres.ErrCodeUniqueViolation {
			return errs.ErrAlreadyExists
		}

		return err
	}

	return nil
}

func (r *PostgresRepository) AreAllAvailable(ctx context.Context, userID uuid.UUID, numbers []int) (bool, error) {
	sql := "SELECT count(*) FROM available_characters WHERE user_id = $1 AND number = ANY ($2);"

	var count int
	err := r.db.Get(ctx, &count, sql, userID, numbers)
	if err != nil {
		return false, err
	}

	return count == len(numbers), nil
}

func (r *PostgresRepository) CharacterLevelExperience(ctx context.Context, userID uuid.UUID, number int) (level int, exp int, err error) {
	sql := "SELECT level, level_experience FROM available_characters WHERE user_id = $1 AND number = $2;"

	var res CharacterLevelExperience
	err = r.db.Get(ctx, &res, sql, userID, number)
	if err != nil {
		return 0, 0, err
	}

	return res.Level, res.LevelExperience, nil
}

func (r *PostgresRepository) UpdateCharacterLevelExperience(ctx context.Context, userID uuid.UUID, number int, level int, exp int) error {
	sql := "UPDATE available_characters SET level = $3, level_experience = $4 WHERE user_id = $1 AND number = $2;"

	_, err := r.db.Exec(ctx, sql, userID, number, level, exp)
	if err != nil {
		return err
	}

	return nil
}
