package router

import "github.com/gorilla/mux"

// Module is a part of API, separated by path.
type Module interface {
	Path() string
	Register(r *mux.Router)
}
