package ml

import (
	"encoding/json"
	"io"

	"github.com/ShmaykhelDuo/battler/internal/bot/minimax"
)

func ExportDataset(out io.Writer, data []minimax.Entry, format Format) error {
	res := make([]map[string]Tensorable, len(data))

	for i, e := range data {
		res[i] = format.Row(e.State)
		res[i]["result"] = TensorableValue[int64]{int64(e.Result[e.State.TurnState][0])}
	}

	return json.NewEncoder(out).Encode(res)
}
