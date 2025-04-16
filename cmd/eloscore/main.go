package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ShmaykhelDuo/battler/internal/bot"
	"github.com/ShmaykhelDuo/battler/internal/bot/alphabeta2"
	"github.com/ShmaykhelDuo/battler/internal/bot/ml"
	mlbot "github.com/ShmaykhelDuo/battler/internal/bot/ml/bot"
	"github.com/ShmaykhelDuo/battler/internal/bot/ml/formats"
	"github.com/ShmaykhelDuo/battler/internal/bot/ml/model"
	"github.com/ShmaykhelDuo/battler/internal/game"
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
	// return minimax.NewBot(minimax.MemOptConcurrentRunner, g.depth)
	return alphabeta2.NewBot(g.depth)
}

type modelPlayerGenerator struct {
	models map[game.CharacterDescription]*model.Model
}

func (g modelPlayerGenerator) Player(desc game.CharacterDescription) match.Player {
	return mlbot.NewBot(g.models[desc])
}

type modelDesc struct {
	prefix string
	format ml.Format
}

func main() {
	modelRoot := "ml/models"

	modelNames := map[string]modelDesc{
		"policy_10_51": modelDesc{
			prefix: "action_0_observation_",
			format: formats.FullStateFormat{},
		},
		"policy_51_10": modelDesc{
			prefix: "action_0_observation_",
			format: formats.FullStateFormat{},
		},
	}

	models := make(map[string]*model.Model)
	for name, desc := range modelNames {
		path := filepath.Join(modelRoot, name)
		var err error
		models[name], err = model.LoadModel(path, desc.format, desc.prefix)
		if err != nil {
			log.Fatalf("failed to load model %s from %s: %v", name, path, err)
		}
	}

	players := map[string]playerGenerator{
		"random":   randomPlayerGenerator{},
		"minimax2": miniMaxPlayerGenerator{depth: 2},
		"minimax4": miniMaxPlayerGenerator{depth: 4},
		// "minimax6":  miniMaxPlayerGenerator{depth: 6},
		// "minimax8":  miniMaxPlayerGenerator{depth: 8},
		// "minimax10": miniMaxPlayerGenerator{depth: 10},
		"mldqn": modelPlayerGenerator{
			models: map[game.CharacterDescription]*model.Model{
				ruby.CharacterRuby.Desc:     models["policy_10_51"],
				milana.CharacterMilana.Desc: models["policy_51_10"],
			},
		},
	}

	var res []map[string]int

	for range 10 {
		s := newEloScoring(ruby.CharacterRuby, milana.CharacterMilana, players)
		err := s.run(50)
		if err != nil {
			log.Fatal(err)
		}

		res = append(res, s.ratings)
	}

	fmt.Printf("%v\n", res)

	f, err := os.Create("ratings.csv")
	if err != nil {
		log.Fatalf("rip file: %v", err)
	}
	defer f.Close()

	p := make([]string, 0, len(players))
	for k := range players {
		p = append(p, k)
	}

	w := csv.NewWriter(f)
	w.Write(p)
	for _, row := range res {
		r := make([]string, len(p))
		for i, k := range p {
			r[i] = strconv.Itoa(row[k])
		}
		w.Write(r)
	}
	w.Flush()
}
