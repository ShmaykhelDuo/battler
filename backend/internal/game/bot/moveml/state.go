package moveml

import (
	"fmt"

	"github.com/ShmaykhelDuo/battler/backend/internal/game/match"
)

type State struct {
	First     bool
	Skills    []int
	OppSkills []int
}

func NewState(state match.GameState) State {
	s := State{
		First: state.Context.IsGoingFirst,
	}

	ourAdd := 0
	oppAdd := 1
	if !s.First {
		oppAdd = 0
		ourAdd = 1
	}
	for i := range 10 {
		if i*2+ourAdd < len(state.SkillLog) {
			s.Skills = append(s.Skills, state.SkillLog[i*2+ourAdd].SkillIndex)
		}

		if i*2+oppAdd < len(state.SkillLog) {
			s.OppSkills = append(s.OppSkills, state.SkillLog[i*2+oppAdd].SkillIndex)
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

	skillsToMap(out, s.Skills, "skill")
	skillsToMap(out, s.OppSkills, "oppskill")

	return out
}

func skillsToMap(out map[string]int64, s []int, prefix string) {
	for i := range 9 {
		name := fmt.Sprintf("%s%d", prefix, i+1)
		var res int64 = -1
		if i < len(s) {
			res = int64(s[i])
		}
		out[name] = res
	}
}
