package formats

import (
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type FullStateFormat struct {
}

func (f FullStateFormat) Row(state match.GameState) map[string]any {
	return GetMapState(state)
}
