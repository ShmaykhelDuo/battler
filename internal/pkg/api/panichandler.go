package api

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func PanicHandlerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				HandleError(w, fmt.Errorf("panic occurred: %v: %s", r, debug.Stack()))
			}
		}()

		h.ServeHTTP(w, r)
	})
}
