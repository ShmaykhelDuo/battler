package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShmaykhelDuo/battler/internal/model/api"
	model "github.com/ShmaykhelDuo/battler/internal/model/auth"
	"github.com/ShmaykhelDuo/battler/internal/model/errs"
	"github.com/ShmaykhelDuo/battler/internal/pkg/auth"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u model.User) error
	UserByUsername(ctx context.Context, username string) (model.User, error)
}

type SessionRepository interface {
	CreateSession(ctx context.Context, s model.Session) error
	Session(ctx context.Context, id uuid.UUID) (model.Session, error)
	DeleteSession(ctx context.Context, id uuid.UUID) error
}

type PasswordHasher interface {
	Hash(password string) ([]byte, error)
	Check(hash []byte, password string) error
}

type Service struct {
	ur UserRepository
	sr SessionRepository
	h  PasswordHasher
}

func NewService(ur UserRepository, sr SessionRepository, h PasswordHasher) *Service {
	return &Service{
		ur: ur,
		sr: sr,
		h:  h,
	}
}

func (s *Service) Register(ctx context.Context, username string, password string) (string, error) {
	passwordHash, err := s.h.Hash(password)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	u := model.User{
		ID:           uuid.Must(uuid.NewV7()),
		Username:     username,
		PasswordHash: passwordHash,
	}

	err = s.ur.CreateUser(ctx, u)
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			return "", api.Error{
				Kind:    api.KindAlreadyExists,
				Message: "user with this username already exists",
			}
		}

		return "", fmt.Errorf("create user: %w", err)
	}

	return s.createNewSession(ctx, u.ID)
}

func (s *Service) SignIn(ctx context.Context, username string, password string) (string, error) {
	u, err := s.ur.UserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return "", api.Error{
				Kind: api.KindInvalidCredentials,
			}
		}

		return "", fmt.Errorf("find user by username: %w", err)
	}

	err = s.h.Check(u.PasswordHash, password)
	if err != nil {
		return "", api.Error{
			Kind: api.KindInvalidCredentials,
		}
	}

	return s.createNewSession(ctx, u.ID)
}

func (s *Service) SignOut(ctx context.Context) error {
	session, err := auth.Session(ctx)
	if err != nil {
		return api.Error{
			Kind: api.KindUnauthenticated,
		}
	}

	err = s.sr.DeleteSession(ctx, session.ID)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}

func (s *Service) createNewSession(ctx context.Context, userID uuid.UUID) (string, error) {
	session := model.Session{
		ID:     uuid.Must(uuid.NewV7()),
		UserID: userID,
	}

	err := s.sr.CreateSession(ctx, session)
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	return session.ID.String(), nil
}
