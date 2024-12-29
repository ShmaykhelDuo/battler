package auth

import (
	"fmt"
	"net/http"

	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
)

func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) {
	err := h.s.SignOut(r.Context())
	if err != nil {
		api.HandleError(w, fmt.Errorf("sign out: %w", err))
		return
	}

	h.removeSessionCookie(w)

	w.WriteHeader(http.StatusOK)
	return
}
