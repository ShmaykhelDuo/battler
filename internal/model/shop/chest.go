package shop

import (
	"github.com/ShmaykhelDuo/battler/internal/model/game"
	"github.com/ShmaykhelDuo/battler/internal/model/money"
)

type Chest struct {
	ID              int
	Name            string
	PriceCurrency   money.Currency
	PriceAmount     int64
	CharacterRarity game.CharacterRarity
	Available       bool
}
