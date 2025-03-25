package matchmaker

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/ShmaykhelDuo/battler/internal/bot/alphabeta2"
	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type CharacterRepository interface {
	Character(number int) (*game.CharacterData, error)
	Characters() []int
}

type matchRequest struct {
	conn            match.Player
	main, secondary int
}

type Matchmaker struct {
	cr CharacterRepository
	in chan matchRequest
}

func New(cr CharacterRepository) *Matchmaker {
	return &Matchmaker{
		cr: cr,
		in: make(chan matchRequest),
	}
}

func (m *Matchmaker) Run(ctx context.Context) error {
	for {
		select {
		case player1Req := <-m.in:
			select {
			case player2Req := <-m.in:
				var char1Num, char2Num int
				if player1Req.main != player2Req.main {
					char1Num = player1Req.main
					char2Num = player2Req.main
				} else if player1Req.secondary != player2Req.secondary {
					char1Num = player1Req.secondary
					char2Num = player2Req.secondary
				} else {
					char1Num = player1Req.main
					char2Num = player2Req.secondary
				}

				err := m.createMatch(ctx, player1Req.conn, char1Num, player2Req.conn, char2Num)
				if err != nil {
					slog.Error("error while creating match", "err", err)
				}

			case <-time.After(5 * time.Second):
				botCharNum := m.selectBotCharacter()
				var charNum int
				if player1Req.main != botCharNum {
					charNum = player1Req.main
				} else {
					charNum = player1Req.secondary
				}

				bot := alphabeta2.NewBot(5)

				err := m.createMatch(ctx, player1Req.conn, charNum, bot, botCharNum)
				if err != nil {
					slog.Error("error while creating match", "err", err)
				}
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (m *Matchmaker) createMatch(ctx context.Context, p1 match.Player, char1Number int, p2 match.Player, char2Number int) error {
	char1, err := m.cr.Character(char1Number)
	if err != nil {
		return fmt.Errorf("get character: %w", err)
	}

	char2, err := m.cr.Character(char2Number)
	if err != nil {
		return fmt.Errorf("get character: %w", err)
	}

	player1 := match.CharacterPlayer{
		Character: game.NewCharacter(char1),
		Player:    p1,
	}

	player2 := match.CharacterPlayer{
		Character: game.NewCharacter(char2),
		Player:    p2,
	}

	match := match.New(player1, player2, false)
	go match.Run(ctx)

	return nil
}

func (m *Matchmaker) MakeMatch(ctx context.Context, conn match.Player, main, secondary int) error {
	req := matchRequest{
		conn:      conn,
		main:      main,
		secondary: secondary,
	}

	select {
	case m.in <- req:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *Matchmaker) selectBotCharacter() int {
	chars := m.cr.Characters()

	i := rand.IntN(len(chars))
	return chars[i]
}
