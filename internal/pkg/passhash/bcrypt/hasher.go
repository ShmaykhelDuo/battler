package bcrypt

import (
	"errors"

	"github.com/ShmaykhelDuo/battler/internal/pkg/passhash"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct {
	cost int
}

func NewPasswordHasher(cost int) (PasswordHasher, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return PasswordHasher{}, bcrypt.InvalidCostError(cost)
	}

	return PasswordHasher{cost: cost}, nil
}

func (h PasswordHasher) Hash(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return nil, passhash.UnsupportedPasswordError{Err: err}
		}

		return nil, err
	}

	return hash, nil
}

func (h PasswordHasher) Check(hash []byte, password string) error {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return passhash.ErrPasswordMismatch
		}

		return err
	}

	return nil
}
