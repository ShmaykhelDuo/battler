package ml

import "github.com/ShmaykhelDuo/battler/internal/game/match"

type Format interface {
	Row(state match.GameState) map[string]any
}
