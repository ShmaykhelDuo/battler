package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	model "github.com/ShmaykhelDuo/battler/internal/model/auth"
	"github.com/ShmaykhelDuo/battler/internal/model/errs"
	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
	"github.com/ShmaykhelDuo/battler/internal/pkg/auth"
	"github.com/google/uuid"
)

type SessionRepository interface {
	Session(ctx context.Context, id uuid.UUID) (model.Session, error)
}

type AuthMiddleware struct {
	repo SessionRepository
}

func NewAuthMiddleware(repo SessionRepository) *AuthMiddleware {
	return &AuthMiddleware{repo: repo}
}

func (m *AuthMiddleware) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := m.getSession(r)
		if err == nil {
			ctx := auth.ContextWithSession(r.Context(), session)
			r = r.WithContext(ctx)
		} else if !errors.Is(err, errs.ErrNotFound) {
			api.HandleError(w, fmt.Errorf("get session: %w", err))
			return
		}

		h.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) getSession(r *http.Request) (model.Session, error) {
	c, err := r.Cookie("session_id")
	if err != nil {
		return model.Session{}, errs.ErrNotFound
	}

	sessionID, err := uuid.Parse(c.Value)
	if err != nil {
		return model.Session{}, errs.ErrNotFound
	}
	session, err := m.repo.Session(r.Context(), sessionID)
	if err != nil {
		return model.Session{}, err
	}

	return session, nil
}
