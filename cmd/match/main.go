package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ShmaykhelDuo/battler/internal/bot"
	"github.com/ShmaykhelDuo/battler/internal/bot/minimax"
	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/milana"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/ruby"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

func main() {
	// path := flag.String("path", "", "model path")
	// flag.Parse()

	// model, err := moveml.LoadModel(*path)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// bot1 := moveml.NewBot(model)
	bot1 := minimax.NewBot(minimax.TimeOptConcurrentRunner, 8)
	// bot2 := minimax.NewBot(4)
	bot2 := &bot.RandomBot{}

	// c1, c2 := getRandomPair()
	c1 := game.NewCharacter(ruby.CharacterRuby)
	c2 := game.NewCharacter(milana.CharacterMilana)

	p1 := match.CharacterPlayer{
		Character: c1,
		Player:    bot1,
	}
	p2 := match.CharacterPlayer{
		Character: c2,
		Player:    bot2,
	}

	m := match.New(p1, p2, false)

	err := m.Run(context.Background())
	if err != nil {
		log.Printf("match run: %v\n", err)
		return
	}

	res, err := m.Result()
	if err != nil {
		log.Printf("match result: %v\n", err)
		return
	}

	switch res {
	case match.ResultWonSecond:
		fmt.Println("Lost")
	case match.ResultWonFirst:
		fmt.Println("Won")
	default:
		fmt.Println("Draw")
	}
	fmt.Printf("%d : %d\n", c1.HP(), c2.HP())
}
