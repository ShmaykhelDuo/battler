package botfactory

import (
	"fmt"
	"path/filepath"

	"github.com/ShmaykhelDuo/battler/internal/bot/ml/bot"
	"github.com/ShmaykhelDuo/battler/internal/bot/ml/formats"
	"github.com/ShmaykhelDuo/battler/internal/bot/ml/model"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type CharacterRepository interface {
	Characters() []int
}

type Factory struct {
	models map[int]map[int]*model.Model
}

func New(cr CharacterRepository, path string) (*Factory, error) {
	chars := cr.Characters()

	res := &Factory{
		models: make(map[int]map[int]*model.Model, len(chars)),
	}

	prefix := "action_0_observation_"
	format := formats.FullStateFormat{}

	for _, botChar := range chars {
		res.models[botChar] = make(map[int]*model.Model, len(chars)-1)

		for _, playerChar := range chars {
			if playerChar == botChar {
				continue
			}

			path := filepath.Join(path, fmt.Sprintf("policy_%d_%d", botChar, playerChar))

			var err error
			res.models[botChar][playerChar], err = model.LoadModel(path, format, prefix)
			if err != nil {
				return nil, fmt.Errorf("load model %s: %w", path, err)
			}
		}
	}

	return res, nil
}

func (f *Factory) Bot(botChar, playerChar int) (match.Player, error) {
	m, ok := f.models[botChar]
	if !ok {
		return nil, fmt.Errorf("invalid bot character: %d", botChar)
	}

	model, ok := m[playerChar]
	if !ok {
		return nil, fmt.Errorf("invalid player character: %d", playerChar)
	}

	return bot.NewBot(model), nil
}
