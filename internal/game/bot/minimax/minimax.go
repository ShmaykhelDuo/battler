package minimax

import (
	"context"
	"errors"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
	"golang.org/x/sync/errgroup"
)

// ErrNoAvailableSkills is returned when no skills are available.
var ErrNoAvailableSkills = errors.New("no available skills")

// Entry is an entry used for further analysis.
type Entry struct {
	State  match.GameState // input game state
	Result match.SkillLog  // minimax result for provided state
}

// Result is a result of minimax.
type Result struct {
	Score    int            // estimated score
	Strategy match.SkillLog // estimated strategy
	Entries  []Entry        // accumulated entries
}

// Runner contains configuration required to run minimax.
type Runner struct {
	MinConcDepth int // minimum concurrent depth
	MaxConcDepth int // maximum concurrent depth
}

// SequentialRunner is a [Runner] configured to be run in a parent goroutine.
var SequentialRunner = Runner{}

// MemOptConcurrentRunner is a [Runner] configured to be run concurrently with the minimum viable peak memory usage.
var MemOptConcurrentRunner = Runner{
	MinConcDepth: 3,
	MaxConcDepth: 5,
}

// TimeOptConcurrentRunner is a [Runner] configured to be run concurrently in the minimum time.
var TimeOptConcurrentRunner = Runner{
	MinConcDepth: 3,
	MaxConcDepth: 7,
}

// MiniMax computes the most optimal strategy within number of turns equal to depth.
func (r Runner) MiniMax(ctx context.Context, state match.GameState, depth int) (Result, error) {
	if depth == 0 || state.IsEnd() {
		return Result{
			Score:    state.Character.HP() - state.Opponent.HP(),
			Strategy: make(match.SkillLog),
		}, nil
	}

	if state.SkillsLeft == 0 {
		return r.turnEnd(ctx, state, depth)
	}

	// Делаем плохо, если ходит противник (сам или за нас)
	worst := state.AsOpp != state.Character.IsControlledByOpp()

	skills := getSkills(state, !worst)
	if len(skills) == 0 {
		return Result{}, ErrNoAvailableSkills
	}

	skillResults, err := r.handleSkills(ctx, state, skills, depth)
	if err != nil {
		return Result{}, err
	}

	return r.accumulate(state, skills, skillResults, worst), nil
}

func (r Runner) turnEnd(ctx context.Context, state match.GameState, depth int) (Result, error) {
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
	return r.MiniMax(ctx, state, depth)
}

func getSkills(state match.GameState, filterAppropriate bool) []int {
	var playC, playOpp *game.Character
	if state.AsOpp {
		playC = state.Opponent
		playOpp = state.Character
	} else {
		playC = state.Character
		playOpp = state.Opponent
	}

	moves := make([]int, 0, 4)

	if filterAppropriate {
		for i, s := range playC.Skills() {
			if s.IsAppropriate(playC, playOpp, state.TurnState) {
				moves = append(moves, i)
			}
		}

		if len(moves) > 0 {
			return moves
		}
	}

	for i, s := range playC.Skills() {
		if s.IsAvailable(playC, playOpp, state.TurnState) {
			moves = append(moves, i)
		}
	}

	return moves
}

func (r Runner) handleSkills(ctx context.Context, state match.GameState, skills []int, depth int) ([]Result, error) {
	if depth > r.MinConcDepth && depth <= r.MaxConcDepth {
		return r.handleSkillsConcurrent(ctx, state, skills, depth)
	} else {
		return r.handleSkillsSequential(ctx, state, skills, depth)
	}
}

func (r Runner) handleSkillsConcurrent(ctx context.Context, state match.GameState, skills []int, depth int) ([]Result, error) {
	skillResults := make([]Result, len(skills))
	eg, egCtx := errgroup.WithContext(ctx)

	for i, skillNum := range skills {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			eg.Go(func() error {
				var err error
				skillResults[i], err = r.handleSkill(egCtx, state.CloneWithSkill(i), skillNum, depth)
				return err
			})
		}
	}

	err := eg.Wait()
	if err != nil {
		return nil, err
	}

	return skillResults, nil
}

func (r Runner) handleSkillsSequential(ctx context.Context, state match.GameState, skills []int, depth int) ([]Result, error) {
	skillResults := make([]Result, len(skills))

	for i, skillNum := range skills {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			var err error
			skillResults[i], err = r.handleSkill(ctx, state.CloneWithSkill(i), skillNum, depth)
			if err != nil {
				return nil, err
			}
		}
	}

	return skillResults, nil
}

func (r Runner) handleSkill(ctx context.Context, state match.GameState, skillNum int, depth int) (res Result, err error) {
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
	clonedS.Use(clonedPlayC, clonedPlayOpp, turnState)

	state.SkillLog[turnState] = append(state.SkillLog[turnState], skillNum)
	state.SkillsLeft--

	return r.MiniMax(ctx, state, depth)
}

func (r Runner) accumulate(state match.GameState, skills []int, results []Result, worst bool) Result {
	var res Result

	if worst {
		res.Score = 1000
	} else {
		res.Score = -1000
	}

	var selected int

	for i, skillRes := range results {
		if (worst && skillRes.Score < res.Score) || (!worst && skillRes.Score > res.Score) {
			selected = skills[i]

			res.Score = skillRes.Score
			res.Strategy = skillRes.Strategy

			if !worst {
				res.Entries = skillRes.Entries
			}
		}

		if worst {
			res.Entries = append(res.Entries, skillRes.Entries...)
		}
	}

	res.Strategy[state.TurnState] = append([]int{selected}, res.Strategy[state.TurnState]...)

	newEntry := Entry{
		State:  state,
		Result: res.Strategy,
	}
	res.Entries = append(res.Entries, newEntry)

	return res
}
