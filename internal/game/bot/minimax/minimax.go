package minimax

import (
	"context"
	"errors"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

var ErrNoAvailableSkills = errors.New("no available moves")

type MiniMaxEntry struct {
	State  match.GameState
	Result match.SkillLog
}

type MiniMaxResult struct {
	Score    int
	Strategy match.SkillLog
	Entries  []MiniMaxEntry
}

func MiniMax(ctx context.Context, state match.GameState, depth int) (MiniMaxResult, error) {
	if depth == 0 || hasGameEnded(state) {
		return MiniMaxResult{
			Score:    state.Character.HP() - state.Opponent.HP(),
			Strategy: make(match.SkillLog),
		}, nil
	}

	if state.SkillsLeft == 0 {
		return miniMaxTurnEnd(ctx, state, depth)
	}

	// Делаем плохо, если ходит противник (сам или за нас)
	worst := state.AsOpp != state.Character.IsControlledByOpp()

	skills := getSkills(state, !worst)

	if len(skills) == 0 {
		return MiniMaxResult{}, ErrNoAvailableSkills
	}

	skillResults := make([]MiniMaxResult, len(skills))

	for i, skillNum := range skills {
		var err error
		skillResults[i], err = miniMaxSkill(ctx, state.Clone(), skillNum, depth)
		if err != nil {
			return MiniMaxResult{}, err
		}
	}

	return miniMaxAccumulate(state, skills, skillResults, worst), nil
}

func miniMaxTurnEnd(ctx context.Context, state match.GameState, depth int) (MiniMaxResult, error) {
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
	return MiniMax(ctx, state, depth)
}

func hasGameEnded(state match.GameState) bool {
	return state.TurnState.TurnNum > game.MaxTurnNumber || state.Character.HP() <= 0 || state.Opponent.HP() <= 0
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
			if s.IsAppropriate(playOpp, state.TurnState) {
				moves = append(moves, i)
			}
		}

		if len(moves) > 0 {
			return moves
		}
	}

	for i, s := range playC.Skills() {
		if s.IsAvailable(playOpp, state.TurnState) {
			moves = append(moves, i)
		}
	}

	return moves
}

func miniMaxSkill(ctx context.Context, state match.GameState, skillNum int, depth int) (res MiniMaxResult, err error) {
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
	clonedS.Use(clonedPlayOpp, turnState)

	state.SkillLog[turnState] = append(state.SkillLog[turnState], skillNum)
	state.SkillsLeft--

	return MiniMax(ctx, state, depth)
}

func miniMaxAccumulate(state match.GameState, skills []int, results []MiniMaxResult, worst bool) MiniMaxResult {
	var res MiniMaxResult

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

	newEntry := MiniMaxEntry{
		State:  state,
		Result: res.Strategy,
	}
	res.Entries = append(res.Entries, newEntry)

	return res
}
