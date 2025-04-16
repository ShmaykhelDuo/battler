package character

import (
	"math/rand/v2"

	"github.com/ShmaykhelDuo/battler/internal/model/game"
)

type CharacterRepository interface {
	Characters() []int
	CharactersOfRarity(rarity game.CharacterRarity) []int
}

type Picker struct {
	cr CharacterRepository
}

func NewPicker(cr CharacterRepository) *Picker {
	return &Picker{cr: cr}
}

func (p *Picker) RandomCharacter() int {
	chars := p.cr.Characters()
	i := rand.IntN(len(chars))
	return chars[i]
}

func (p *Picker) RandomCharacters(n int) []int {
	chars := p.cr.Characters()
	perm := rand.Perm(len(chars))

	res := make([]int, n)
	for i := range n {
		res[i] = chars[perm[i]]
	}
	return res
}

func (p *Picker) RandomCharacterOfRarity(rarity game.CharacterRarity) int {
	chars := p.cr.CharactersOfRarity(rarity)
	i := rand.IntN(len(chars))
	return chars[i]
}
