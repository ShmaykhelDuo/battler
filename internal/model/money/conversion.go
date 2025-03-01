package money

import (
	"time"

	"github.com/google/uuid"
)

type CurrencyConversion struct {
	ID             uuid.UUID
	StartTime      time.Time
	FinishTime     time.Time
	TargetCurrency Currency
	TargetAmount   int64
}
