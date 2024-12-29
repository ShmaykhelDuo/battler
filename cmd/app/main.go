package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	authhandler "github.com/ShmaykhelDuo/battler/internal/app/auth"
	gamehandler "github.com/ShmaykhelDuo/battler/internal/app/game"
	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
	authhttp "github.com/ShmaykhelDuo/battler/internal/pkg/auth/http"
	"github.com/ShmaykhelDuo/battler/internal/pkg/character"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db/postgres"
	"github.com/ShmaykhelDuo/battler/internal/pkg/passhash/bcrypt"
	"github.com/ShmaykhelDuo/battler/internal/repository/auth/session"
	"github.com/ShmaykhelDuo/battler/internal/repository/auth/user"
	"github.com/ShmaykhelDuo/battler/internal/repository/game/available"
	characterrepo "github.com/ShmaykhelDuo/battler/internal/repository/game/character"
	authservice "github.com/ShmaykhelDuo/battler/internal/service/auth"
	"github.com/ShmaykhelDuo/battler/internal/service/game"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

func main() {
	err := run()
	if err != nil {
		slog.Error("fatal error occurred", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(h))

	mux, err := constructDependencies(ctx)
	if err != nil {
		return fmt.Errorf("construct dependencies: %w", err)
	}

	serv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		err := serv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server listen: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		err = serv.Shutdown(shutdownCtx)
		if err != nil {
			return fmt.Errorf("server shutdown: %w", err)
		}

		slog.Info("server was shut down successfully")
		return nil
	})

	slog.Info("server is started", "addr", serv.Addr)

	return eg.Wait()
}

func constructDependencies(ctx context.Context) (http.Handler, error) {
	connString := os.Getenv("DB_CONN")

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("pgxpool: %w", err)
	}

	db := postgres.NewDB(pool)
	tm := postgres.NewTransactionManager(pool)

	userRepo := user.NewPostgresRepository(db)
	sessionRepo := session.NewInMemoryRepository()
	availableCharRepo := available.NewPostgresRepository(db)
	characterRepo := characterrepo.NewGameRepository()

	passwordHasher, err := bcrypt.NewPasswordHasher(10)
	if err != nil {
		return nil, fmt.Errorf("bcrypt: %w", err)
	}

	characterPicker := character.NewPicker(characterRepo)

	authService := authservice.NewService(userRepo, sessionRepo, passwordHasher)
	authHandler := authhandler.NewHandler(authService)

	gameService := game.NewService(availableCharRepo, characterPicker, tm)
	gameHandler := gamehandler.NewHandler(gameService)

	authMiddleware := authhttp.NewAuthMiddleware(sessionRepo)

	mux := http.NewServeMux()
	mux.Handle("/auth/", http.StripPrefix("/auth", authhandler.Mux(authHandler)))
	mux.Handle("/game/", http.StripPrefix("/game", gamehandler.Mux(gameHandler)))

	return api.PanicHandlerMiddleware(authMiddleware.Middleware(mux)), nil
}
