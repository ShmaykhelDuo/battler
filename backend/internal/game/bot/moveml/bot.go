package moveml

import (
	"errors"
	"fmt"
	"sort"

	"github.com/ShmaykhelDuo/battler/backend/internal/game/match"
)

type Bot struct {
	model   *Model
	actions []int
}

func NewBot(model *Model) *Bot {
	return &Bot{model: model}
}

func (b *Bot) SendState(state match.GameState) error {
	if !state.PlayerTurn {
		return nil
	}

	res, err := b.model.Predict(NewState(state))
	if err != nil {
		return fmt.Errorf("model: %w", err)
	}

	b.actions = scoresToRating(res)
	return nil
}

func (b *Bot) SendError() error {
	if len(b.actions) == 0 {
		return errors.New("no actions left")
	}

	b.actions = b.actions[1:]
	return nil
}

func (b *Bot) SendEnd() error {
	return nil
}

func (b *Bot) RequestSkill() (int, error) {
	return b.actions[0], nil
}

type sortableIndexSlice struct {
	sort.Interface
	idx []int
}

func newSortableIndexSlice(sortable sort.Interface) *sortableIndexSlice {
	s := &sortableIndexSlice{
		Interface: sortable,
		idx:       make([]int, sortable.Len()),
	}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}

func (s *sortableIndexSlice) Swap(i, j int) {
	s.Interface.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

func scoresToRating(scores []float64) []int {
	s := newSortableIndexSlice(sort.Reverse(sort.Float64Slice(scores)))
	sort.Sort(s)
	return s.idx
}
