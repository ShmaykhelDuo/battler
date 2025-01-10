package gamestate

import (
	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type Turn struct {
	First  []int `json:"first,omitempty"`
	Second []int `json:"second,omitempty"`
}

type SkillLog struct {
	Turns []Turn `json:"turns"`
}

func NewSkillLog(l match.SkillLog, turnState game.TurnState) SkillLog {
	res := SkillLog{
		Turns: make([]Turn, turnState.TurnNum),
	}

	for i := range turnState.TurnNum {
		ts := game.NewTurnState(i + 1)
		res.Turns[i].First = l[ts.WithGoingFirst(true)]
		res.Turns[i].Second = l[ts.WithGoingFirst(false)]
	}

	return res
}
