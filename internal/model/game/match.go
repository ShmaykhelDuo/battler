package game

import (
	"github.com/ShmaykhelDuo/battler/internal/game/match"
	"github.com/google/uuid"
)

type MatchRequest struct {
	UserID          uuid.UUID
	Conn            match.Player
	Main, Secondary int
}

type PlayerID struct {
	IsBot  bool
	UserID uuid.UUID
}

type MatchPlayer struct {
	PlayerID  PlayerID
	Conn      match.Player
	Character int
}

type MatchPlayerEndResult struct {
	Result         match.ResultStatus
	PrevLevel      int
	PrevExperience int
	Level          int
	Experience     int
	Reward         int
}
