package social

import "github.com/google/uuid"

type Profile struct {
	ID       uuid.UUID
	Username string
}

type ProfileStatistics struct {
	ID         uuid.UUID
	Username   string
	MatchCount int
	WinCount   int
}
