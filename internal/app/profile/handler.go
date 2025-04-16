package profile

import (
	"net/http"

	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
	service "github.com/ShmaykhelDuo/battler/internal/service/profile"
	"github.com/google/uuid"
)

type Profile struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
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

	res := Profile{
		ID:       profile.ID,
		Username: profile.Username,
	}

	api.WriteJSONResponse(w, http.StatusCreated, res)
}
