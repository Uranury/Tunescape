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
	"gitlab.com/Uranury/tunescape/internal/analytics"
	"gitlab.com/Uranury/tunescape/internal/auth"
	"gitlab.com/Uranury/tunescape/internal/friends"
	"gitlab.com/Uranury/tunescape/internal/infra"
	"gitlab.com/Uranury/tunescape/internal/leaderboard"
	"gitlab.com/Uranury/tunescape/internal/middleware"
	"gitlab.com/Uranury/tunescape/internal/playlist"
	"gitlab.com/Uranury/tunescape/internal/report"
	"gitlab.com/Uranury/tunescape/internal/snapshot"
	"gitlab.com/Uranury/tunescape/internal/spotify"
	"gitlab.com/Uranury/tunescape/internal/trends"
	"gitlab.com/Uranury/tunescape/internal/user"
	"gitlab.com/Uranury/tunescape/internal/worker"
)

type Server struct {
	router             *gin.Engine
	logger             *slog.Logger
	httpServer         *http.Server
	authHandler        *auth.Handler
	spotifyHandler     *spotify.Handler
	snapshotHandler    *snapshot.Handler
	analyticsHandler   *analytics.Handler
	leaderboardHandler *leaderboard.Handler
	trendsHandler      *trends.Handler
	reportHandler      *report.Handler
	userHandler        *user.Handler
	playlistHandler    *playlist.Handler
	friendHandler      *friends.Handler
	authMiddleware     *middleware.Auth
	rateLimiter        *middleware.RateLimiter
	snapshotWorker     *worker.SnapshotWorker
}

func NewServer(
	deps *infra.Deps,
	authHandler *auth.Handler,
	spotifyHandler *spotify.Handler,
	snapshotHandler *snapshot.Handler,
	analyticsHandler *analytics.Handler,
	leaderboardHandler *leaderboard.Handler,
	trendsHandler *trends.Handler,
	reportHandler *report.Handler,
	userHandler *user.Handler,
	playlistHandler *playlist.Handler,
	friendHandler *friends.Handler,
	authMiddleware *middleware.Auth,
	rateLimiter *middleware.RateLimiter,
	snapshotWorker *worker.SnapshotWorker,
) *Server {
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
		router:             router,
		logger:             deps.Logger,
		authHandler:        authHandler,
		spotifyHandler:     spotifyHandler,
		snapshotHandler:    snapshotHandler,
		analyticsHandler:   analyticsHandler,
		leaderboardHandler: leaderboardHandler,
		trendsHandler:      trendsHandler,
		reportHandler:      reportHandler,
		userHandler:        userHandler,
		playlistHandler:    playlistHandler,
		friendHandler:      friendHandler,
		authMiddleware:     authMiddleware,
		rateLimiter:        rateLimiter,
		snapshotWorker:     snapshotWorker,
		httpServer: &http.Server{
			Addr:         deps.Config.ListenAddr,
			Handler:      router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	server.registerRoutes()
	return server
}

func (s *Server) Run() error {
	s.logger.Info("server starting", "addr", s.httpServer.Addr)
	s.snapshotWorker.Start()
	s.logger.Info("snapshot worker started")
	if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.snapshotWorker.Stop()
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
	s.router.Static("/css", "./frontend/css")
	s.router.Static("/js", "./frontend/js")
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

	meGroup := s.router.Group("/me", s.authMiddleware.JWTAuth())
	{
		meGroup.GET("/profile", s.userHandler.GetProfile)
		meGroup.DELETE("/spotify", s.spotifyHandler.DisconnectHandler)
		meGroup.POST("/snapshots", s.snapshotHandler.CreateSnapshot)
		meGroup.GET("/snapshots", s.snapshotHandler.ListSnapshots)
		meGroup.GET("/snapshots/:id", s.snapshotHandler.GetSnapshot)
		meGroup.GET("/trends", s.trendsHandler.GetTrends)
		meGroup.GET("/report", s.reportHandler.GetReport)
		meGroup.POST("/playlists/top-tracks", s.playlistHandler.CreateFromSnapshot)
	}

	analyticsGroup := s.router.Group("/analytics", s.authMiddleware.JWTAuth())
	{
		analyticsGroup.GET("/top-tracks", s.analyticsHandler.GetMusicTaste)
	}

	s.router.GET("/leaderboards/:feature", s.rateLimiter.PerIP(), s.leaderboardHandler.GetLeaderboard)

	usersGroup := s.router.Group("/users", s.authMiddleware.JWTAuth())
	{
		usersGroup.GET("/lookup", s.userHandler.LookupUser)
	}

	friendsGroup := s.router.Group("/friends", s.authMiddleware.JWTAuth())
	{
		friendsGroup.POST("/requests", s.friendHandler.SendRequest)
		friendsGroup.GET("/requests", s.friendHandler.ListIncoming)
		friendsGroup.POST("/requests/:id/accept", s.friendHandler.AcceptRequest)
		friendsGroup.POST("/requests/:id/reject", s.friendHandler.RejectRequest)
		friendsGroup.GET("", s.friendHandler.ListFriends)
		friendsGroup.GET("/:friend_id/compare", s.friendHandler.CompareTastes)
		friendsGroup.GET("/:friend_id/playlists", s.friendHandler.GetFriendPlaylists)
	}
}
