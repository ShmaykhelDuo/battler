package session

import (
	"context"
	"errors"
	"fmt"

	model "github.com/ShmaykhelDuo/battler/internal/model/auth"
	"github.com/ShmaykhelDuo/battler/internal/model/errs"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	cli *redis.Client
}

func NewRedisRepository(cli *redis.Client) *RedisRepository {
	return &RedisRepository{
		cli: cli,
	}
}

func (r *RedisRepository) CreateSession(ctx context.Context, s model.Session) error {
	key := fmt.Sprintf("sessions:%s", s.ID)
	err := r.cli.Set(ctx, key, s.UserID.String(), 0).Err()
	if err != nil {
		return fmt.Errorf("redis: %w", err)
	}

	return nil
}

func (r *RedisRepository) Session(ctx context.Context, id uuid.UUID) (model.Session, error) {
	key := fmt.Sprintf("sessions:%s", id)
	res, err := r.cli.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.Session{}, errs.ErrNotFound
		}
		return model.Session{}, fmt.Errorf("redis: %w", err)
	}

	userID, err := uuid.Parse(res)
	if err != nil {
		return model.Session{}, fmt.Errorf("uuid: %w", err)
	}

	return model.Session{
		ID:     id,
		UserID: userID,
	}, nil
}

func (r *RedisRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	key := fmt.Sprintf("sessions:%s", id)
	err := r.cli.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis: %w", err)
	}

	return nil
}
