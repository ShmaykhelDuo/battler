package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	apimodel "github.com/ShmaykhelDuo/battler/internal/model/api"
	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		api.HandleError(w, apimodel.Error{
			Kind:    apimodel.KindInvalidRequest,
			Message: "invalid body",
		})
		return
	}

	sessionID, err := h.s.Register(r.Context(), req.Username, req.Password)
	if err != nil {
		api.HandleError(w, fmt.Errorf("register: %w", err))
		return
	}

	h.setSessionCookie(w, sessionID)

	w.WriteHeader(http.StatusCreated)
	return
}
