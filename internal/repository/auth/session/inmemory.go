package session

import (
	"context"
	"errors"
	"sync"

	model "github.com/ShmaykhelDuo/battler/internal/model/auth"
	"github.com/ShmaykhelDuo/battler/internal/model/errs"
	"github.com/google/uuid"
)

type InMemoryRepository struct {
	sessions *sync.Map
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		sessions: &sync.Map{},
	}
}

func (r *InMemoryRepository) CreateSession(ctx context.Context, s model.Session) error {
	r.sessions.Store(s.ID, s)
	return nil
}

func (r *InMemoryRepository) Session(ctx context.Context, id uuid.UUID) (model.Session, error) {
	val, ok := r.sessions.Load(id)
	if !ok {
		return model.Session{}, errs.ErrNotFound
	}

	s, ok := val.(model.Session)
	if !ok {
		return model.Session{}, errors.New("invalid type stored")
	}

	return s, nil
}

func (r *InMemoryRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	r.sessions.Delete(id)
	return nil
}
