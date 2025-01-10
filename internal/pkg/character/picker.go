package character

import (
	"math/rand/v2"
)

type CharacterRepository interface {
	Characters() []int
}

type Picker struct {
	cr CharacterRepository
}

func NewPicker(cr CharacterRepository) *Picker {
	return &Picker{cr: cr}
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
