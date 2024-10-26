package moveml

import (
	"fmt"

	tf "github.com/wamuir/graft/tensorflow"
)

type Model struct {
	model *tf.SavedModel
}

func LoadModel(path string) (*Model, error) {
	model, err := tf.LoadSavedModel(path, []string{"serve"}, nil)
	if err != nil {
		return nil, fmt.Errorf("tensorflow: %w", err)
	}

	return &Model{model: model}, nil
}

func (m *Model) Predict(s State) ([]float64, error) {
	in := make(map[tf.Output]*tf.Tensor)

	for key, val := range s.toInputMap() {
		name := "serving_default_" + key
		output := m.model.Graph.Operation(name).Output(0)
		tensor, err := tf.NewTensor([1][1]int64{{int64(val)}})
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
