package main

import (
	"math/rand/v2"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/euphoria"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/milana"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/ruby"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/speed"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/storyteller"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/structure"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/z89"
)

var chars = []*game.CharacterData{
	storyteller.CharacterStoryteller,
	z89.CharacterZ89,
	euphoria.CharacterEuphoria,
	ruby.CharacterRuby,
	speed.CharacterSpeed,
	milana.CharacterMilana,
	structure.CharacterStructure,
}

func getRandomPair() (c, opp *game.Character) {
	n := len(chars)
	i1 := rand.IntN(n)
	i2 := rand.IntN(n - 1)
	if i2 >= i1 {
		i2 += 1
	}

	c = game.NewCharacter(chars[i1])
	opp = game.NewCharacter(chars[i2])
	return
}

func getRandomChar(except game.CharacterData) *game.Character {
	for {
		i := rand.IntN(len(chars))
		if chars[i].Desc == except.Desc {
			continue
		}
		return game.NewCharacter(chars[i])
	}
}
