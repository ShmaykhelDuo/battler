package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/ShmaykhelDuo/battler/internal/router"
	"github.com/ShmaykhelDuo/battler/internal/serv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	modules := make([]router.Module, 0)

	srv := serv.New(
		os.Getenv("LISTEN_ADDR"),
		modules,
		5*time.Second,
	)
	err := srv.Run(ctx)
	if err != nil {
		slog.Error("error running server", "err", err)
	}
}
