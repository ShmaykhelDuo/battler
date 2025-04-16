package shop

import (
	"fmt"
	"net/http"
	"strconv"

	apimodel "github.com/ShmaykhelDuo/battler/internal/model/api"
	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
)

type Chest struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	PriceCurrencyID int    `json:"currency_id"`
	PriceAmount     int64  `json:"price"`
	Available       bool   `json:"available"`
}

type Character struct {
	Number int `json:"number"`
}

func (h *Handler) Chests(w http.ResponseWriter, r *http.Request) {
	chests, err := h.s.Chests(r.Context())
	if err != nil {
		api.HandleError(w, fmt.Errorf("chests: %w", err))
		return
	}

	res := make([]Chest, len(chests))
	for i, c := range chests {
		res[i] = Chest{
			ID:              c.ID,
			Name:            c.Name,
			PriceCurrencyID: int(c.PriceCurrency),
			PriceAmount:     c.PriceAmount,
			Available:       c.Available,
		}
	}

	api.WriteJSONResponse(w, http.StatusOK, res)
}

func (h *Handler) BuyChest(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		api.HandleError(w, apimodel.Error{
			Kind:    apimodel.KindNotFound,
			Message: "invalid chest id",
		})
		return
	}

	char, err := h.s.BuyChest(r.Context(), id)
	if err != nil {
		api.HandleError(w, fmt.Errorf("buy chest: %w", err))
		return
	}

	res := Character{
		Number: char.Number,
	}
	api.WriteJSONResponse(w, http.StatusOK, res)
}
