package money

import (
	"math/big"
	"time"
)

type Currency int

const (
	CurrencyWhiteDust  Currency = 1
	CurrencyBlueDust   Currency = 2
	CurrencyYellowDust Currency = 3
	CurrencyPurpleDust Currency = 4
	CurrencyStarDust   Currency = 5
)

type CurrencyConversionSpec struct {
	NextCurrency Currency
	Rate         *big.Rat
	TimePerUnit  time.Duration
}

var CurrencyConversionSpecs = map[Currency]CurrencyConversionSpec{
	CurrencyWhiteDust: {
		NextCurrency: CurrencyBlueDust,
		Rate:         big.NewRat(1, 2),
		TimePerUnit:  24 * time.Second,
	},
	CurrencyBlueDust: {
		NextCurrency: CurrencyYellowDust,
		Rate:         big.NewRat(2, 5),
		TimePerUnit:  30 * time.Second,
	},
	CurrencyYellowDust: {
		NextCurrency: CurrencyPurpleDust,
		Rate:         big.NewRat(1, 4),
		TimePerUnit:  45 * time.Second,
	},
	CurrencyPurpleDust: {
		NextCurrency: CurrencyStarDust,
		Rate:         big.NewRat(1, 5),
		TimePerUnit:  60 * time.Second,
	},
}
