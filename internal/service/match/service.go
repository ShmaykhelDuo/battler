package match

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
	"github.com/ShmaykhelDuo/battler/internal/model/api"
	model "github.com/ShmaykhelDuo/battler/internal/model/game"
	"github.com/ShmaykhelDuo/battler/internal/model/money"
	"github.com/ShmaykhelDuo/battler/internal/model/social"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db"
	"github.com/google/uuid"
)

type ConnectionRepository interface {
	CreateConnection(ctx context.Context, conn *model.Connection) error
	Connection(ctx context.Context, userID uuid.UUID) (*model.Connection, error)
}

type AvailableCharacterRepository interface {
	AreAllAvailable(ctx context.Context, userID uuid.UUID, numbers []int) (bool, error)
	CharacterLevelExperience(ctx context.Context, userID uuid.UUID, number int) (level int, exp int, err error)
	UpdateCharacterLevelExperience(ctx context.Context, userID uuid.UUID, number int, level int, exp int) error
}

type Matchmaker interface {
	RequestMatch() chan<- model.MatchRequest
	Match() <-chan [2]model.MatchPlayer
}

type BalanceRepository interface {
	CurrencyBalance(ctx context.Context, userID uuid.UUID, currency money.Currency) (int64, error)
	SetBalance(ctx context.Context, userID uuid.UUID, currency money.Currency, amount int64) error
}

type TransactionManager interface {
	Transact(ctx context.Context, isolation db.TxIsolation, f func(context.Context) error) error
}

type CharacterRepository interface {
	Character(number int) (model.Character, error)
}

type MatchRepository interface {
	CreateMatch(ctx context.Context, id uuid.UUID) error
	CreateMatchParticipant(ctx context.Context, userID uuid.UUID, matchID uuid.UUID, characterNum int, res match.ResultPlayer) error
}

type ProfileRepository interface {
	Profile(ctx context.Context, id uuid.UUID) (social.Profile, error)
}

type Service struct {
	connRepo  ConnectionRepository
	acr       AvailableCharacterRepository
	mm        Matchmaker
	br        BalanceRepository
	tm        TransactionManager
	charRepo  CharacterRepository
	matchRepo MatchRepository
	profRepo  ProfileRepository
}

func NewService(connRepo ConnectionRepository, acr AvailableCharacterRepository, mm Matchmaker, br BalanceRepository, tm TransactionManager, charRepo CharacterRepository, matchRepo MatchRepository, profRepo ProfileRepository) *Service {
	return &Service{
		connRepo:  connRepo,
		acr:       acr,
		mm:        mm,
		br:        br,
		tm:        tm,
		charRepo:  charRepo,
		matchRepo: matchRepo,
		profRepo:  profRepo,
	}
}

