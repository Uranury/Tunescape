package friends

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

type Repository interface {
	SendRequest(ctx context.Context, senderID, receiverID uuid.UUID) error
	GetRequest(ctx context.Context, requestID int64) (*FriendRequest, error)
	AcceptRequest(ctx context.Context, requestID int64, senderID, receiverID uuid.UUID) error
	RejectRequest(ctx context.Context, requestID int64, receiverID uuid.UUID) error
	ListIncoming(ctx context.Context, userID uuid.UUID) ([]IncomingRequest, error)
	ListFriends(ctx context.Context, userID uuid.UUID) ([]FriendProfile, error)
	AreFriends(ctx context.Context, userID, friendID uuid.UUID) (bool, error)
	RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error
	HasSpotifyConnected(ctx context.Context, userID uuid.UUID) (bool, error)
}

type repository struct {
	db database.Executor
}

func NewRepository(db database.Executor) Repository {
	return &repository{db: db}
}

func (r *repository) SendRequest(ctx context.Context, senderID, receiverID uuid.UUID) error {
	const q = `
		INSERT INTO friend_requests (sender_id, receiver_id)
		VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, q, senderID, receiverID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return apperrors.ErrRequestAlreadySent
		}
		return err
	}
	return nil
}

func (r *repository) GetRequest(ctx context.Context, requestID int64) (*FriendRequest, error) {
	const q = `SELECT id, sender_id, receiver_id, status, created_at, updated_at
	           FROM friend_requests WHERE id = $1`
	var req FriendRequest
	if err := r.db.GetContext(ctx, &req, q, requestID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrRequestNotFound
		}
		return nil, err
	}
	return &req, nil
}

func (r *repository) AcceptRequest(ctx context.Context, requestID int64, senderID, receiverID uuid.UUID) error {
	const updateQ = `
		UPDATE friend_requests SET status = 'accepted', updated_at = NOW()
		WHERE id = $1 AND receiver_id = $2 AND status = 'pending'`
	res, err := r.db.ExecContext(ctx, updateQ, requestID, receiverID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return apperrors.ErrRequestNotFound
	}

	const insertQ = `
		INSERT INTO friends (user_id, friend_id) VALUES ($1, $2), ($2, $1)
		ON CONFLICT DO NOTHING`
	_, err = r.db.ExecContext(ctx, insertQ, senderID, receiverID)
	return err
}

func (r *repository) RejectRequest(ctx context.Context, requestID int64, receiverID uuid.UUID) error {
	const q = `
		UPDATE friend_requests SET status = 'rejected', updated_at = NOW()
		WHERE id = $1 AND receiver_id = $2 AND status = 'pending'`
	res, err := r.db.ExecContext(ctx, q, requestID, receiverID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return apperrors.ErrRequestNotFound
	}
	return nil
}

func (r *repository) ListIncoming(ctx context.Context, userID uuid.UUID) ([]IncomingRequest, error) {
	const q = `
		SELECT fr.id AS request_id, fr.sender_id, u.display_name, fr.created_at
		FROM friend_requests fr
		JOIN users u ON u.id = fr.sender_id
		WHERE fr.receiver_id = $1 AND fr.status = 'pending'
		ORDER BY fr.created_at DESC`
	var out []IncomingRequest
	if err := r.db.SelectContext(ctx, &out, q, userID); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repository) ListFriends(ctx context.Context, userID uuid.UUID) ([]FriendProfile, error) {
	const q = `
		SELECT u.id AS user_id,
		       u.display_name,
		       u.spotify_id IS NOT NULL AS spotify_connected,
		       u.spotify_id
		FROM friends f
		JOIN users u ON u.id = f.friend_id
		WHERE f.user_id = $1
		ORDER BY u.display_name`
	var out []FriendProfile
	if err := r.db.SelectContext(ctx, &out, q, userID); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repository) AreFriends(ctx context.Context, userID, friendID uuid.UUID) (bool, error) {
	const q = `SELECT EXISTS(SELECT 1 FROM friends WHERE user_id = $1 AND friend_id = $2)`
	var ok bool
	if err := r.db.GetContext(ctx, &ok, q, userID, friendID); err != nil {
		return false, err
	}
	return ok, nil
}

func (r *repository) RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	const q = `DELETE FROM friends WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)`
	_, err := r.db.ExecContext(ctx, q, userID, friendID)
	return err
}

func (r *repository) HasSpotifyConnected(ctx context.Context, userID uuid.UUID) (bool, error) {
	const q = `SELECT spotify_id IS NOT NULL FROM users WHERE id = $1`
	var ok bool
	if err := r.db.GetContext(ctx, &ok, q, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, apperrors.ErrNotFound
		}
		return false, err
	}
	return ok, nil
}
