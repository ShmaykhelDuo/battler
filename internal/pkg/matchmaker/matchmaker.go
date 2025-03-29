package matchmaker

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/ShmaykhelDuo/battler/internal/bot/alphabeta2"
	model "github.com/ShmaykhelDuo/battler/internal/model/game"
)

type CharacterRepository interface {
	Characters() []int
}

type Matchmaker struct {
	in  chan model.MatchRequest
	out chan [2]model.MatchPlayer
	cr  CharacterRepository
}

func New(cr CharacterRepository) *Matchmaker {
	return &Matchmaker{
		in:  make(chan model.MatchRequest),
		out: make(chan [2]model.MatchPlayer),
		cr:  cr,
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
					res[0].Character = player1Req.Main
					res[1].Character = player2Req.Main
				} else if player1Req.Secondary != player2Req.Secondary {
					res[0].Character = player1Req.Secondary
					res[1].Character = player2Req.Secondary
				} else {
					res[0].Character = player1Req.Main
					res[1].Character = player2Req.Secondary
				}

			case <-time.After(5 * time.Second):
				res[1].PlayerID.IsBot = true
				res[1].Conn = alphabeta2.NewBot(5)

				res[1].Character = m.selectBotCharacter()
				if player1Req.Main != res[1].Character {
					res[0].Character = player1Req.Main
				} else {
					res[0].Character = player1Req.Secondary
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
