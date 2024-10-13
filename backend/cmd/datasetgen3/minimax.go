package main

import (
	"fmt"
	"slices"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
)

type Out struct {
	PrevMoves []int
	Strategy  []int
	First     bool
}

func MiniMax(c, opp *game.Character, gameCtx game.Context, skillsLeft int, depth int, asOpp bool, prevMoves []int, out chan<- Out) (score int, strategy []int) {
	// c - кому делаем хорошо
	// opp - кому делаем плохо
	// asOpp - если сейчас ход противника

	if skillsLeft == 0 {
		endCtx := gameCtx
		endCtx.IsTurnEnd = true
		c.OnTurnEnd(opp, endCtx)
		opp.OnTurnEnd(c, endCtx)

		nextCtx := gameCtx.AddTurns(0, true)
		if asOpp {
			depth -= 1
		}
		return MiniMax(c, opp, nextCtx, opp.SkillsPerTurn(), depth, !asOpp, prevMoves, out)
	}

	if depth == 0 || hasGameEnded(c, opp, gameCtx) {
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
		appropriate[i] = s.IsAppropriate(playOpp, gameCtx)
	}
	filterAppropriate := appropriate[0] || appropriate[1] || appropriate[2] || appropriate[3]

	for i, s := range playC.Skills() {
		if depth > 7 {
			fmt.Printf("depth:%d, i:%d\n", depth, i)
		}
		if !s.IsAvailable(playOpp, gameCtx) {
			continue
		}
		if filterAppropriate && !c.IsControlledByOpp() && !s.IsAppropriate(playOpp, gameCtx) {
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
		clonedS.Use(clonedPlayOpp, gameCtx)

		moves := slices.Clone(prevMoves)
		moves = append(moves, i)
		skillScore, skillStrategy := MiniMax(clonedC, clonedOpp, gameCtx, skillsLeft-1, depth, asOpp, moves, out)
		if !asOpp && len(skillStrategy) > 0 {
			o := Out{
				PrevMoves: prevMoves,
				Strategy:  skillStrategy,
				First:     gameCtx.IsGoingFirst,
			}
			out <- o
		}

		if (worst && skillScore < score) || (!worst && skillScore > score) {
			score = skillScore
			strategy = append([]int{i}, skillStrategy...)
		}
	}

	return
}

func hasGameEnded(c, opp *game.Character, gameCtx game.Context) bool {
	return gameCtx.TurnNum > game.MaxTurnNumber || c.HP() <= 0 || opp.HP() <= 0
}
