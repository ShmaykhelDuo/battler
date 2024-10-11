package main

import (
	"math/rand/v2"

	"github.com/ShmaykhelDuo/battler/backend/internal/bot"
	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/ruby"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/structure"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/match"
)

type PlayerState struct {
	State match.GameState
	End   bool
	Win   bool
}

type Player struct {
	Out chan<- PlayerState
	In  <-chan int
}

// var oppCharacters = []game.CharacterData{
// 	storyteller.CharacterStoryteller,
// 	z89.CharacterZ89,
// 	euphoria.CharacterEuphoria,
// 	ruby.CharacterRuby,
// 	speed.CharacterSpeed,
// 	milana.CharacterMilana,
// 	structure.CharacterStructure,
// }

// func selectRandomCharacter() game.CharacterData {
// 	i := rand.IntN(len(oppCharacters))
// 	return oppCharacters[i]
// }

func Game(player, bot Player) {
	c := game.NewCharacter(structure.CharacterStructure)
	opp := game.NewCharacter(ruby.CharacterRuby)

	goingFirst := rand.IntN(2) == 0

	var p1, p2 Player
	var c1, c2 *game.Character

	if goingFirst {
		p1 = player
		c1 = c
		p2 = bot
		c2 = opp
	} else {
		p1 = bot
		c1 = opp
		p2 = player
		c2 = c
	}

	for turnNum := 1; turnNum <= 10; turnNum++ {
		if Turn(p1, p2, c1, c2, game.Context{
			TurnNum:      turnNum,
			IsGoingFirst: true,
		}) {
			return
		}
		if Turn(p2, p1, c2, c1, game.Context{
			TurnNum:      turnNum,
			IsGoingFirst: false,
		}) {
			return
		}
	}

	endCtx := game.Context{
		TurnNum:      10,
		IsGoingFirst: false,
		IsTurnEnd:    true,
	}
	if c.HP() > opp.HP() {
		sendEnd(player, bot, c, opp, endCtx)
	} else {
		sendEnd(bot, player, opp, c, endCtx)
	}
}

func Turn(p, oppP Player, c, opp *game.Character, gameCtx game.Context) bool {
	controlPlayer := p
	sendC := c
	sendOpp := opp
	if c.IsControlledByOpp() {
		controlPlayer = oppP
		sendC = opp
		sendOpp = c
	}

	for range c.SkillsPerTurn() {
		state := PlayerState{
			State: match.GameState{
				Character: sendC,
				Opponent:  sendOpp,
				Context:   gameCtx,
			},
		}
		controlPlayer.Out <- state
		i := <-controlPlayer.In
		c.Skills()[i].Use(opp, gameCtx)
		if checkEnd(p, oppP, c, opp, gameCtx) {
			return true
		}
	}

	endCtx := gameCtx
	endCtx.IsTurnEnd = true
	c.OnTurnEnd(opp, endCtx)
	opp.OnTurnEnd(opp, endCtx)

	return checkEnd(p, oppP, c, opp, gameCtx)
}

func checkEnd(p, oppP Player, c, opp *game.Character, gameCtx game.Context) bool {
	if opp.HP() <= 0 {
		sendEnd(p, oppP, c, opp, gameCtx)
		return true
	}

	if c.HP() <= 0 {
		sendEnd(oppP, p, opp, c, gameCtx)
		return true
	}

	return false
}

func sendEnd(winP, loseP Player, winC, loseC *game.Character, gameCtx game.Context) {
	winP.Out <- PlayerState{
		State: match.GameState{
			Character: winC,
			Opponent:  loseC,
			Context:   gameCtx,
		},
		End: true,
		Win: true,
	}
	close(winP.Out)
	loseP.Out <- PlayerState{
		State: match.GameState{
			Character: loseC,
			Opponent:  winC,
			Context:   gameCtx,
		},
		End: true,
		Win: false,
	}
	close(loseP.Out)
}

func Bot() Player {
	out := make(chan PlayerState)
	in := make(chan int)
	bot := bot.RandomBot{}

	p := Player{
		Out: out,
		In:  in,
	}

	go func() {
		for {
			s, ok := <-out
			if !ok {
				break
			}

			if !s.End {
				in <- bot.Skill(s.State)
			}
		}
	}()

	return p
}
