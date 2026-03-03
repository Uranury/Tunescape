package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.com/Uranury/tunescape/internal/infra"
	"log/slog"
)

type Server struct {
	router     *gin.Engine
	logger     *slog.Logger
	httpServer *http.Server
	// TODO: expand with services
}

func NewServer(deps *infra.Deps) *Server {
	router := gin.New()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		cors.New(cors.Config{
			AllowOrigins:     deps.Config.AllowedOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}),
	)

	server := &Server{
		router: router,
		logger: deps.Logger,
		httpServer: &http.Server{
			Addr:         deps.Config.ListenAddr,
			Handler:      router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	server.registerRoutes()
	return server
}

func (s *Server) Run() error {
	s.logger.Info("server starting", "addr", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	s.logger.Info("Shutting down HTTP server")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("Failed to gracefully shutdown", "error", err)
		return err
	}

	s.logger.Info("HTTP server shut down gracefully")
	return nil
}

func (s *Server) registerRoutes() {
	// routes go here
}
