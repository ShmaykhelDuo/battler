package match

import (
	"github.com/ShmaykhelDuo/battler/backend/internal/game"
)

func Turn(p, oppP Player, c, opp *game.Character, turnState game.TurnState, skillLog *SkillLog) (end bool, err error) {
	controlPlayer := p
	observePlayer := oppP
	controlCharacter := c
	observeCharacter := opp

	asOpp := c.IsControlledByOpp()
	if asOpp {
		controlPlayer = oppP
		observePlayer = p
		controlCharacter = opp
		observeCharacter = c
	}

	skills := c.SkillsPerTurn()
	for range skills {
		err = sendState(controlPlayer, controlCharacter, observeCharacter, turnState, skillLog, true, asOpp)
		if err != nil {
			return true, err
		}
		err = sendState(observePlayer, observeCharacter, controlCharacter, turnState, skillLog, false, asOpp)
		if err != nil {
			return true, err
		}

		for {
			i, err := controlPlayer.RequestSkill()
			if err != nil {
				return true, err
			}

			// log.Printf("Player %s has selected skill %d\n", c.Desc().Name, i)

			err = c.Skills()[i].Use(opp, turnState)
			if err == nil {
				skillLog.Add(c, i)
				break
			}

			err = controlPlayer.SendError()
			if err != nil {
				return true, err
			}
		}

		if isEnd(c, opp) {
			return true, nil
		}
	}

	endState := turnState.WithTurnEnd()
	c.OnTurnEnd(opp, endState)
	opp.OnTurnEnd(opp, endState)
	return isEnd(c, opp), nil
}

func isEnd(c, opp *game.Character) bool {
	return c.HP() <= 0 || opp.HP() <= 0
}
