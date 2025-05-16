package connection

import (
	"context"
	"errors"
	"sync"

	"github.com/ShmaykhelDuo/battler/internal/model/errs"
	model "github.com/ShmaykhelDuo/battler/internal/model/game"
	"github.com/google/uuid"
)

type InMemoryRepository struct {
	conns sync.Map
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		conns: sync.Map{},
	}
}

func (r *InMemoryRepository) CreateConnection(ctx context.Context, conn *model.Connection) error {
	userID := conn.UserID()
	r.conns.Store(userID, conn)
	return nil
}

func (r *InMemoryRepository) Connection(ctx context.Context, userID uuid.UUID) (*model.Connection, error) {
	val, ok := r.conns.Load(userID)
	if !ok {
		return nil, errs.ErrNotFound
	}

	conn, ok := val.(*model.Connection)
	if !ok {
		return nil, errors.New("invalid conn type")
	}

	return conn, nil
}

func (r *InMemoryRepository) RemoveConnection(ctx context.Context, userID uuid.UUID) error {
	r.conns.Delete(userID)
	return nil
}
