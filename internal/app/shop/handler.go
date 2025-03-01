package shop

import service "github.com/ShmaykhelDuo/battler/internal/service/shop"

type Handler struct {
	s *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{s: s}
}
