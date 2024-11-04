package ml

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/ShmaykhelDuo/battler/internal/bot/minimax"
)

var (
	start = []byte("[")
	sep   = []byte(",")
	end   = []byte("]")
)

func ExportDataset(out io.Writer, data []minimax.Entry, format Format) error {
	_, err := out.Write(start)
	if err != nil {
		return fmt.Errorf("write start: %w", err)
	}

	enc := json.NewEncoder(out)

	for i, e := range data {
		if i != 0 {
			_, err := out.Write(sep)
			if err != nil {
				return fmt.Errorf("write sep: %w", err)
			}
		}

		res := format.Row(e.State)
		res["result"] = TensorableValue[int64]{int64(e.Result[e.State.TurnState][0])}
		err := enc.Encode(res)
		if err != nil {
			return fmt.Errorf("encode: %w", err)
		}
	}

	_, err = out.Write(end)
	if err != nil {
		return fmt.Errorf("write end: %w", err)
	}

	return nil
}
