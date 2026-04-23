package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "gitlab.com/Uranury/tunescape/docs"
	"gitlab.com/Uranury/tunescape/internal/analytics"
	"gitlab.com/Uranury/tunescape/internal/app"
	"gitlab.com/Uranury/tunescape/internal/auth"
	"gitlab.com/Uranury/tunescape/internal/cache"
	"gitlab.com/Uranury/tunescape/internal/infra"
	"gitlab.com/Uranury/tunescape/internal/leaderboard"
	"gitlab.com/Uranury/tunescape/internal/middleware"
	"gitlab.com/Uranury/tunescape/internal/reccobeats"
	"gitlab.com/Uranury/tunescape/internal/report"
	"gitlab.com/Uranury/tunescape/internal/snapshot"
	"gitlab.com/Uranury/tunescape/internal/spotify"
	"gitlab.com/Uranury/tunescape/internal/trends"
	"gitlab.com/Uranury/tunescape/internal/user"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

// @title Tunescape API
// @version 1.0
// @description Tunescape backend API.
// @BasePath /
func main() {
	deps, cleanup, err := infra.New(context.Background())
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer cleanup()

	txProvider := database.NewTxProvider(deps.DBConn)
	redisCache := cache.New(deps.RedisClient)

	userRepo := user.NewRepository(deps.DBConn)
	userSvc := user.NewService(userRepo)
	userHandler := user.NewHandler(userSvc)

	authRepo := auth.NewRepository(deps.DBConn)
	tokenSvc := auth.NewTokenService([]byte(deps.Config.JWTKey), deps.Logger)
	authSvc := auth.NewRefreshService(tokenSvc, authRepo, txProvider, deps.Logger)
	authHandler := auth.NewHandler(authSvc, tokenSvc, userSvc, deps.Config.IsProd())

	spotifyClient := spotify.NewClient(deps.Config.Spotify, deps.HTTPClient)
	spotifyRepo := spotify.NewRepository(deps.DBConn)
	spotifySvc := spotify.NewService(spotifyRepo, userRepo, spotifyClient, txProvider, deps.Logger)
	spotifyHandler := spotify.NewHandler(spotifySvc, authSvc, deps.Logger, deps.Config.IsProd(), deps.Config.FrontendURL)

	snapshotRepo := snapshot.NewRepository(deps.DBConn)
	snapshotSvc := snapshot.NewService(snapshotRepo, spotifySvc, txProvider, redisCache, deps.Logger)
	snapshotHandler := snapshot.NewHandler(snapshotSvc)

	leaderboardStore := leaderboard.NewStore(deps.RedisClient)
	leaderboardSvc := leaderboard.NewService(leaderboardStore, userRepo, redisCache)
	leaderboardHandler := leaderboard.NewHandler(leaderboardSvc)

	reccobeatsClient := reccobeats.NewClient(deps.Config.Reccobeats, deps.HTTPClient)
	reccobeatsService := reccobeats.NewService(reccobeatsClient)
	analyticsRepo := analytics.NewRepository(deps.DBConn)
	analyticsSvc := analytics.NewService(analyticsRepo, reccobeatsService, txProvider, deps.Logger, redisCache, leaderboardSvc)
	analyticsHandler := analytics.NewHandler(analyticsSvc)

	trendsRepo := trends.NewRepository(deps.DBConn)
	trendsSvc := trends.NewService(trendsRepo)
	trendsHandler := trends.NewHandler(trendsSvc)

	reportRepo := report.NewRepository(deps.DBConn)
	reportSvc := report.NewService(reportRepo, leaderboardSvc, userRepo)
	reportHandler := report.NewHandler(reportSvc)

	authMiddleware := middleware.NewAuth(tokenSvc)
	rateLimiter := middleware.NewRateLimiter(deps.RedisClient, 60, time.Minute)

	server := app.NewServer(
		deps,
		authHandler,
		spotifyHandler,
		snapshotHandler,
		analyticsHandler,
		leaderboardHandler,
		trendsHandler,
		reportHandler,
		userHandler,
		authMiddleware,
		rateLimiter,
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	authSvc.StartCleanup(ctx)

	go func() {
		if err := server.Run(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Shutdown failed: %v", err)
	}

	log.Println("Server exited cleanly")
}
