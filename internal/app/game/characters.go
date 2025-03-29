package game

import (
	"fmt"
	"net/http"

	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
)

type Character struct {
	Number          int `json:"number"`
	Level           int `json:"level"`
	LevelExperience int `json:"level_experience"`
	MatchCount      int `json:"match_count"`
	WinCount        int `json:"win_count"`
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
			Number:          c.Number,
			Level:           c.Level,
			LevelExperience: c.LevelExperience,
			MatchCount:      c.MatchCount,
			WinCount:        c.WinCount,
		}
	}

	api.WriteJSONResponse(w, http.StatusOK, dto)
}
