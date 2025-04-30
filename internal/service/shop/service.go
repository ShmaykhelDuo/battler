package shop

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShmaykhelDuo/battler/internal/model/api"
	"github.com/ShmaykhelDuo/battler/internal/model/errs"
	"github.com/ShmaykhelDuo/battler/internal/model/game"
	"github.com/ShmaykhelDuo/battler/internal/model/money"
	model "github.com/ShmaykhelDuo/battler/internal/model/shop"
	"github.com/ShmaykhelDuo/battler/internal/pkg/auth"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db"
	"github.com/google/uuid"
)

type ChestRepository interface {
	Chests(ctx context.Context) ([]model.Chest, error)
	Chest(ctx context.Context, id int) (model.Chest, error)
}

type BalanceRepository interface {
	CurrencyBalance(ctx context.Context, userID uuid.UUID, currency money.Currency) (int64, error)
	SetBalance(ctx context.Context, userID uuid.UUID, currency money.Currency, amount int64) error
}

type CharacterPicker interface {
	RandomCharacter() int
	RandomCharacterOfRarity(rarity game.CharacterRarity) int
}

type AvailableCharacterRepository interface {
	AddCharacters(ctx context.Context, userID uuid.UUID, chars []int) error
	AreAllAvailable(ctx context.Context, userID uuid.UUID, numbers []int) (bool, error)
}

type TransactionManager interface {
	Transact(ctx context.Context, isolation db.TxIsolation, f func(context.Context) error) error
}

type Service struct {
	chestRepo     ChestRepository
	balanceRepo   BalanceRepository
	charPicker    CharacterPicker
	availCharRepo AvailableCharacterRepository
	tm            TransactionManager
}

func NewService(chestRepo ChestRepository, balanceRepo BalanceRepository, charPicker CharacterPicker, availCharRepo AvailableCharacterRepository, tm TransactionManager) *Service {
	return &Service{
		chestRepo:     chestRepo,
		balanceRepo:   balanceRepo,
		charPicker:    charPicker,
		availCharRepo: availCharRepo,
		tm:            tm,
	}
}

func (s *Service) Chests(ctx context.Context) ([]model.Chest, error) {
	return s.chestRepo.Chests(ctx)
}

func (s *Service) BuyChest(ctx context.Context, chestID int) (game.AvailableCharacter, error) {
	session, err := auth.Session(ctx)
	if err != nil {
		return game.AvailableCharacter{}, api.Error{Kind: api.KindUnauthenticated}
	}

	var charNum int
	err = s.tm.Transact(ctx, db.TxIsolationRepeatableRead, func(ctx context.Context) error {
		chest, err := s.chestRepo.Chest(ctx, chestID)
		if err != nil {
			if errors.Is(err, errs.ErrNotFound) {
				return api.Error{
					Kind:    api.KindNotFound,
					Message: "chest with provided id not found",
				}
			}
			return fmt.Errorf("get chest: %w", err)
		}

		if !chest.Available {
			return api.Error{
				Kind:    api.KindInvalidRequest,
				Message: "chest is not available for purchase",
			}
		}

		balance, err := s.balanceRepo.CurrencyBalance(ctx, session.UserID, chest.PriceCurrency)
		if err != nil {
			return fmt.Errorf("get balance: %w", err)
		}

		if balance < chest.PriceAmount {
			return api.Error{
				Kind:    api.KindInvalidRequest,
				Message: "not enough currency for purchase",
			}
		}

		balance -= chest.PriceAmount
		err = s.balanceRepo.SetBalance(ctx, session.UserID, chest.PriceCurrency, balance)
		if err != nil {
			return fmt.Errorf("set balance: %w", err)
		}

		charNum = s.charPicker.RandomCharacterOfRarity(chest.CharacterRarity)

		isAvail, err := s.availCharRepo.AreAllAvailable(ctx, session.UserID, []int{charNum})
		if err != nil {
			return fmt.Errorf("is character avail: %w", err)
		}

		if !isAvail {
			err = s.availCharRepo.AddCharacters(ctx, session.UserID, []int{charNum})
			if err != nil {
				return fmt.Errorf("add characters: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return game.AvailableCharacter{}, err
	}

	return game.AvailableCharacter{
		Number:          charNum,
		Level:           1,
		LevelExperience: 0,
		MatchCount:      0,
		WinCount:        0,
	}, nil
}
