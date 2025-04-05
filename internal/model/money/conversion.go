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

type CurrencyConversionNotification struct {
	ConversionID   uuid.UUID `json:"conversion_id"`
	TargetCurrency Currency  `json:"currency_id"`
	TargetAmount   int64     `json:"amount"`
}
