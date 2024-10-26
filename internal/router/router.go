package router

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Router compiles handler from Module definitions.
func Router(modules []Module) http.Handler {
	r := mux.NewRouter()
	for _, m := range modules {
		s := r.PathPrefix(m.Path()).Subrouter()
		m.Register(s)
	}
	return r
}
