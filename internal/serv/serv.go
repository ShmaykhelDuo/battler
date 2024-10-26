package serv

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/ShmaykhelDuo/battler/internal/router"
)

// Serv contains HTTP server boilerplate code and adds support of context.Context.
type Serv struct {
	addr            string
	modules         []router.Module
	shutdownTimeout time.Duration
}

func New(addr string, modules []router.Module, shutdownTimeout time.Duration) *Serv {
	return &Serv{
		addr:            addr,
		modules:         modules,
		shutdownTimeout: shutdownTimeout,
	}
}

func (s *Serv) Run(ctx context.Context) error {
	srv := &http.Server{
		Handler: router.Router(s.modules),
	}

	err := s.startServer(srv)
	if err != nil {
		return err
	}

	return s.awaitShutdown(ctx, srv)
}

func (s *Serv) startServer(srv *http.Server) error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	go func() {
		err := srv.Serve(l)
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("serve error", "err", err)
		}
	}()

	slog.Info("http server started", "addr", l.Addr())
	return nil
}

func (s *Serv) awaitShutdown(ctx context.Context, srv *http.Server) error {
	<-ctx.Done()
	slog.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	err := srv.Shutdown(shutdownCtx)
	if err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}
