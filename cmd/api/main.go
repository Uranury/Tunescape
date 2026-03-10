package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/Uranury/tunescape/internal/app"
	"gitlab.com/Uranury/tunescape/internal/auth"
	"gitlab.com/Uranury/tunescape/internal/infra"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

func main() {
	deps, cleanup, err := infra.New()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer cleanup()

	txProvider := database.NewTxProvider(deps.DBConn)

	authRepo := auth.NewRepository(deps.DBConn)
	tokenSvc := auth.NewTokenService([]byte(deps.Config.JWTKey), deps.Logger)
	authSvc := auth.NewRefreshService(tokenSvc, authRepo, txProvider, deps.Logger)

	server := app.NewServer(deps)

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
