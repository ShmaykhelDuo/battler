package formats

import (
	"fmt"

	"github.com/ShmaykhelDuo/battler/internal/bot/ml"
	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type PrevMovesFormat struct {
}

func (f PrevMovesFormat) Row(state match.GameState) map[string]ml.Tensorable {
	isGoingFirst := state.TurnState.IsGoingFirst

	res := map[string]ml.Tensorable{
		"first": ml.TensorableValue[bool]{Item: state.TurnState.IsGoingFirst},
		"asopp": ml.TensorableValue[bool]{Item: state.AsOpp},
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

func (f PrevMovesFormat) setTurn(turnState game.TurnState, state match.GameState) ml.Tensorable {
	skills, ok := state.SkillLog[turnState]
	if ok {
		val := make([]int64, len(skills))
		for i, s := range skills {
			val[i] = int64(s)
		}

		return ml.TensorableSlice[int64]{Items: val}
	}

	return ml.TensorableSlice[int64]{Items: []int64{}}
}
