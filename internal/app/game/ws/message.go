package ws

import (
	"encoding/json"
	"errors"

	"github.com/ShmaykhelDuo/battler/internal/app/game/ws/gamestate"
)

type MessageType int

const (
	MessageTypeUnknown        MessageType = 0
	MessageTypeMatchRequest   MessageType = 1
	MessageTypeMatchReconnect MessageType = 2
	MessageTypeError          MessageType = 3
	MessageTypeGameState      MessageType = 4
	MessageTypeMove           MessageType = 5
	MessageTypeGiveUp         MessageType = 6
	MessageTypeGameEnd        MessageType = 7
)

type MessageMatchRequest struct {
	MainCharacter      int `json:"main"`
	SecondaryCharacter int `json:"secondary"`
}

type MessageMatchReconnect struct {
}

type MessageError struct {
	ID      int    `json:"id"`
	Kind    string `json:"kind"`
	Message string `json:"message,omitempty"`
}

type MessageGameState struct {
	State gamestate.GameState `json:"state"`
}

type MessageMove struct {
	Move int `json:"move"`
}

type MessageGiveUp struct {
}

type MessageGameEnd struct {
	Result         int `json:"result"`
	PrevLevel      int `json:"prev_level"`
	PrevExperience int `json:"prev_experience"`
	Level          int `json:"level"`
	Experience     int `json:"experience"`
	Reward         int `json:"reward"`
}

var ErrInvalidMessageType = errors.New("invalid message type")

type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func NewMessage(payload any) (Message, error) {
	var t MessageType
	switch payload.(type) {
	case MessageMatchRequest:
		t = MessageTypeMatchRequest
	case MessageMatchReconnect:
		t = MessageTypeMatchReconnect
	case MessageError:
		t = MessageTypeError
	case MessageGameState:
		t = MessageTypeGameState
	case MessageMove:
		t = MessageTypeMove
	case MessageGiveUp:
		t = MessageTypeGiveUp
	case MessageGameEnd:
		t = MessageTypeGameEnd
	default:
		return Message{}, ErrInvalidMessageType
	}

	p, err := json.Marshal(payload)
	if err != nil {
		return Message{}, err
	}

	return Message{
		Type:    t,
		Payload: p,
	}, nil
}

func (m Message) UnmarshalPayload() (any, error) {
	switch m.Type {
	case MessageTypeMatchRequest:
		return unmarshalPayload[MessageMatchRequest](m.Payload)
	case MessageTypeMatchReconnect:
		return unmarshalPayload[MessageMatchReconnect](m.Payload)
	case MessageTypeError:
		return unmarshalPayload[MessageError](m.Payload)
	case MessageTypeGameState:
		return unmarshalPayload[MessageGameState](m.Payload)
	case MessageTypeMove:
		return unmarshalPayload[MessageMove](m.Payload)
	case MessageTypeGiveUp:
		return unmarshalPayload[MessageGiveUp](m.Payload)
	case MessageTypeGameEnd:
		return unmarshalPayload[MessageGameEnd](m.Payload)
	default:
		return nil, ErrInvalidMessageType
	}
}

func unmarshalPayload[T any](payload []byte) (T, error) {
	var res T
	err := json.Unmarshal(payload, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}
