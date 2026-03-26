package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "gitlab.com/Uranury/tunescape/docs"
	"gitlab.com/Uranury/tunescape/internal/app"
	"gitlab.com/Uranury/tunescape/internal/auth"
	"gitlab.com/Uranury/tunescape/internal/infra"
	"gitlab.com/Uranury/tunescape/internal/spotify"
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

	userRepo := user.NewRepository(deps.DBConn)
	userSvc := user.NewService(userRepo)

	authRepo := auth.NewRepository(deps.DBConn)
	tokenSvc := auth.NewTokenService([]byte(deps.Config.JWTKey), deps.Logger)
	authSvc := auth.NewRefreshService(tokenSvc, authRepo, txProvider, deps.Logger)
	authHandler := auth.NewHandler(authSvc, tokenSvc, userSvc, deps.Config.IsProd())

	spotifyClient := spotify.NewClient(deps.Config.Spotify, deps.HTTPClient)
	spotifyRepo := spotify.NewRepository(deps.DBConn)
	spotifySvc := spotify.NewService(spotifyRepo, userRepo, spotifyClient)
	spotifyHandler := spotify.NewHandler(
		spotifySvc,
		authSvc,
		deps.Logger,
		deps.Config.IsProd(),
		deps.Config.FrontendURL,
	)

	server := app.NewServer(deps, authHandler, spotifyHandler)

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
