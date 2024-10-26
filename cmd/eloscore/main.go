package main

import (
	"fmt"
	"log"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/bot"
	"github.com/ShmaykhelDuo/battler/internal/game/bot/minimax"
	"github.com/ShmaykhelDuo/battler/internal/game/bot/moveml"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/milana"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/ruby"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type playerGenerator interface {
	Player(desc game.CharacterDescription) match.Player
}

type randomPlayerGenerator struct {
}

func (g randomPlayerGenerator) Player(desc game.CharacterDescription) match.Player {
	return &bot.RandomBot{}
}

type miniMaxPlayerGenerator struct {
	depth int
}

func (g miniMaxPlayerGenerator) Player(desc game.CharacterDescription) match.Player {
	return minimax.NewBot(g.depth)
}

type moveMLPlayerGenerator struct {
	models map[game.CharacterDescription]*moveml.Model
}

func (g moveMLPlayerGenerator) Player(desc game.CharacterDescription) match.Player {
	return moveml.NewBot(g.models[desc])
}

func main() {
	rubyModel, err := moveml.LoadModel("../../../ml/models/ruby-vs-milana")
	if err != nil {
		log.Fatalf("failed to load Ruby model: %v\n", err)
	}

	milanaModel, err := moveml.LoadModel("../../../ml/models/milana-vs-ruby")
	if err != nil {
		log.Fatalf("failed to load Milana model: %v\n", err)
	}

	players := map[string]playerGenerator{
		"random":   randomPlayerGenerator{},
		"minimax3": miniMaxPlayerGenerator{depth: 3},
		"minimax5": miniMaxPlayerGenerator{depth: 5},
		"minimax7": miniMaxPlayerGenerator{depth: 7},
		// "minimax8": miniMaxPlayerGenerator{depth: 8},
		"ml": moveMLPlayerGenerator{
			models: map[game.CharacterDescription]*moveml.Model{
				ruby.CharacterRuby.Desc:     rubyModel,
				milana.CharacterMilana.Desc: milanaModel,
			},
		},
	}

	s := newEloScoring(ruby.CharacterRuby, milana.CharacterMilana, players)
	err = s.run(25)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Result: %v\n", s.ratings)
}
