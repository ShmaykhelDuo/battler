package passhash

import (
	"errors"
	"fmt"
)

var ErrPasswordMismatch = errors.New("password mismatch")

type UnsupportedPasswordError struct {
	Err error
}

func (e UnsupportedPasswordError) Error() string {
	return fmt.Sprintf("password is not supported by hashing algorithm: %s", e.Err)
}
