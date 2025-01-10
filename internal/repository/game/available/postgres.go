package available

import (
	"context"
	"fmt"
	"strings"

	model "github.com/ShmaykhelDuo/battler/internal/model/game"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db/postgres"
	"github.com/google/uuid"
)

type Character struct {
	Number int
}

type PostgresRepository struct {
	db *postgres.DB
}

func NewPostgresRepository(db *postgres.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) AvailableCharacters(ctx context.Context, userID uuid.UUID) ([]model.Character, error) {
	sql := "SELECT number FROM available_characters WHERE user_id = $1;"

	var dto []Character
	err := r.db.Select(ctx, &dto, sql, userID)
	if err != nil {
		return nil, err
	}

	chars := make([]model.Character, len(dto))
	for i, c := range dto {
		chars[i] = model.Character{
			Number: c.Number,
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

func (r *PostgresRepository) AddCharacters(ctx context.Context, userID uuid.UUID, chars []model.Character) error {
	valuesSQL := make([]string, len(chars))
	args := make([]any, 2*len(chars))
	for i, c := range chars {
		valuesSQL[i] = fmt.Sprintf("($%d, $%d)", 2*i+1, 2*i+2)
		args[2*i] = userID
		args[2*i+1] = c.Number
	}

	sql := fmt.Sprintf("INSERT INTO available_characters (user_id, number) VALUES %s;", strings.Join(valuesSQL, ", "))

	_, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
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
