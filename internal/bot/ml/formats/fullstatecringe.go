package formats

import (
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type FullStateCringeFormat struct {
}

func (f FullStateCringeFormat) Row(state match.GameState) map[string]any {
	return GetMapStateCringe(state)
}
