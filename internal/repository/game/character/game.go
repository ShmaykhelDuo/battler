package character

import (
	"maps"
	"slices"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/euphoria"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/milana"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/ruby"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/speed"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/storyteller"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/structure"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/z89"
	"github.com/ShmaykhelDuo/battler/internal/model/errs"
)

var characters = map[int]*game.CharacterData{
	1:   storyteller.CharacterStoryteller,
	8:   z89.CharacterZ89,
	9:   euphoria.CharacterEuphoria,
	10:  ruby.CharacterRuby,
	33:  speed.CharacterSpeed,
	51:  milana.CharacterMilana,
	119: structure.CharacterStructure,
}

type GameRepository struct {
}

func NewGameRepository() GameRepository {
	return GameRepository{}
}

func (r GameRepository) Characters() []int {
	return slices.Collect(maps.Keys(characters))
}

func (r GameRepository) Character(number int) (*game.CharacterData, error) {
	char, ok := characters[number]
	if !ok {
		return nil, errs.ErrNotFound
	}

	return char, nil
}
