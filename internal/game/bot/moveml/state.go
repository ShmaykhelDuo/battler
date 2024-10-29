package moveml

import (
	"fmt"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type State struct {
	First     bool
	Skills    []int
	OppSkills []int
}

func NewState(state match.GameState) State {
	s := State{
		First: state.TurnState.IsGoingFirst,
	}

	for i := range 10 {
		turn := game.TurnState{
			TurnNum:      i,
			IsGoingFirst: s.First,
		}
		if skills, ok := state.SkillLog[turn]; ok {
			s.Skills = append(s.Skills, skills[0])
		}

		turn = game.TurnState{
			TurnNum:      i,
			IsGoingFirst: !s.First,
		}
		if skills, ok := state.SkillLog[turn]; ok {
			s.OppSkills = append(s.OppSkills, skills[0])
		}
	}
	return s
}

func (s State) toInputMap() map[string]int64 {
	out := make(map[string]int64)

	if s.First {
		out["first"] = 1
	} else {
		out["first"] = 0
	}

	skillsToMap(out, s.Skills, 9, "skill")
	skillsToMap(out, s.OppSkills, 10, "oppskill")

	return out
}

func skillsToMap(out map[string]int64, s []int, n int, prefix string) {
	for i := range n {
		name := fmt.Sprintf("%s%d", prefix, i+1)
		var res int64 = -1
		if i < len(s) {
			res = int64(s[i])
		}
		out[name] = res
	}
}
