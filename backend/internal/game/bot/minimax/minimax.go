package minimax

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

func MiniMax(c, opp *game.Character, gameCtx game.Context, skillsLeft int, depth int, asOpp bool) (score int, strategy []int) {
	if skillsLeft == 0 {
		endCtx := gameCtx
		endCtx.IsTurnEnd = true
		c.OnTurnEnd(opp, endCtx)
		opp.OnTurnEnd(c, endCtx)

		nextCtx := gameCtx.AddTurns(0, true)
		return MiniMax(opp, c, nextCtx, opp.SkillsPerTurn(), depth-1, !asOpp)
	}

	if depth == 0 || hasGameEnded(c, opp, gameCtx) {
		score = c.HP() - opp.HP()
		return
	}

	worst := asOpp != c.IsControlledByOpp()

	if worst {
		score = 1000
	} else {
		score = -1000
	}

	for i, s := range c.Skills() {
		if !s.IsAvailable(opp, gameCtx) {
			continue
		}

		clonedC, clonedOpp := game.Clone(c, opp)

		clonedS := clonedC.Skills()[i]
		clonedS.Use(clonedOpp, gameCtx)

		skillScore, skillStrategy := MiniMax(clonedOpp, clonedC, gameCtx, skillsLeft-1, depth, asOpp)

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
