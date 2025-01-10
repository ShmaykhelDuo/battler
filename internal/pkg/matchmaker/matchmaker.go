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
}

func New(cr CharacterRepository) *Matchmaker {
	return &Matchmaker{
		cr: cr,
	}
}

func (m *Matchmaker) Run(ctx context.Context) error {
	return nil
}

func (m *Matchmaker) CreateMatch(ctx context.Context, conn match.Player, main, secondary int) error {
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

	match := match.New(player, bot, false)

	go match.Run(ctx)
	return nil
}

func (m *Matchmaker) selectBotCharacter() int {
	chars := m.cr.Characters()

	i := rand.IntN(len(chars))
	return chars[i]
}
