package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/Uranury/tunescape/internal/app"
	"gitlab.com/Uranury/tunescape/internal/infra"
)

func main() {
	deps, cleanup, err := infra.New()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer cleanup()

	server := app.NewServer(deps)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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
