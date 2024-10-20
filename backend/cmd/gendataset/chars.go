package main

import (
	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/euphoria"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/milana"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/ruby"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/speed"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/storyteller"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/structure"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/z89"
)

var chars = map[int]game.CharacterData{
	1:   storyteller.CharacterStoryteller,
	8:   z89.CharacterZ89,
	9:   euphoria.CharacterEuphoria,
	10:  ruby.CharacterRuby,
	33:  speed.CharacterSpeed,
	51:  milana.CharacterMilana,
	119: structure.CharacterStructure,
}
