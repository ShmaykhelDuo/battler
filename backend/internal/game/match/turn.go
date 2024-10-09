package match

import "github.com/ShmaykhelDuo/battler/backend/internal/game"

type GameState struct {
	Character *game.Character
	Opponent  *game.Character
	Context   game.Context
}

type Player interface {
	Skill(state GameState) int
}

func Turn(p, oppP Player, c, opp *game.Character, gameCtx game.Context) {
	controlPlayer := p
	if c.IsControlledByOpp() {
		controlPlayer = oppP
	}

	for range c.SkillsPerTurn() {
		state := GameState{
			Character: c,
			Opponent:  opp,
			Context:   gameCtx,
		}
		i := controlPlayer.Skill(state)
		c.Skills()[i].Use(opp, gameCtx)
	}

	c.OnTurnEnd(opp, gameCtx)
	opp.OnTurnEnd(opp, gameCtx)
}