func (s *Service) HandleMatches(ctx context.Context) error {
	for {
		select {
		case m := <-s.mm.Match():
			match, err := s.createMatch(m[0], m[1])
			if err != nil {
				slog.Error("failed to create match", "err", err)
				continue
			}
			go match.Run(ctx)
			go s.handleMatchResult(ctx, match, m)

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *Service) createMatch(p1, p2 model.MatchPlayer) (*match.Match, error) {
	char1, err := s.charRepo.Character(p1.Character)
	if err != nil {
		return nil, fmt.Errorf("get character: %w", err)
	}

	char2, err := s.charRepo.Character(p2.Character)
	if err != nil {
		return nil, fmt.Errorf("get character: %w", err)
	}

	player1 := match.CharacterPlayer{
		Character: game.NewCharacter(char1.Character),
		Player:    p1.Conn,
	}

	player2 := match.CharacterPlayer{
		Character: game.NewCharacter(char2.Character),
		Player:    p2.Conn,
	}

	match := match.New(player1, player2, false)
	return match, nil
}

func (s *Service) handleMatchResult(ctx context.Context, m *match.Match, players [2]model.MatchPlayer) {
	select {
	case res := <-m.Result():
		if res.Err != nil {
			slog.Error("match resulted in error", "err", res.Err)
			return
		}

		var ret1, ret2 model.MatchPlayerEndResult
		var prof1, prof2 social.Profile

		matchID := uuid.Must(uuid.NewV7())
		err := s.tm.Transact(ctx, db.TxIsolationRepeatableRead, func(ctx context.Context) error {
			err := s.matchRepo.CreateMatch(ctx, matchID)
			if err != nil {
				return fmt.Errorf("create match: %w", err)
			}

			if !players[0].PlayerID.IsBot {
				slog.Debug("saving match participant", "id", players[0].PlayerID.UserID)

				var err error
				ret1, err = s.updateMatchParticipant(ctx, players[0].PlayerID.UserID, matchID, players[0].Character, res.Res.Player1)
				if err != nil {
					return fmt.Errorf("update match participant: %w", err)
				}

				prof1, err = s.profRepo.Profile(ctx, players[0].PlayerID.UserID)
				if err != nil {
					return fmt.Errorf("get profile: %w", err)
				}
			}

			if !players[1].PlayerID.IsBot {
				slog.Debug("saving match participant", "id", players[1].PlayerID.UserID)

				var err error
				ret2, err = s.updateMatchParticipant(ctx, players[1].PlayerID.UserID, matchID, players[1].Character, res.Res.Player2)
				if err != nil {
					return fmt.Errorf("update match participant: %w", err)
				}

				prof2, err = s.profRepo.Profile(ctx, players[1].PlayerID.UserID)
				if err != nil {
					return fmt.Errorf("get profile: %w", err)
				}
			}

			return nil
		})
		if err != nil {
			slog.Error("failed to save match results", "err", err)
			return
		}

		ret1.OpponentProfile = prof2
		ret2.OpponentProfile = prof1

		conn1, ok := players[0].Conn.(*model.Connection)
		if ok {
			err := conn1.SendEndResult(ctx, ret1)
			if err != nil {
				return
			}
		}

		conn2, ok := players[1].Conn.(*model.Connection)
		if ok {
			err := conn2.SendEndResult(ctx, ret2)
			if err != nil {
				return
			}
		}

	case <-ctx.Done():
		return
	}
}

func (s *Service) updateMatchParticipant(ctx context.Context, userID uuid.UUID, matchID uuid.UUID, characterNum int, res match.ResultPlayer) (model.MatchPlayerEndResult, error) {
	err := s.matchRepo.CreateMatchParticipant(ctx, userID, matchID, characterNum, res)
	if err != nil {
		return model.MatchPlayerEndResult{}, fmt.Errorf("create match participant: %w", err)
	}

	character, err := s.charRepo.Character(characterNum)
	if err != nil {
		return model.MatchPlayerEndResult{}, fmt.Errorf("get character: %w", err)
	}

	level, exp, err := s.acr.CharacterLevelExperience(ctx, userID, characterNum)
	if err != nil {
		return model.MatchPlayerEndResult{}, fmt.Errorf("get character level experience: %w", err)
	}

	ret := model.MatchPlayerEndResult{
		Result:         res.Status,
		PrevLevel:      level,
		PrevExperience: exp,
	}

	if !res.HasGivenUp && level < len(model.CharacterLevelCaps[character.Rarity]) {
		exp += 1
		if exp >= model.CharacterLevelCaps[character.Rarity][level-1] {
			level += 1
			exp = 0
		}

		err = s.acr.UpdateCharacterLevelExperience(ctx, userID, characterNum, level, exp)
		if err != nil {
			return model.MatchPlayerEndResult{}, fmt.Errorf("update character level experience: %w", err)
		}
	}

	ret.Level = level
	ret.Experience = exp

	if res.HasGivenUp {
		ret.Reward = 0
	} else {
		switch res.Status {
		case match.ResultStatusWon:
			ret.Reward = 10 + model.ResultModifiers[character.Rarity]
		case match.ResultStatusLost:
			ret.Reward = 10 - model.ResultModifiers[character.Rarity]
		default:
			ret.Reward = 10
		}
	}

	if ret.Reward != 0 {
		balance, err := s.br.CurrencyBalance(ctx, userID, money.CurrencyWhiteDust)
		if err != nil {
			return model.MatchPlayerEndResult{}, fmt.Errorf("get currency balance: %w", err)
		}

		balance += int64(ret.Reward)
		err = s.br.SetBalance(ctx, userID, money.CurrencyWhiteDust, balance)
		if err != nil {
			return model.MatchPlayerEndResult{}, fmt.Errorf("set balance: %w", err)
		}
	}

	return ret, nil
}

func (s *Service) ConnectToNewMatch(ctx context.Context, userID uuid.UUID, main, secondary int) (*model.Connection, error) {
	avail, err := s.acr.AreAllAvailable(ctx, userID, []int{main, secondary})
	if err != nil {
		return nil, fmt.Errorf("available characters: %w", err)
	}

	if !avail {
		return nil, api.Error{
			Kind:    api.KindInvalidArgument,
			Message: "character is not available",
		}
	}

	if main == secondary {
		return nil, api.Error{
			Kind:    api.KindInvalidArgument,
			Message: "main and secondary characters must be different",
		}
	}

	conn := model.NewConnection(userID)
	err = s.connRepo.CreateConnection(ctx, conn)
	if err != nil {
		return nil, fmt.Errorf("create connection: %w", err)
	}

	req := model.MatchRequest{
		UserID:    userID,
		Conn:      conn,
		Main:      main,
		Secondary: secondary,
	}

	select {
	case s.mm.RequestMatch() <- req:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return conn, nil
}

func (s *Service) ReconnectToMatch(ctx context.Context, userID uuid.UUID) (*model.Connection, error) {
	conn, err := s.connRepo.Connection(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get connection: %w", err)
	}

	return conn, nil
}
