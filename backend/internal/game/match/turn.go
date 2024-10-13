package match

import (
	"github.com/ShmaykhelDuo/battler/backend/internal/game"
)

func Turn(p, oppP Player, c, opp *game.Character, gameCtx game.Context, skillLog *SkillLog) (end bool, err error) {
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
		err = sendState(controlPlayer, controlCharacter, observeCharacter, gameCtx, skillLog, true, asOpp)
		if err != nil {
			return true, err
		}
		err = sendState(observePlayer, observeCharacter, controlCharacter, gameCtx, skillLog, false, asOpp)
		if err != nil {
			return true, err
		}

		for {
			i, err := controlPlayer.RequestSkill()
			if err != nil {
				return true, err
			}

			err = c.Skills()[i].Use(opp, gameCtx)
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

	endCtx := gameCtx.WithTurnEnd()
	c.OnTurnEnd(opp, endCtx)
	opp.OnTurnEnd(opp, endCtx)
	return isEnd(c, opp), nil
}

func isEnd(c, opp *game.Character) bool {
	return c.HP() <= 0 || opp.HP() <= 0
}
