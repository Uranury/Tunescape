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

func New() (*Deps, func(), error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", "tunescape")

	if err := database.RunMigrations(cfg.Driver, cfg.DSN(), cfg.MigrationsPath, logger); err != nil {
		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	dbConn, err := database.InitDB(cfg.Driver, cfg.DSN(), logger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to init database: %w", err)
	}

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
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
		"db_driver", cfg.Driver,
		"redis_addr", cfg.RedisAddr,
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
