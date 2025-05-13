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
	notificationhandler "github.com/ShmaykhelDuo/battler/internal/app/notification"
	profilehandler "github.com/ShmaykhelDuo/battler/internal/app/profile"
	shophandler "github.com/ShmaykhelDuo/battler/internal/app/shop"
	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
	authhttp "github.com/ShmaykhelDuo/battler/internal/pkg/auth/http"
	"github.com/ShmaykhelDuo/battler/internal/pkg/botfactory"
	"github.com/ShmaykhelDuo/battler/internal/pkg/character"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db/postgres"
	"github.com/ShmaykhelDuo/battler/internal/pkg/matchmaker"
	"github.com/ShmaykhelDuo/battler/internal/pkg/passhash/bcrypt"
	"github.com/ShmaykhelDuo/battler/internal/repository/auth/session"
	"github.com/ShmaykhelDuo/battler/internal/repository/auth/user"
	"github.com/ShmaykhelDuo/battler/internal/repository/game/available"
	characterrepo "github.com/ShmaykhelDuo/battler/internal/repository/game/character"
	connectionrepo "github.com/ShmaykhelDuo/battler/internal/repository/match/connection"
	matchrepo "github.com/ShmaykhelDuo/battler/internal/repository/match/match"
	balancerepo "github.com/ShmaykhelDuo/battler/internal/repository/money/balance"
	currencyconversionrepo "github.com/ShmaykhelDuo/battler/internal/repository/money/conversion"
	notificationrepo "github.com/ShmaykhelDuo/battler/internal/repository/notification/notification"
	"github.com/ShmaykhelDuo/battler/internal/repository/shop/chest"
	friendrepo "github.com/ShmaykhelDuo/battler/internal/repository/social/friends"
	profilerepo "github.com/ShmaykhelDuo/battler/internal/repository/social/profile"
	authservice "github.com/ShmaykhelDuo/battler/internal/service/auth"
	"github.com/ShmaykhelDuo/battler/internal/service/friends"
	"github.com/ShmaykhelDuo/battler/internal/service/game"
	"github.com/ShmaykhelDuo/battler/internal/service/match"
	"github.com/ShmaykhelDuo/battler/internal/service/money"
	"github.com/ShmaykhelDuo/battler/internal/service/notification"
	"github.com/ShmaykhelDuo/battler/internal/service/profile"
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

	mux, ms, mm, err := constructDependencies(ctx)
	if err != nil {
		return fmt.Errorf("construct dependencies: %w", err)
	}

	serv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		err := ms.HandleMatches(ctx)
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	})

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

func constructDependencies(ctx context.Context) (http.Handler, *match.Service, *matchmaker.Matchmaker, error) {
	connString := os.Getenv("DB_CONN")

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("pgxpool: %w", err)
	}

	db := postgres.NewDB(pool)
	tm := postgres.NewTransactionManager(pool)

	cacheUrl := os.Getenv("CACHE_URL")
	opts, err := redis.ParseURL(cacheUrl)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("redis parse url: %w", err)
	}
	redisCli := redis.NewClient(opts)

	userRepo := user.NewPostgresRepository(db)
	// sessionRepo := session.NewInMemoryRepository()
	sessionRepo := session.NewRedisRepository(redisCli)

	availableCharRepo := available.NewPostgresRepository(db)
	characterRepo := characterrepo.NewGameRepository()
	connectionRepo := connectionrepo.NewInMemoryRepository()
	matchRepo := matchrepo.NewPostgresRepository(db)

	balanceRepo := balancerepo.NewPostgresRepository(db)
	currencyConvRepo := currencyconversionrepo.NewPostgresRepository(db)

	chestRepo := chest.NewRepository()

	friendRepo := friendrepo.NewPostgresRepository(db)
	profileRepo := profilerepo.NewPostgresRepository(db)

	notificationRepo := notificationrepo.NewPostgresRepository(db)

	passwordHasher, err := bcrypt.NewPasswordHasher(10)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("bcrypt: %w", err)
	}

	characterPicker := character.NewPicker(characterRepo)

	authService := authservice.NewService(userRepo, sessionRepo, passwordHasher, characterPicker, tm, availableCharRepo)
	authHandler := authhandler.NewHandler(authService)

	botFactory, err := botfactory.New(characterRepo, "ml/models")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("bot factory: %w", err)
	}

	matchmaker := matchmaker.New(characterRepo, botFactory)

	gameService := game.NewService(availableCharRepo)
	matchService := match.NewService(connectionRepo, availableCharRepo, matchmaker, balanceRepo, tm, characterRepo, matchRepo, profileRepo)
	gameHandler := gamehandler.NewHandler(gameService, matchService)

	moneyService := money.NewService(balanceRepo, currencyConvRepo, tm, notificationRepo)
	moneyHandler := moneyhandler.NewHandler(moneyService)

	shopService := shop.NewService(chestRepo, balanceRepo, characterPicker, availableCharRepo, tm)
	shopHandler := shophandler.NewHandler(shopService)

	friendService := friends.NewService(friendRepo, profileRepo, tm, notificationRepo)
	friendHandler := friendhandler.NewHandler(friendService)

	notificationService := notification.NewService(notificationRepo, tm)
	notificationHandler := notificationhandler.NewHandler(notificationService)

	profileService := profile.NewService(profileRepo)
	profileHandler := profilehandler.NewHandler(profileService)

	authMiddleware := authhttp.NewAuthMiddleware(sessionRepo)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/auth/", http.StripPrefix("/api/v1/auth", authhandler.Mux(authHandler)))
	mux.Handle("/api/v1/game/", http.StripPrefix("/api/v1/game", gamehandler.Mux(gameHandler)))
	mux.Handle("/api/v1/money/", http.StripPrefix("/api/v1/money", moneyhandler.Mux(moneyHandler)))
	mux.Handle("/api/v1/shop/", http.StripPrefix("/api/v1/shop", shophandler.Mux(shopHandler)))
	mux.Handle("/api/v1/friends/", http.StripPrefix("/api/v1/friends", friendhandler.Mux(friendHandler)))
	mux.Handle("/api/v1/notifications/", http.StripPrefix("/api/v1/notifications", notificationhandler.Mux(notificationHandler)))
	mux.Handle("/api/v1/profile/", http.StripPrefix("/api/v1/profile", profilehandler.Mux(profileHandler)))

	mux.Handle("/", http.FileServerFS(web.FS))

	return api.PanicHandlerMiddleware(authMiddleware.Middleware(mux)), matchService, matchmaker, nil
}
