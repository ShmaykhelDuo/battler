package model

import (
	"fmt"

	"github.com/ShmaykhelDuo/battler/internal/bot/ml"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
	tf "github.com/wamuir/graft/tensorflow"
)

type Model struct {
	model  *tf.SavedModel
	format ml.Format
}

func LoadModel(path string, format ml.Format) (*Model, error) {
	model, err := tf.LoadSavedModel(path, []string{"serve"}, nil)
	if err != nil {
		return nil, fmt.Errorf("tensorflow: %w", err)
	}

	return &Model{
		model:  model,
		format: format,
	}, nil
}

func (m *Model) Predict(state match.GameState) ([]float64, error) {
	in := make(map[tf.Output]*tf.Tensor)

	for key, val := range m.format.Row(state) {
		name := "serving_default_" + key
		op := m.model.Graph.Operation(name)
		if op == nil {
			return nil, fmt.Errorf("no input is defined: %s", key)
		}

		output := op.Output(0)
		tensor, err := tf.NewTensor(val.Value())
		if err != nil {
			return nil, fmt.Errorf("tensor: %w", err)
		}

		in[output] = tensor
	}

	out := m.model.Graph.Operation("StatefulPartitionedCall_1").Output(0)

	res, err := m.model.Session.Run(in, []tf.Output{out}, nil)
	if err != nil {
		return nil, fmt.Errorf("run: %w", err)
	}

	vals := res[0].Value().([][]float32)[0]
	actRes := []float64{
		float64(vals[0]),
		float64(vals[1]),
		float64(vals[2]),
		float64(vals[3]),
	}
	return actRes, nil
}
