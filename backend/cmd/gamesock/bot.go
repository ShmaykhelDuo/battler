package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/ShmaykhelDuo/battler/backend/internal/game/bot/ml"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/match"
)

type stateMsg struct {
	State  []int `json:"state"`
	End    bool  `json:"end"`
	Reward int   `json:"reward"`
}

type actionMsg struct {
	Action int `json:"action"`
}

type DQLLearnerBot struct {
	encoder     *json.Encoder
	decoder     *json.Decoder
	state       match.GameState
	totalReward int
}

func NewDQLLearnerBot(conn net.Conn) *DQLLearnerBot {
	return &DQLLearnerBot{
		encoder: json.NewEncoder(conn),
		decoder: json.NewDecoder(conn),
	}
}

func (b *DQLLearnerBot) SendState(state match.GameState) error {
	b.state = state

	if !state.PlayerTurn {
		return nil
	}

	return b.send(state, false, false)
}

func (b *DQLLearnerBot) SendError() error {
	return b.send(b.state, false, true)
}

func (b *DQLLearnerBot) SendEnd() error {
	return b.send(b.state, true, false)
}

func (b *DQLLearnerBot) RequestSkill() (int, error) {
	var msg actionMsg
	err := b.decoder.Decode(&msg)
	if err != nil {
		return 0, fmt.Errorf("recv: %w", err)
	}
	fmt.Printf("Got action %#v\n", msg)
	return msg.Action, nil
}

func (b *DQLLearnerBot) send(state match.GameState, end bool, hasErr bool) error {
	reward := state.Character.HP() - state.Opponent.HP()
	// if end {
	// 	if reward > 0 {
	// 		reward += 10
	// 	} else if reward < 0 {
	// 		reward -= 10
	// 	}
	// }

	var sendReward int
	if hasErr {
		sendReward = 0
	} else {
		sendReward = reward - b.totalReward
		b.totalReward = reward
	}

	sMsg := stateMsg{
		State:  ml.NewStateV2(state).ToSlice(),
		End:    end,
		Reward: sendReward,
	}
	fmt.Printf("Sending %#v\n", sMsg)
	err := b.encoder.Encode(sMsg)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}

	return nil
}
