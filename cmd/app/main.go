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
	friendhandler "github.com/ShmaykhelDuo/battler/internal/app/friends"
	gamehandler "github.com/ShmaykhelDuo/battler/internal/app/game"
	moneyhandler "github.com/ShmaykhelDuo/battler/internal/app/money"
	shophandler "github.com/ShmaykhelDuo/battler/internal/app/shop"
	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
	authhttp "github.com/ShmaykhelDuo/battler/internal/pkg/auth/http"
	"github.com/ShmaykhelDuo/battler/internal/pkg/character"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db/postgres"
	"github.com/ShmaykhelDuo/battler/internal/pkg/matchmaker"
	"github.com/ShmaykhelDuo/battler/internal/pkg/passhash/bcrypt"
	"github.com/ShmaykhelDuo/battler/internal/repository/auth/session"
	"github.com/ShmaykhelDuo/battler/internal/repository/auth/user"
	"github.com/ShmaykhelDuo/battler/internal/repository/game/available"
	characterrepo "github.com/ShmaykhelDuo/battler/internal/repository/game/character"
	connectionrepo "github.com/ShmaykhelDuo/battler/internal/repository/match/connection"
	balancerepo "github.com/ShmaykhelDuo/battler/internal/repository/money/balance"
	currencyconversionrepo "github.com/ShmaykhelDuo/battler/internal/repository/money/conversion"
	"github.com/ShmaykhelDuo/battler/internal/repository/shop/chest"
	friendrepo "github.com/ShmaykhelDuo/battler/internal/repository/social/friends"
	authservice "github.com/ShmaykhelDuo/battler/internal/service/auth"
	"github.com/ShmaykhelDuo/battler/internal/service/friends"
	"github.com/ShmaykhelDuo/battler/internal/service/game"
	"github.com/ShmaykhelDuo/battler/internal/service/match"
	"github.com/ShmaykhelDuo/battler/internal/service/money"
	"github.com/ShmaykhelDuo/battler/internal/service/shop"
	"github.com/ShmaykhelDuo/battler/web"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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

	mux, mm, err := constructDependencies(ctx)
	if err != nil {
		return fmt.Errorf("construct dependencies: %w", err)
	}

	serv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		err := mm.Run(ctx)
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	})

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

func constructDependencies(ctx context.Context) (http.Handler, *matchmaker.Matchmaker, error) {
	connString := os.Getenv("DB_CONN")

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, nil, fmt.Errorf("pgxpool: %w", err)
	}

	db := postgres.NewDB(pool)
	tm := postgres.NewTransactionManager(pool)

	cacheUrl := os.Getenv("CACHE_URL")
	opts, err := redis.ParseURL(cacheUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("redis parse url: %w", err)
	}
	redisCli := redis.NewClient(opts)

	userRepo := user.NewPostgresRepository(db)
	// sessionRepo := session.NewInMemoryRepository()
	sessionRepo := session.NewRedisRepository(redisCli)

	availableCharRepo := available.NewPostgresRepository(db)
	characterRepo := characterrepo.NewGameRepository()
	connectionRepo := connectionrepo.NewInMemoryRepository()

	balanceRepo := balancerepo.NewPostgresRepository(db)
	currencyConvRepo := currencyconversionrepo.NewPostgresRepository(db)

	chestRepo := chest.NewRepository()

	friendRepo := friendrepo.NewPostgresRepository(db)

	passwordHasher, err := bcrypt.NewPasswordHasher(10)
	if err != nil {
		return nil, nil, fmt.Errorf("bcrypt: %w", err)
	}

	characterPicker := character.NewPicker(characterRepo)

	authService := authservice.NewService(userRepo, sessionRepo, passwordHasher)
	authHandler := authhandler.NewHandler(authService)

	matchmaker := matchmaker.New(characterRepo)

	gameService := game.NewService(availableCharRepo, characterPicker, tm)
	matchService := match.NewService(connectionRepo, availableCharRepo, matchmaker, balanceRepo, tm)
	gameHandler := gamehandler.NewHandler(gameService, matchService)

	moneyService := money.NewService(balanceRepo, currencyConvRepo, tm)
	moneyHandler := moneyhandler.NewHandler(moneyService)

	shopService := shop.NewService(chestRepo, balanceRepo, characterPicker, availableCharRepo, tm)
	shopHandler := shophandler.NewHandler(shopService)

	friendService := friends.NewService(friendRepo)
	friendHandler := friendhandler.NewHandler(friendService)

	authMiddleware := authhttp.NewAuthMiddleware(sessionRepo)

	mux := http.NewServeMux()
	mux.Handle("/auth/", http.StripPrefix("/auth", authhandler.Mux(authHandler)))
	mux.Handle("/game/", http.StripPrefix("/game", gamehandler.Mux(gameHandler)))
	mux.Handle("/money/", http.StripPrefix("/money", moneyhandler.Mux(moneyHandler)))
	mux.Handle("/shop/", http.StripPrefix("/shop", shophandler.Mux(shopHandler)))
	mux.Handle("/friends/", http.StripPrefix("/friends", friendhandler.Mux(friendHandler)))

	mux.Handle("/web/", http.StripPrefix("/web", http.FileServerFS(web.FS)))

	return api.PanicHandlerMiddleware(authMiddleware.Middleware(mux)), matchmaker, nil
}
