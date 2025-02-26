package matchmaker

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/bot"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type CharacterRepository interface {
	Character(number int) (*game.CharacterData, error)
	Characters() []int
}

type Matchmaker struct {
	cr CharacterRepository
	in chan [2]match.CharacterPlayer
}

func New(cr CharacterRepository) *Matchmaker {
	return &Matchmaker{
		cr: cr,
		in: make(chan [2]match.CharacterPlayer),
	}
}

func (m *Matchmaker) Run(ctx context.Context) error {
	for {
		select {
		case players := <-m.in:
			match := match.New(players[0], players[1], false)
			go match.Run(ctx)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (m *Matchmaker) MakeMatch(ctx context.Context, conn match.Player, main, secondary int) error {
	playerCharNum := main
	botCharNum := m.selectBotCharacter()

	if playerCharNum == botCharNum {
		playerCharNum = secondary
	}

	playerChar, err := m.cr.Character(playerCharNum)
	if err != nil {
		return fmt.Errorf("get character: %w", err)
	}

	botChar, err := m.cr.Character(botCharNum)
	if err != nil {
		return fmt.Errorf("get character: %w", err)
	}

	player := match.CharacterPlayer{
		Character: game.NewCharacter(playerChar),
		Player:    conn,
	}

	bot := match.CharacterPlayer{
		Character: game.NewCharacter(botChar),
		Player:    &bot.RandomBot{},
	}

	select {
	case m.in <- [2]match.CharacterPlayer{player, bot}:
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
