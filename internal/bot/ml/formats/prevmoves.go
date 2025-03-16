package formats

import (
	"fmt"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type PrevMovesFormat struct {
}

func (f PrevMovesFormat) Row(state match.GameState) map[string]any {
	isGoingFirst := state.TurnState.IsGoingFirst

	res := map[string]any{
		"first": state.TurnState.IsGoingFirst,
		"asopp": state.AsOpp,
	}

	for turnNum := game.MinTurnNumber; turnNum <= game.MaxTurnNumber; turnNum++ {
		turnState := game.TurnState{
			TurnNum:      turnNum,
			IsGoingFirst: isGoingFirst,
		}
		name := fmt.Sprintf("skill%d", turnNum)
		res[name] = f.setTurn(turnState, state)

		turnState = game.TurnState{
			TurnNum:      turnNum,
			IsGoingFirst: !isGoingFirst,
		}
		name = fmt.Sprintf("oppskill%d", turnNum)
		res[name] = f.setTurn(turnState, state)
	}

	return res
}

func (f PrevMovesFormat) setTurn(turnState game.TurnState, state match.GameState) []int64 {
	skills, ok := state.SkillLog[turnState]
	if ok {
		val := make([]int64, len(skills))
		for i, s := range skills {
			val[i] = int64(s)
		}

		return val
	}

	return []int64{}
}
