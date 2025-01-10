package game

import "net/http"

func Mux(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /characters", h.Characters)
	mux.HandleFunc("POST /characters/unlock", h.UnlockInitialCharacters)
	mux.HandleFunc("GET /match", h.StartMatch)

	return mux
}
