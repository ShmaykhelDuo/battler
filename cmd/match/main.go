package main

import (
	"fmt"
	"log"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/bot/minimax"
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
	bot1 := minimax.NewBot(8)
	bot2 := minimax.NewBot(8)

	// c1, c2 := getRandomPair()
	c1 := game.NewCharacter(ruby.CharacterRuby)
	c2 := game.NewCharacter(milana.CharacterMilana)

	res, err := match.Match(c1, c2, bot1, bot2)
	if err != nil {
		log.Printf("match: %v\n", err)
		return
	}
	switch res {
	case 1:
		fmt.Println("Lost")
	case -1:
		fmt.Println("Won")
	default:
		fmt.Println("Draw")
	}
	fmt.Printf("%d : %d\n", c1.HP(), c2.HP())
}
