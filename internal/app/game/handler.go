package game

import (
	service "github.com/ShmaykhelDuo/battler/internal/service/game"
	"github.com/ShmaykhelDuo/battler/internal/service/match"
	"github.com/gorilla/websocket"
)

type Handler struct {
	gameSvc  *service.Service
	matchSvc *match.Service
	upgrader websocket.Upgrader
}

func NewHandler(gameSvc *service.Service, matchSvc *match.Service) *Handler {
	return &Handler{
		gameSvc:  gameSvc,
		matchSvc: matchSvc,
	}
}
