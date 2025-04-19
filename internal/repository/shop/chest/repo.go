package chest

import (
	"context"
	"maps"
	"slices"

	"github.com/ShmaykhelDuo/battler/internal/model/errs"
	"github.com/ShmaykhelDuo/battler/internal/model/game"
	"github.com/ShmaykhelDuo/battler/internal/model/money"
	"github.com/ShmaykhelDuo/battler/internal/model/shop"
)

var chests = map[int]shop.Chest{
	1: ChestSTPack,
	2: ChestADPack,
	3: ChestSPPack,
	4: ChestRPPack,
	5: ChestLFPack,
}

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Chests(ctx context.Context) ([]shop.Chest, error) {
	return slices.Collect(maps.Values(chests)), nil
}

func (r *Repository) Chest(ctx context.Context, id int) (shop.Chest, error) {
	chest, ok := chests[id]
	if !ok {
		return shop.Chest{}, errs.ErrNotFound
	}

	return chest, nil
}

var ChestSTPack = shop.Chest{
	ID:              1,
	Name:            "ST Pack",
	PriceCurrency:   money.CurrencyWhiteDust,
	PriceAmount:     300,
	CharacterRarity: game.CharacterRarityST,
}

var ChestADPack = shop.Chest{
	ID:              2,
	Name:            "AD Pack",
	PriceCurrency:   money.CurrencyBlueDust,
	PriceAmount:     183,
	CharacterRarity: game.CharacterRarityAD,
	Available:       true,
}

var ChestSPPack = shop.Chest{
	ID:              3,
	Name:            "SP Pack",
	PriceCurrency:   money.CurrencyYellowDust,
	PriceAmount:     153,
	CharacterRarity: game.CharacterRaritySP,
	Available:       true,
}

var ChestRPPack = shop.Chest{
	ID:              4,
	Name:            "RP Pack",
	PriceCurrency:   money.CurrencyPurpleDust,
	PriceAmount:     72,
	CharacterRarity: game.CharacterRarityRP,
}

var ChestLFPack = shop.Chest{
	ID:              5,
	Name:            "LF Pack",
	PriceCurrency:   money.CurrencyStarDust,
	PriceAmount:     17,
	CharacterRarity: game.CharacterRarityLF,
	Available:       true,
}
