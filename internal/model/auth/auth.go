package auth

import (
	"errors"

	"github.com/google/uuid"
)

var ErrInvalidCredentials = errors.New("invalid username or password")
var ErrNoSession = errors.New("no session is found")

type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash []byte
}

type Session struct {
	ID     uuid.UUID
	UserID uuid.UUID
}
