package money

import "net/http"

func Mux(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /balance", h.Balance)
	mux.HandleFunc("POST /conversion", h.Convert)
	mux.HandleFunc("GET /conversion", h.ConversionStatus)
	mux.HandleFunc("POST /conversion/claim", h.ClaimConversion)

	return mux
}
