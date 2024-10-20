package match

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

type GameState struct {
	Character *game.Character
	Opponent  *game.Character
	Context   game.TurnState
}

type Player interface {
	Skill(state GameState) int
}

func Turn(p, oppP Player, c, opp *game.Character, turnState game.TurnState) {
	controlPlayer := p
	if c.IsControlledByOpp() {
		controlPlayer = oppP
	}

	for range c.SkillsPerTurn() {
		state := GameState{
			Character: c,
			Opponent:  opp,
			Context:   turnState,
		}
		i := controlPlayer.Skill(state)
		c.Skills()[i].Use(opp, turnState)
	}

	c.OnTurnEnd(opp, turnState)
	opp.OnTurnEnd(opp, turnState)
}
