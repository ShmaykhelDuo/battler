package minimax

import "github.com/ShmaykhelDuo/battler/internal/game"

func MiniMax(c, opp *game.Character, turnState game.TurnState, skillsLeft int, depth int, asOpp bool) (score int, strategy []int) {
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
		return MiniMax(c, opp, nextCtx, opp.SkillsPerTurn(), depth, !asOpp)
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

	clonedC := makeClones(c, len(moves))
	clonedOpp := makeClones(opp, len(moves))

	var clonedPlayC, clonedPlayOpp []*game.Character
	if asOpp {
		clonedPlayC = clonedOpp
		clonedPlayOpp = clonedC
	} else {
		clonedPlayC = clonedC
		clonedPlayOpp = clonedOpp
	}

	for i, skillNum := range moves {
		clonedS := clonedPlayC[i].Skills()[skillNum]
		clonedS.Use(clonedPlayOpp[i], turnState)

		skillScore, skillStrategy := MiniMax(clonedC[i], clonedOpp[i], turnState, skillsLeft-1, depth, asOpp)

		if (worst && skillScore < score) || (!worst && skillScore > score) {
			score = skillScore
			strategy = append([]int{skillNum}, skillStrategy...)
		}
	}

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
