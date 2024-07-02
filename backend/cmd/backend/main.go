package main

import (
	"net/http"
	"os"

	"github.com/ShmaykhelDuo/battler/backend/internal/handler"
)

func main() {
	srv := http.Server{
		Addr:    os.Getenv("LISTEN_ADDR"),
		Handler: handler.Handler(),
	}
	srv.ListenAndServe()
}
