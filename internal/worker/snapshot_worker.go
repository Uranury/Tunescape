package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/snapshot"
	"gitlab.com/Uranury/tunescape/internal/spotify"
	"gitlab.com/Uranury/tunescape/internal/user"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

// SnapshotWorker creates snapshots for all users on a schedule
type SnapshotWorker struct {
	userRepo       user.Repository
	snapshotSvc    snapshot.Service
	txProvider     database.TxProvider
	logger         *slog.Logger
	interval       time.Duration
	ctx            context.Context
	cancel         context.CancelFunc
	done           chan struct{}
}

// NewSnapshotWorker creates a new snapshot worker
func NewSnapshotWorker(
	userRepo user.Repository,
	snapshotSvc snapshot.Service,
	txProvider database.TxProvider,
	logger *slog.Logger,
	interval time.Duration,
) *SnapshotWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &SnapshotWorker{
		userRepo:    userRepo,
		snapshotSvc: snapshotSvc,
		txProvider:  txProvider,
		logger:      logger,
		interval:    interval,
		ctx:         ctx,
		cancel:      cancel,
		done:        make(chan struct{}),
	}
}

// Start begins the worker loop
func (w *SnapshotWorker) Start() {
	go w.run()
}

// Stop gracefully stops the worker
func (w *SnapshotWorker) Stop() {
	w.logger.Info("stopping snapshot worker")
	w.cancel()
	<-w.done
	w.logger.Info("snapshot worker stopped")
}

func (w *SnapshotWorker) run() {
	defer close(w.done)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Run immediately on start
	w.createSnapshotsForAllUsers()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			w.createSnapshotsForAllUsers()
		}
	}
}

func (w *SnapshotWorker) createSnapshotsForAllUsers() {
	w.logger.Info("starting snapshot creation for all users")

	users, err := w.userRepo.FindAll(w.ctx)
	if err != nil {
		w.logger.Error("failed to fetch all users", "error", err)
		return
	}

	w.logger.Info("fetched users", "count", len(users))

	successCount := 0
	failureCount := 0

	for _, u := range users {
		userID := u.ID

		// Check if user has Spotify connected
		if u.SpotifyID == nil {
			w.logger.Debug("skipping user without spotify connection", "user_id", userID)
			continue
		}

		if err := w.createSnapshot(userID); err != nil {
			w.logger.Error("failed to create snapshot", "user_id", userID, "error", err)
			failureCount++
			continue
		}

		successCount++
	}

	w.logger.Info("snapshot creation batch completed",
		"success", successCount,
		"failed", failureCount,
		"skipped", len(users)-successCount-failureCount,
	)
}

func (w *SnapshotWorker) createSnapshot(userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(w.ctx, 30*time.Second)
	defer cancel()

	_, err := w.snapshotSvc.CreateSnapshot(ctx, userID, spotify.MediumTerm)
	if err != nil {
		return fmt.Errorf("create snapshot: %w", err)
	}

	return nil
}
