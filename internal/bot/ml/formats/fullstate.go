package formats

import (
	"github.com/ShmaykhelDuo/battler/internal/bot/ml"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type FullStateFormat struct {
}

func (f FullStateFormat) Row(state match.GameState) map[string]ml.Tensorable {
	return GetMapState(state)
}
