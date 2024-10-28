package minimax

import (
	"context"
	"errors"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

var ErrNoAvailableMoves = errors.New("no available moves")

type MiniMaxEntry struct {
	state  match.GameState
	result match.SkillLog
}

func MiniMax(ctx context.Context, state match.GameState, depth int) (score int, strategy match.SkillLog, entries []MiniMaxEntry, err error) {
	// c - кому делаем хорошо
	// opp - кому делаем плохо
	// asOpp - если сейчас ход противника

	c := state.Character
	opp := state.Opponent
	turnState := state.TurnState
	skillsLeft := state.SkillsLeft
	asOpp := state.AsOpp

	if skillsLeft == 0 {
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

	if depth == 0 || hasGameEnded(c, opp, turnState) {
		score = c.HP() - opp.HP()
		strategy = make(match.SkillLog)
		return
	}

	// Делаем плохо, если ходит противник (сам или за нас)
	worst := asOpp != c.IsControlledByOpp()

	if worst {
		score = 1000
	} else {
		score = -1000
	}

	var playC, playOpp *game.Character
	if asOpp {
		playC = opp
		playOpp = c
	} else {
		playC = c
		playOpp = opp
	}

	moves := make([]int, 0, 4)

	if !worst {
		for i, s := range playC.Skills() {
			if s.IsAppropriate(playOpp, turnState) {
				moves = append(moves, i)
			}
		}
	}

	if len(moves) == 0 {
		for i, s := range playC.Skills() {
			if s.IsAvailable(playOpp, turnState) {
				moves = append(moves, i)
			}
		}
	}

	if len(moves) == 0 {
		err = ErrNoAvailableMoves
		return
	}

	for _, skillNum := range moves {
		newState := state.Clone()

		var clonedPlayC, clonedPlayOpp *game.Character
		if asOpp {
			clonedPlayC = newState.Opponent
			clonedPlayOpp = newState.Character
		} else {
			clonedPlayC = newState.Character
			clonedPlayOpp = newState.Opponent
		}

		clonedS := clonedPlayC.Skills()[skillNum]
		clonedS.Use(clonedPlayOpp, turnState)

		newState.SkillLog[turnState] = append(newState.SkillLog[turnState], skillNum)
		newState.SkillsLeft--

		var skillScore int
		var skillStrategy match.SkillLog
		var skillEntries []MiniMaxEntry
		skillScore, skillStrategy, skillEntries, err = MiniMax(ctx, newState, depth)
		if err != nil {
			return
		}

		if (worst && skillScore < score) || (!worst && skillScore > score) {
			score = skillScore

			skillStrategy[turnState] = append([]int{skillNum}, skillStrategy[turnState]...)
			strategy = skillStrategy

			if !worst {
				entries = skillEntries
			}
		}

		if worst {
			entries = append(entries, skillEntries...)
		}
	}

	newEntry := MiniMaxEntry{
		state:  state,
		result: strategy,
	}
	entries = append(entries, newEntry)

	return
}

func hasGameEnded(c, opp *game.Character, turnState game.TurnState) bool {
	return turnState.TurnNum > game.MaxTurnNumber || c.HP() <= 0 || opp.HP() <= 0
}

func makeClones(c *game.Character, n int) []*game.Character {
	res := make([]*game.Character, n)

	for i := range n {
		if i == 0 {
			res[i] = c
			continue
		}

		res[i] = c.Clone()
	}

	return res
}
