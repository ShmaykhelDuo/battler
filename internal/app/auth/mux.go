package auth

import "net/http"

func Mux(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", h.Register)
	mux.HandleFunc("POST /signin", h.SignIn)
	mux.HandleFunc("POST /signout", h.SignOut)

	return mux
}
