package router

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testModule struct {
	P   string
	Reg func(r *mux.Router)
}

func (m testModule) Path() string {
	return m.P
}

func (m testModule) Register(r *mux.Router) {
	m.Reg(r)
}

func handlerResponse(t *testing.T, h http.Handler, path string) string {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
	return rr.Body.String()
}

func TestRouter(t *testing.T) {
	modules := []Module{
		testModule{
			P: "/path1",
			Reg: func(r *mux.Router) {
				r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("hello world"))
				}).Methods("GET")
			},
		},
		testModule{
			P: "/path2",
			Reg: func(r *mux.Router) {
				r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("hello world x2"))
				}).Methods("GET")
			},
		},
	}

	h := Router(modules)
	assert.Equal(t, "hello world", handlerResponse(t, h, "/path1/test"))
	assert.Equal(t, "hello world x2", handlerResponse(t, h, "/path2/test"))
}
