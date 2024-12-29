package auth

import (
	service "github.com/ShmaykhelDuo/battler/internal/service/auth"
)

type Handler struct {
	s *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{s: s}
}
