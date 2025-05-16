package game

import (
	"github.com/ShmaykhelDuo/battler/internal/game/match"
	"github.com/ShmaykhelDuo/battler/internal/model/social"
	"github.com/google/uuid"
)

type MatchRequestCharacter struct {
	Number int
	Level  int
}

type MatchRequest struct {
	UserID          uuid.UUID
	Conn            match.Player
	Main, Secondary MatchRequestCharacter
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
	Result          match.ResultStatus
	PrevLevel       int
	PrevExperience  int
	Level           int
	Experience      int
	Reward          int64
	OpponentProfile social.Profile
}
