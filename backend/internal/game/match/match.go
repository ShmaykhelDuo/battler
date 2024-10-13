package match

import (
	"math/rand/v2"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
)

func Match(c1, c2 *game.Character, p1, p2 Player) (int, error) {
	goingFirst := rand.IntN(2) == 0

	if !goingFirst {
		c1, c2 = c2, c1
		p1, p2 = p2, p1
	}

	skillLog := NewSkillLog()

	var gameCtx game.Context
	for turnNum := game.MinTurnNumber; turnNum <= game.MaxTurnNumber; turnNum++ {
		gameCtx = game.TurnCtx(turnNum)

		gameCtx = gameCtx.WithGoingFirst(true)
		end, err := Turn(p1, p2, c1, c2, gameCtx, skillLog)
		if err != nil {
			return 0, err
		}
		if end {
			break
		}

		gameCtx = gameCtx.WithGoingFirst(false)
		end, err = Turn(p2, p1, c2, c1, gameCtx, skillLog)
		if err != nil {
			return 0, err
		}
		if end {
			break
		}
	}

	err := sendState(p1, c1, c2, gameCtx.WithTurnEnd(), skillLog, false, false)
	if err != nil {
		return 0, err
	}
	err = sendState(p2, c2, c1, gameCtx.WithTurnEnd(), skillLog, false, false)
	if err != nil {
		return 0, err
	}
	err = p1.SendEnd()
	if err != nil {
		return 0, err
	}
	err = p2.SendEnd()
	if err != nil {
		return 0, err
	}

	res := 0
	if c1.HP() < c2.HP() {
		res = 1
	} else if c2.HP() < c1.HP() {
		res = -1
	}

	if !goingFirst {
		res *= -1
	}

	return res, nil
}
