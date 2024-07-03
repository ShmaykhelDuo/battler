package main

import (
	"context"
	"github.com/ShmaykhelDuo/battler/backend/internal/router"
	"github.com/ShmaykhelDuo/battler/backend/internal/serv"
	"log/slog"
	"os"
	"os/signal"
	"time"
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
