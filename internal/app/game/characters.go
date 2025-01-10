package game

import (
	"fmt"
	"net/http"

	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
)

type Character struct {
	Number int `json:"number"`
}

func (h *Handler) Characters(w http.ResponseWriter, r *http.Request) {
	chars, err := h.gameSvc.AvailableCharacters(r.Context())
	if err != nil {
		api.HandleError(w, fmt.Errorf("characters: %w", err))
		return
	}

	dto := make([]Character, len(chars))
	for i, c := range chars {
		dto[i] = Character{
			Number: c.Number,
		}
	}

	api.WriteJSONResponse(w, http.StatusOK, dto)
}
