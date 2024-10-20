package ml2

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"

	"github.com/ShmaykhelDuo/battler/backend/internal/game/match"
)

type stateMsg struct {
	State []int `json:"state"`
}

type actionMsg struct {
	ActionScores []float64 `json:"actions"`
}

type Bot struct {
	encoder *json.Encoder
	decoder *json.Decoder
	state   match.GameState
	actions []int
	hasErr  bool
}

func NewBot() (*Bot, error) {
	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		return nil, err
	}

	return &Bot{
		encoder: json.NewEncoder(conn),
		decoder: json.NewDecoder(conn),
	}, nil
}

func (b *Bot) SendState(state match.GameState) error {
	b.state = state

	if !state.PlayerTurn {
		return nil
	}

	return b.send(state)
}

func (b *Bot) SendError() error {
	b.actions = b.actions[1:]
	b.hasErr = true
	return nil
}

func (b *Bot) SendEnd() error {
	return nil
}

func (b *Bot) RequestSkill() (int, error) {
	if b.hasErr {
		return b.actions[0], nil
	}

	var msg actionMsg
	err := b.decoder.Decode(&msg)
	if err != nil {
		return 0, fmt.Errorf("recv: %w", err)
	}
	fmt.Printf("Got action %#v\n", msg)
	b.actions = scoresToRating(msg.ActionScores)
	return b.actions[0], nil
}

func (b *Bot) send(state match.GameState) error {
	s := State{
		PrevMoves: nil,
		First:     state.TurnState.IsGoingFirst,
	}
	for _, it := range state.SkillLog {
		s.PrevMoves = append(s.PrevMoves, it.SkillIndex)
	}
	sMsg := stateMsg{
		State: s.ToSlice(),
	}
	fmt.Printf("Sending %#v\n", sMsg)
	err := b.encoder.Encode(sMsg)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}

	return nil
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
