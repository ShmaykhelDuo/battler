package character

import (
	"maps"
	"slices"

	"github.com/ShmaykhelDuo/battler/internal/game/characters/euphoria"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/milana"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/ruby"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/speed"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/storyteller"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/structure"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/z89"
	"github.com/ShmaykhelDuo/battler/internal/model/errs"
	model "github.com/ShmaykhelDuo/battler/internal/model/game"
)

var characters = map[int]model.Character{
	1: {
		Character: storyteller.CharacterStoryteller,
		Rarity:    model.CharacterRarityLF,
	},
	8: {
		Character: z89.CharacterZ89,
		Rarity:    model.CharacterRarityAD,
	},
	9: {
		Character: euphoria.CharacterEuphoria,
		Rarity:    model.CharacterRarityAD,
	},
	10: {
		Character: ruby.CharacterRuby,
		Rarity:    model.CharacterRaritySP,
	},
	33: {
		Character: speed.CharacterSpeed,
		Rarity:    model.CharacterRarityAD,
	},
	51: {
		Character: milana.CharacterMilana,
		Rarity:    model.CharacterRaritySP,
	},
	119: {
		Character: structure.CharacterStructure,
		Rarity:    model.CharacterRarityLF,
	},
}

type GameRepository struct {
}

func NewGameRepository() GameRepository {
	return GameRepository{}
}

func (r GameRepository) Characters() []int {
	return slices.Collect(maps.Keys(characters))
}

func (r GameRepository) CharactersOfRarity(rarity model.CharacterRarity) []int {
	var res []int

	for num, c := range characters {
		if c.Rarity == rarity {
			res = append(res, num)
		}
	}

	return res
}

func (r GameRepository) Character(number int) (model.Character, error) {
	char, ok := characters[number]
	if !ok {
		return model.Character{}, errs.ErrNotFound
	}

	return char, nil
}
