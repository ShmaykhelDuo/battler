package friends

import "net/http"

func Mux(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", h.Friends)
	mux.HandleFunc("GET /incoming", h.IncomingFriendRequests)
	mux.HandleFunc("GET /outgoing", h.OutgoingFriendRequests)
	mux.HandleFunc("POST /{id}", h.CreateFriendLink)
	mux.HandleFunc("DELETE /{id}", h.RemoveFriendLink)
	mux.HandleFunc("GET /{id}", h.FriendshipStatus)

	return mux
}
