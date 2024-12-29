package game

import (
	"fmt"
	"net/http"

	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
)

func (h *Handler) UnlockInitialCharacters(w http.ResponseWriter, r *http.Request) {
	err := h.s.UnlockInitialCharacters(r.Context())
	if err != nil {
		api.HandleError(w, fmt.Errorf("characters: %w", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
