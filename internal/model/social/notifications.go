package social

import "github.com/google/uuid"

type NewFriendRequestNotification struct {
	ID       uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
}

type FriendRequestAcceptedNotification struct {
	ID       uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
}
