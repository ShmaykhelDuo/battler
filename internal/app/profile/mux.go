package profile

import "net/http"

func Mux(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", h.Profile)

	return mux
}
