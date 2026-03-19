package infra

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"gitlab.com/Uranury/tunescape/pkg/config"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

type Deps struct {
	DBConn      *sqlx.DB
	Logger      *slog.Logger
	Config      *config.Config
	RedisClient *redis.Client
	HTTPClient  *http.Client
}

func New(ctx context.Context) (*Deps, func(), error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	var handler slog.Handler
	if cfg.Env == "development" {
		handler = slog.NewTextHandler(os.Stdout, nil)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	}
	logger := slog.New(handler).With("service", "tunescape")

	if err := database.RunMigrations(cfg.Database.Driver, cfg.Database.DSN(), cfg.MigrationsPath, logger); err != nil {
		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	dbConn, err := database.InitDB(ctx, cfg.Database.Driver, cfg.Database.DSN(), logger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to init database: %w", err)
	}

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Warn("redis ping failed", "error", err)
	}

	deps := &Deps{
		DBConn:      dbConn,
		Logger:      logger,
		Config:      cfg,
		RedisClient: redisClient,
		HTTPClient:  httpClient,
	}

	logger.Info("infrastructure initialized",
		"db_driver", cfg.Database.Driver,
		"redis_addr", cfg.Redis.Addr,
	)

	cleanup := func() {
		if err := dbConn.Close(); err != nil {
			logger.Warn("failed to close database connection", "error", err)
		}
		if err := redisClient.Close(); err != nil {
			logger.Warn("failed to close redis client", "error", err)
		}
		logger.Info("infra cleaned up")
	}

	return deps, cleanup, nil
}
