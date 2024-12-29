package game

import service "github.com/ShmaykhelDuo/battler/internal/service/game"

type Handler struct {
	s *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{s: s}
}
