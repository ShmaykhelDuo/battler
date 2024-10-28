package match

import (
	"slices"

	"github.com/ShmaykhelDuo/battler/internal/game"
)

// SkillLog is a log of all skills used.
type SkillLog map[game.TurnState][]int

// Clones returns a deep clone of the log.
func (l SkillLog) Clone() SkillLog {
	res := make(SkillLog, len(l))

	for state, skills := range l {
		res[state] = slices.Clone(skills)
	}

	return res
}

// GameState is a representation of a state of a game.
type GameState struct {
	Character  *game.Character
	Opponent   *game.Character
	TurnState  game.TurnState
	SkillsLeft int
	SkillLog   SkillLog
	PlayerTurn bool
	AsOpp      bool
}

// Clone returns a clone of the state.
func (s GameState) Clone() GameState {
	return GameState{
		Character:  s.Character.Clone(),
		Opponent:   s.Opponent.Clone(),
		TurnState:  s.TurnState,
		SkillsLeft: s.SkillsLeft,
		SkillLog:   s.SkillLog.Clone(),
		PlayerTurn: s.PlayerTurn,
		AsOpp:      s.AsOpp,
	}
}
