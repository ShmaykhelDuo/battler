package notification

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID         uuid.UUID
	Type       Type
	Payload    json.RawMessage
	CreateTime time.Time
}

type Type int

const (
	TypeUnknown                    Type = 0
	TypeCurrencyConversionFinished Type = 1
	TypeNewFriendRequest           Type = 2
	TypeFriendRequestAccepted      Type = 3
)
