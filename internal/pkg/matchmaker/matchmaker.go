package matchmaker

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/ShmaykhelDuo/battler/internal/game/match"
	model "github.com/ShmaykhelDuo/battler/internal/model/game"
)

type CharacterRepository interface {
	Characters() []int
}

type BotFactory interface {
	Bot(botChar, playerChar int, playerLevel int) (match.Player, error)
}

type Matchmaker struct {
	in  chan model.MatchRequest
	out chan [2]model.MatchPlayer
	cr  CharacterRepository
	bf  BotFactory
}

func New(cr CharacterRepository, bf BotFactory) *Matchmaker {
	return &Matchmaker{
		in:  make(chan model.MatchRequest),
		out: make(chan [2]model.MatchPlayer),
		cr:  cr,
		bf:  bf,
	}
}

func (m *Matchmaker) Run(ctx context.Context) error {
	for {
		select {
		case player1Req := <-m.in:
			var res [2]model.MatchPlayer
			res[0].PlayerID.UserID = player1Req.UserID
			res[0].Conn = player1Req.Conn

			select {
			case player2Req := <-m.in:
				res[1].PlayerID.UserID = player2Req.UserID
				res[1].Conn = player2Req.Conn

				if player1Req.Main != player2Req.Main {
					res[0].Character = player1Req.Main.Number
					res[1].Character = player2Req.Main.Number
				} else if player1Req.Secondary != player2Req.Secondary {
					res[0].Character = player1Req.Secondary.Number
					res[1].Character = player2Req.Secondary.Number
				} else {
					res[0].Character = player1Req.Main.Number
					res[1].Character = player2Req.Secondary.Number
				}

			case <-time.After(5 * time.Second):
				res[1].PlayerID.IsBot = true

				res[1].Character = m.selectBotCharacter()

				var playerLevel int
				if player1Req.Main.Number != res[1].Character {
					res[0].Character = player1Req.Main.Number
					playerLevel = player1Req.Main.Level
				} else {
					res[0].Character = player1Req.Secondary.Number
					playerLevel = player1Req.Secondary.Level
				}

				var err error
				res[1].Conn, err = m.bf.Bot(res[1].Character, res[0].Character, playerLevel)
				if err != nil {
					slog.Error("failed to create bot", "err", err)
				}

			case <-ctx.Done():
				return ctx.Err()
			}

			select {
			case m.out <- res:
			case <-ctx.Done():
				return ctx.Err()
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (m *Matchmaker) RequestMatch() chan<- model.MatchRequest {
	return m.in
}

func (m *Matchmaker) Match() <-chan [2]model.MatchPlayer {
	return m.out
}

func (m *Matchmaker) selectBotCharacter() int {
	chars := m.cr.Characters()

	i := rand.IntN(len(chars))
	return chars[i]
}
