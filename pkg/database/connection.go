package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// InitDB connects to database with retries and returns the connection
func InitDB(driverName, dsn string, logger *slog.Logger) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Connect(driverName, dsn)
		if err == nil {
			break
		}

		waitTime := time.Duration(i+1) * time.Second
		logger.Warn("Failed to connect to database, retrying...",
			"attempt", i+1,
			"max_retries", maxRetries,
			"wait_time", waitTime,
			"error", err,
		)
		time.Sleep(waitTime)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established successfully")
	return db, nil
}

// RunMigrations runs database migrations using a separate connection
func RunMigrations(driverName, dsn, migrationsPath string, logger *slog.Logger) error {
	migrationURL := migrationsPath
	if !strings.HasPrefix(migrationsPath, "file://") {
		migrationURL = "file://" + migrationsPath
	}

	// Open a separate connection just for migrations
	migrationDB, err := sql.Open(driverName, dsn)
	if err != nil {
		return fmt.Errorf("failed to open migration database: %w", err)
	}
	defer migrationDB.Close() // ← Safe to close, it's a separate connection

	driver, err := postgres.WithInstance(migrationDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationURL,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}
	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			logger.Error("migration source close error", "error", sourceErr)
		}
		if dbErr != nil {
			logger.Error("migration database close error", "error", dbErr)
		}
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("Database is already up to date")
			return nil
		}
		return fmt.Errorf("migration failed: %w", err)
	}

	logger.Info("Database migrations completed successfully")
	return nil
}
