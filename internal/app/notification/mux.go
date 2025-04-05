package notification

import "net/http"

func Mux(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", h.ReceiveNotifications)

	return mux
}
