package model

import (
	"fmt"
	"math"

	"github.com/ShmaykhelDuo/battler/internal/bot/ml"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
	tf "github.com/wamuir/graft/tensorflow"
)

type Model struct {
	model  *tf.SavedModel
	format ml.Format
	prefix string
}

func LoadModel(path string, format ml.Format, prefix string) (*Model, error) {
	model, err := tf.LoadSavedModel(path, []string{"serve"}, nil)
	if err != nil {
		return nil, fmt.Errorf("tensorflow: %w", err)
	}

	return &Model{
		model:  model,
		format: format,
		prefix: prefix,
	}, nil
}

func (m *Model) Predict(state match.GameState) (int, error) {
	in := make(map[tf.Output]*tf.Tensor)

	for key, val := range m.format.Row(state) {
		name := m.prefix + key
		op := m.model.Graph.Operation(name)
		if op == nil {
			return 0, fmt.Errorf("no input is defined: %s", key)
		}

		output := op.Output(0)
		tensor, err := NewTensor(val)
		if err != nil {
			return 0, fmt.Errorf("tensor: %w", err)
		}

		in[output] = tensor
	}

	// out := m.model.Graph.Operation("StatefulPartitionedCall_1").Output(0)
	out := m.model.Graph.Operation("StatefulPartitionedCall").Output(0)

	res, err := m.model.Session.Run(in, []tf.Output{out}, nil)
	if err != nil {
		return 0, fmt.Errorf("run: %w", err)
	}

	switch val := res[0].Value().(type) {
	case [][]float32:
		maxIndex := -1
		maxVal := float32(math.Inf(-1))

		for i, v := range val[0] {
			if v > maxVal {
				maxVal = v
				maxIndex = i
			}
		}

		return maxIndex, nil

	case []int64:
		return int(val[0]), nil

	default:
		return 0, fmt.Errorf("invalid output format: %T", val)
	}
}
