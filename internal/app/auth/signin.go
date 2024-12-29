package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	apimodel "github.com/ShmaykhelDuo/battler/internal/model/api"
	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
)

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		api.HandleError(w, apimodel.Error{
			Kind:    apimodel.KindInvalidRequest,
			Message: "invalid body",
		})
		return
	}

	sessionID, err := h.s.SignIn(r.Context(), req.Username, req.Password)
	if err != nil {
		api.HandleError(w, fmt.Errorf("sign in: %w", err))
		return
	}

	h.setSessionCookie(w, sessionID)

	w.WriteHeader(http.StatusOK)
	return
}
