package social

import "github.com/google/uuid"

type Profile struct {
	ID       uuid.UUID
	Username string
}
