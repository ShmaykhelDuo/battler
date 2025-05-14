package bot_test

import (
	"context"
	"log"
	"math/rand/v2"
	"path/filepath"
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/bot"
	"github.com/ShmaykhelDuo/battler/internal/bot/alphabeta2"
	"github.com/ShmaykhelDuo/battler/internal/bot/minimax"
	"github.com/ShmaykhelDuo/battler/internal/bot/ml"
	mlbot "github.com/ShmaykhelDuo/battler/internal/bot/ml/bot"
	"github.com/ShmaykhelDuo/battler/internal/bot/ml/formats"
	"github.com/ShmaykhelDuo/battler/internal/bot/ml/model"
	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/milana"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/ruby"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

func benchmarkBot(b *testing.B, bot match.Player, state match.GameState) {
	for i := 0; i < b.N; i++ {
		bot.SendState(context.Background(), state)
		bot.RequestSkill(context.Background())
	}
}

func BenchmarkBots(b *testing.B) {
	bots := players()
	// s := states()

	state := match.GameState{
		Character: game.NewCharacter(ruby.CharacterRuby),
		Opponent:  game.NewCharacter(milana.CharacterMilana),
		TurnState: game.TurnState{
			TurnNum:      1,
			IsGoingFirst: true,
			IsTurnEnd:    false,
		},
		SkillsLeft: 1,
		SkillLog:   make(match.SkillLog),
		PlayerTurn: true,
		AsOpp:      false,
	}

	for key, bot := range bots {
		b.Run(key, func(b *testing.B) {
			benchmarkBot(b, bot, state)
		})
	}
}

type modelDesc struct {
	prefix string
	format ml.Format
}

func players() map[string]match.Player {
	modelRoot := "../../ml/models"

	modelNames := map[string]modelDesc{
		// "ruby-vs-milana-val20": modelDesc{
		// 	prefix: "serve_",
		// 	format: formats.FullStateCringeFormat{},
		// },
		// "ruby-vs-milana-noval": modelDesc{
		// 	prefix: "serve_",
		// 	format: formats.FullStateCringeFormat{},
		// },
		"policy_10_51": modelDesc{
			prefix: "action_0_observation_",
			format: formats.FullStateFormat{},
		},
		// "milana-vs-ruby-val20": modelDesc{
		// 	prefix: "serve_",
		// 	format: formats.FullStateCringeFormat{},
		// },
		// "milana-vs-ruby-noval": modelDesc{
		// 	prefix: "serve_",
		// 	format: formats.FullStateCringeFormat{},
		// },
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

	players := map[string]match.Player{
		"random":   &bot.RandomBot{},
		"minimax2": minimax.NewBot(minimax.SequentialRunner, 2),
		"minimax4": minimax.NewBot(minimax.SequentialRunner, 4),
		// "minimax6":    minimax.NewBot(minimax.SequentialRunner, 6),
		// "minimax8":    minimax.NewBot(minimax.SequentialRunner, 8),
		"alphabeta2": alphabeta2.NewBot(2),
		"alphabeta4": alphabeta2.NewBot(4),
		// "alphabeta6":  alphabeta2.NewBot(6),
		// "alphabeta8":  alphabeta2.NewBot(8),
		// "alphabeta10": alphabeta2.NewBot(10),
		// "mlval20":     mlbot.NewBot(models["ruby-vs-milana-val20"]),
		// "mlnoval":     mlbot.NewBot(models["ruby-vs-milana-noval"]),
		"mldqn": mlbot.NewBot(models["policy_10_51"]),
	}

	return players
}

type collector struct {
	states []match.GameState
	bot    match.Player
}

func (c *collector) SendState(ctx context.Context, state match.GameState) error {
	c.states = append(c.states, state)
	return c.bot.SendState(ctx, state)
}

func (c *collector) SendError(ctx context.Context, err error) error {
	return c.bot.SendError(ctx, err)
}

func (c *collector) SendEnd(ctx context.Context) error {
	return c.bot.SendEnd(ctx)
}

func (c *collector) RequestSkill(ctx context.Context) (int, error) {
	return c.bot.RequestSkill(ctx)
}

func (c *collector) GivenUp() <-chan any {
	return nil
}

func states() []match.GameState {
	var states []match.GameState

	for range 100 {
		hmm := rand.IntN(2) == 1

		c := game.NewCharacter(milana.CharacterMilana)
		opp := game.NewCharacter(ruby.CharacterRuby)

		p := &collector{
			bot: &bot.RandomBot{},
		}
		b := &bot.RandomBot{}

		m := match.New(match.CharacterPlayer{
			Character: c,
			Player:    p,
		}, match.CharacterPlayer{
			Character: opp,
			Player:    b,
		}, hmm)

		m.Run(context.Background())

		states = append(states, p.states...)
	}

	return states
}
