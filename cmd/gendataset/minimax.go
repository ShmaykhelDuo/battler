package main

import (
	"context"
	"slices"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"golang.org/x/sync/errgroup"
)

type Out struct {
	PrevMoves []int
	Strategy  []int
	First     bool
}

func MiniMax(ctx context.Context, c, opp *game.Character, turnState game.TurnState, skillsLeft int, depth int, asOpp bool, prevMoves []int) (score int, strategy []int, res []Out, err error) {
	// c - кому делаем хорошо
	// opp - кому делаем плохо
	// asOpp - если сейчас ход противника

	err = ctx.Err()
	if err != nil {
		return
	}

	if skillsLeft == 0 {
		return miniMaxEndOfTurn(ctx, c, opp, turnState, depth, asOpp, prevMoves)
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
		appropriate[i] = s.IsAppropriate(playC, playOpp, turnState)
	}
	filterAppropriate := appropriate[0] || appropriate[1] || appropriate[2] || appropriate[3]

	eg, egCtx := errgroup.WithContext(ctx)

	skillScores := make([]int, 4)
	skillStrategies := make([][]int, 4)
	results := make([][]Out, 4)

	for i, s := range playC.Skills() {
		err = ctx.Err()
		if err != nil {
			return
		}

		skillScores[i] = score

		if !s.IsAvailable(playC, playOpp, turnState) {
			continue
		}
		if filterAppropriate && !worst && !appropriate[i] {
			continue
		}

		clonedC := c.Clone()
		clonedOpp := opp.Clone()

		var clonedPlayC, clonedPlayOpp *game.Character
		if asOpp {
			clonedPlayC = clonedOpp
			clonedPlayOpp = clonedC
		} else {
			clonedPlayC = clonedC
			clonedPlayOpp = clonedOpp
		}

		clonedS := clonedPlayC.Skills()[i]
		clonedS.Use(clonedPlayC, clonedPlayOpp, turnState)

		moves := slices.Clone(prevMoves)
		moves = append(moves, i)

		f := func() error {
			var err error
			skillScores[i], skillStrategies[i], results[i], err = MiniMax(egCtx, clonedC, clonedOpp, turnState, skillsLeft-1, depth, asOpp, moves)
			return err
		}

		if depth > 4 && depth < 8 {
			eg.Go(f)
		} else {
			err = f()
			if err != nil {
				return
			}
		}
	}

	err = eg.Wait()
	if err != nil {
		return
	}

	for i := range 4 {
		if (worst && skillScores[i] < score) || (!worst && skillScores[i] > score) {
			score = skillScores[i]
			strategy = append([]int{i}, skillStrategies[i]...)
			if !worst {
				res = results[i]
			}
		}

		if worst {
			res = append(res, results[i]...)
		}
	}

	if !worst {
		o := Out{
			PrevMoves: prevMoves,
			Strategy:  strategy,
			First:     turnState.IsGoingFirst,
		}
		res = append(res, o)
	}

	return
}

func miniMaxEndOfTurn(ctx context.Context, c, opp *game.Character, turnState game.TurnState, depth int, asOpp bool, prevMoves []int) (score int, strategy []int, res []Out, err error) {
	endCtx := turnState.WithTurnEnd()
	c.OnTurnEnd(opp, endCtx)
	opp.OnTurnEnd(c, endCtx)

	nextCtx := turnState.AddTurns(0, true)
	if asOpp {
		depth -= 1
	}
	return MiniMax(ctx, c, opp, nextCtx, opp.SkillsPerTurn(), depth, !asOpp, prevMoves)
}

func hasGameEnded(c, opp *game.Character, turnState game.TurnState) bool {
	return turnState.TurnNum > game.MaxTurnNumber || c.HP() <= 0 || opp.HP() <= 0
}
