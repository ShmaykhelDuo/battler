package auth

import "net/http"

func (h *Handler) setSessionCookie(w http.ResponseWriter, sessionID string) {
	sessionCookie := http.Cookie{
		Name:  "session_id",
		Value: sessionID,
		Path:  "/",
		// Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &sessionCookie)
}

func (h *Handler) removeSessionCookie(w http.ResponseWriter) {
	sessionCookie := http.Cookie{
		Name:   "session_id",
		Path:   "/",
		MaxAge: -1,
		// Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &sessionCookie)
}
