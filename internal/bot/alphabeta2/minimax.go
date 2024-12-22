package alphabeta2

import (
	"context"

	"github.com/ShmaykhelDuo/battler/internal/bot/walk"
	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type State struct {
	state match.GameState
	depth int
}

// IsTerminal reports whether the node is terminal.
func (s State) IsTerminal() bool {
	return s.depth == 0 || s.state.IsEnd()
}

// IsParallel reports whether the walk needs to be done in parallel.
func (s State) IsParallel() bool {
	return false
	// return s.depth > r.MinParallelDepth && s.depth <= r.MaxParallelDepth
}

// Value returns result of the node.
func (s State) Value() Result {
	return Result{
		Score: s.state.Character.HP() - s.state.Opponent.HP(),
	}
}

// Children returns a slice of children nodes.
func (s State) Children() ([]walk.Move[State], error) {
	skills := s.skills(!s.wantsWorst())

	res := make([]walk.Move[State], len(skills))
	for i, skill := range skills {
		var err error
		res[i].Action = skill
		res[i].State, err = s.handleSkill(skill)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (s State) wantsWorst() bool {
	return s.state.AsOpp != s.state.Character.IsControlledByOpp()
}

func (s State) skills(filterAppropriate bool) []int {
	var playC, playOpp *game.Character
	if s.state.AsOpp {
		playC = s.state.Opponent
		playOpp = s.state.Character
	} else {
		playC = s.state.Character
		playOpp = s.state.Opponent
	}

	moves := make([]int, 0, 4)

	if filterAppropriate {
		for i, skill := range playC.Skills() {
			if skill.IsAppropriate(playC, playOpp, s.state.TurnState) {
				moves = append(moves, i)
			}
		}

		if len(moves) > 0 {
			return moves
		}
	}

	for i, skill := range playC.Skills() {
		if skill.IsAvailable(playC, playOpp, s.state.TurnState) {
			moves = append(moves, i)
		}
	}

	return moves
}

func (s State) handleSkill(skillNum int) (State, error) {
	state := s.state.CloneWithSkill(skillNum)
	depth := s.depth

	var clonedPlayC, clonedPlayOpp *game.Character
	if state.AsOpp {
		clonedPlayC = state.Opponent
		clonedPlayOpp = state.Character
	} else {
		clonedPlayC = state.Character
		clonedPlayOpp = state.Opponent
	}

	turnState := state.TurnState

	clonedS := clonedPlayC.Skills()[skillNum]
	err := clonedS.Use(clonedPlayC, clonedPlayOpp, turnState)
	if err != nil {
		return State{}, err
	}

	state.SkillLog.Append(turnState, skillNum)
	state.SkillsLeft--

	if state.SkillsLeft == 0 {
		c := state.Character
		opp := state.Opponent
		turnState := state.TurnState
		asOpp := state.AsOpp

		endCtx := turnState.WithTurnEnd()
		c.OnTurnEnd(opp, endCtx)
		opp.OnTurnEnd(c, endCtx)

		state.TurnState = turnState.Next()
		state.SkillsLeft = opp.SkillsPerTurn()
		state.AsOpp = !asOpp
		if asOpp {
			depth -= 1
		}
	}

	return State{
		state: state,
		depth: depth,
	}, nil
}

type Params struct {
	alpha int
	beta  int
}

// InitVars returns an initialized vars.
func (p Params) InitVars(s State) Vars {
	worst := s.wantsWorst()

	var score int
	if worst {
		score = 1000
	} else {
		score = -1000
	}

	return Vars{
		worst: worst,
		score: score,
		alpha: p.alpha,
		beta:  p.beta,
	}
}

type Vars struct {
	worst bool
	score int
	skill int
	alpha int
	beta  int
}

// Params returns params for nested walk call.
func (v Vars) Params() Params {
	return Params{
		alpha: v.alpha,
		beta:  v.beta,
	}
}

// Accumulate updates vars based on nested call result.
func (v Vars) Accumulate(action int, res Result) (update Vars, cutoff bool) {
	if (v.worst && res.Score < v.score) || (!v.worst && res.Score > v.score) {
		v.skill = action
		v.score = res.Score

		if !v.worst {
			if v.score > v.alpha {
				v.alpha = v.score
			}

			if v.score >= v.beta {
				return v, true
			}
		} else {
			if v.score < v.beta {
				v.beta = v.score
			}

			if v.score <= v.alpha {
				return v, true
			}
		}
	}

	return v, false
}

// Result returns accumulated result.
func (v Vars) Result() Result {
	return Result{
		Score: v.score,
		Skill: v.skill,
	}
}

type Result struct {
	Score int
	Skill int
}

func Minimax(ctx context.Context, state match.GameState, depth int) (Result, error) {
	r := walk.Runner[State, Params, Vars, Result]{}
	s := State{
		state: state,
		depth: depth,
	}
	p := Params{
		alpha: -1000,
		beta:  1000,
	}
	return r.Walk(ctx, s, p)
}
