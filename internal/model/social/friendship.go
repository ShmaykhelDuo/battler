package social

import "github.com/google/uuid"

type FriendshipStatus int

const (
	FriendshipStatusUnknown         FriendshipStatus = 0
	FriendshipStatusNone            FriendshipStatus = 1
	FriendshipStatusOutgoingRequest FriendshipStatus = 2
	FriendshipStatusIncomingRequest FriendshipStatus = 3
	FriendshipStatusFriends         FriendshipStatus = 4
)

type ProfileFriendshipStatus struct {
	ID               uuid.UUID
	Username         string
	FriendshipStatus FriendshipStatus
}
