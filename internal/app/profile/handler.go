package profile

import (
	"net/http"

	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
	service "github.com/ShmaykhelDuo/battler/internal/service/profile"
	"github.com/google/uuid"
)

type ProfileStatistics struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	MatchCount int       `json:"match_count"`
	WinCount   int       `json:"win_count"`
}

type Handler struct {
	s *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	profile, err := h.s.Profile(r.Context())
	if err != nil {
		api.HandleError(w, err)
		return
	}

	res := ProfileStatistics{
		ID:         profile.ID,
		Username:   profile.Username,
		MatchCount: profile.MatchCount,
		WinCount:   profile.WinCount,
	}

	api.WriteJSONResponse(w, http.StatusCreated, res)
}
