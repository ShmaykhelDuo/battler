package character

import (
	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/euphoria"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/milana"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/ruby"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/speed"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/storyteller"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/structure"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/z89"
)

var characters = []*game.CharacterData{
	storyteller.CharacterStoryteller,
	z89.CharacterZ89,
	euphoria.CharacterEuphoria,
	ruby.CharacterRuby,
	speed.CharacterSpeed,
	milana.CharacterMilana,
	structure.CharacterStructure,
}

type GameRepository struct {
}

func NewGameRepository() GameRepository {
	return GameRepository{}
}

func (r GameRepository) Characters() []int {
	res := make([]int, len(characters))
	for i, c := range characters {
		res[i] = c.Desc.Number
	}
	return res
}
