package money

import (
	"fmt"
	"net/http"

	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
)

func (h *Handler) Balance(w http.ResponseWriter, r *http.Request) {
	balance, err := h.s.Balance(r.Context())
	if err != nil {
		api.HandleError(w, fmt.Errorf("balance: %w", err))
		return
	}

	api.WriteJSONResponse(w, http.StatusOK, balance)
}
