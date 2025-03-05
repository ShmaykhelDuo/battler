package shop

import "net/http"

func Mux(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /chests", h.Chests)
	mux.HandleFunc("POST /chests/{id}", h.BuyChest)

	return mux
}
