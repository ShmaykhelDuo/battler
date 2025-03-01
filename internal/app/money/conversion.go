package money

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	apimodel "github.com/ShmaykhelDuo/battler/internal/model/api"
	"github.com/ShmaykhelDuo/battler/internal/model/money"
	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
)

type ConversionRequest struct {
	SourceCurrencyID int   `json:"currency_id"`
	SourceAmount     int64 `json:"amount"`
}

type Conversion struct {
	StartTime        time.Time `json:"started_at"`
	FinishTime       time.Time `json:"finishes_at"`
	TargetCurrencyID int       `json:"currency_id"`
	TargetAmount     int64     `json:"amount"`
}

func (h *Handler) Convert(w http.ResponseWriter, r *http.Request) {
	var req ConversionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		api.HandleError(w, apimodel.Error{
			Kind:    apimodel.KindInvalidRequest,
			Message: "invalid body",
		})
		return
	}

	conv, err := h.s.Convert(r.Context(), money.Currency(req.SourceCurrencyID), req.SourceAmount)
	if err != nil {
		api.HandleError(w, fmt.Errorf("convert: %w", err))
		return
	}

	api.WriteJSONResponse(w, http.StatusCreated, Conversion{
		StartTime:        conv.StartTime,
		FinishTime:       conv.FinishTime,
		TargetCurrencyID: int(conv.TargetCurrency),
		TargetAmount:     conv.TargetAmount,
	})
}

func (h *Handler) ConversionStatus(w http.ResponseWriter, r *http.Request) {
	conv, err := h.s.ConversionStatus(r.Context())
	if err != nil {
		api.HandleError(w, fmt.Errorf("conversion status: %w", err))
		return
	}

	api.WriteJSONResponse(w, http.StatusOK, Conversion{
		StartTime:        conv.StartTime,
		FinishTime:       conv.FinishTime,
		TargetCurrencyID: int(conv.TargetCurrency),
		TargetAmount:     conv.TargetAmount,
	})
}

func (h *Handler) ClaimConversion(w http.ResponseWriter, r *http.Request) {
	balance, err := h.s.ClaimConversion(r.Context())
	if err != nil {
		api.HandleError(w, fmt.Errorf("claim conversion: %w", err))
		return
	}

	api.WriteJSONResponse(w, http.StatusOK, balance)
}
