package match

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

type Player interface {
	SendState(state GameState) error
	SendError() error
	SendEnd() error
	RequestSkill() (int, error)
}

func sendState(p Player, c, opp *game.Character, turnState game.TurnState, skillLog *SkillLog, playerTurn bool, asOpp bool) error {
	state := GameState{
		Character:  c,
		Opponent:   opp,
		TurnState:  turnState,
		SkillLog:   skillLog.Items(),
		PlayerTurn: playerTurn,
		AsOpp:      asOpp,
	}
	return p.SendState(state)
}
