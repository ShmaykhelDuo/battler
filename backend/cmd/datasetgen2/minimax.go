package main

import (
	"fmt"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/bot/ml"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/match"
)

func MiniMax(c, opp *game.Character, turnState game.TurnState, skillsLeft int, depth int, asOpp bool, out chan<- []int) (score int, strategy []int) {
	// c - кому делаем хорошо
	// opp - кому делаем плохо
	// asOpp - если сейчас ход противника

	if skillsLeft == 0 {
		endCtx := turnState
		endCtx.IsTurnEnd = true
		c.OnTurnEnd(opp, endCtx)
		opp.OnTurnEnd(c, endCtx)

		nextCtx := turnState.AddTurns(0, true)
		if asOpp {
			depth -= 1
		}
		return MiniMax(c, opp, nextCtx, opp.SkillsPerTurn(), depth, !asOpp, out)
	}

	if depth == 0 || hasGameEnded(c, opp, turnState) {
		score = c.HP() - opp.HP()
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
	appropriate := make([]bool, 4)
	for i, s := range playC.Skills() {
		appropriate[i] = s.IsAppropriate(playOpp, turnState)
	}
	filterAppropriate := appropriate[0] || appropriate[1] || appropriate[2] || appropriate[3]

	for i, s := range playC.Skills() {
		if depth > 7 {
			fmt.Printf("depth:%d, i:%d\n", depth, i)
		}
		if !s.IsAvailable(playOpp, turnState) {
			continue
		}
		if filterAppropriate && !c.IsControlledByOpp() && !s.IsAppropriate(playOpp, turnState) {
			continue
		}

		clonedC, clonedOpp := game.Clone(c, opp)

		var clonedPlayC, clonedPlayOpp *game.Character
		if asOpp {
			clonedPlayC = clonedOpp
			clonedPlayOpp = clonedC
		} else {
			clonedPlayC = clonedC
			clonedPlayOpp = clonedOpp
		}

		clonedS := clonedPlayC.Skills()[i]
		clonedS.Use(clonedPlayOpp, turnState)

		state := ml.NewState(match.GameState{
			Character:  c,
			Opponent:   opp,
			TurnState:  turnState,
			PlayerTurn: true,
			AsOpp:      asOpp,
		})
		skillScore, skillStrategy := MiniMax(clonedC, clonedOpp, turnState, skillsLeft-1, depth, asOpp, out)
		if !asOpp && len(skillStrategy) > 0 {
			out <- append(state.ToSlice(), skillStrategy[0])
		}

		if (worst && skillScore < score) || (!worst && skillScore > score) {
			score = skillScore
			strategy = append([]int{i}, skillStrategy...)
		}
	}

	return
}

func hasGameEnded(c, opp *game.Character, turnState game.TurnState) bool {
	return turnState.TurnNum > game.MaxTurnNumber || c.HP() <= 0 || opp.HP() <= 0
}
