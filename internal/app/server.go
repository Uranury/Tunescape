package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gitlab.com/Uranury/tunescape/internal/auth"
	"gitlab.com/Uranury/tunescape/internal/infra"
	"gitlab.com/Uranury/tunescape/internal/spotify"
)

type Server struct {
	router         *gin.Engine
	logger         *slog.Logger
	httpServer     *http.Server
	authHandler    *auth.Handler
	spotifyHandler *spotify.Handler
}

func NewServer(deps *infra.Deps, authHandler *auth.Handler, spotifyHandler *spotify.Handler) *Server {
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
		router:         router,
		logger:         deps.Logger,
		authHandler:    authHandler,
		spotifyHandler: spotifyHandler,
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
	s.router.StaticFile("/", "./frontend/index.html")
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authGroup := s.router.Group("/auth")
	{
		authGroup.POST("/login", s.authHandler.Login)
		authGroup.POST("/signup", s.authHandler.Signup)
		authGroup.POST("/logout", s.authHandler.Logout)
		authGroup.POST("/refresh", s.authHandler.Refresh)
		authGroup.GET("/spotify/login", s.spotifyHandler.LoginHandler)
		authGroup.GET("/spotify/callback", s.spotifyHandler.CallbackHandler)
	}
}
