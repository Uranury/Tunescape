package auth

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"io"
	"log/slog"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gitlab.com/Uranury/tunescape/pkg/database"
)

type mockRefreshRepo struct {
	saveFn          func(ctx context.Context, token *RefreshToken) error
	findByHashFn    func(ctx context.Context, tokenHash string) (*RefreshToken, error)
	findByHashForFn func(ctx context.Context, tokenHash string) (*RefreshToken, error)
	revokeByHashFn  func(ctx context.Context, tokenHash string) error
	deleteExpiredFn func(ctx context.Context) error
}

func (m *mockRefreshRepo) Save(ctx context.Context, token *RefreshToken) error {
	return m.saveFn(ctx, token)
}

func (m *mockRefreshRepo) FindByHash(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	return m.findByHashFn(ctx, tokenHash)
}

func (m *mockRefreshRepo) FindByHashForUpdate(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	return m.findByHashForFn(ctx, tokenHash)
}

func (m *mockRefreshRepo) RevokeByHash(ctx context.Context, tokenHash string) error {
	return m.revokeByHashFn(ctx, tokenHash)
}

func (m *mockRefreshRepo) DeleteExpired(ctx context.Context) error {
	return m.deleteExpiredFn(ctx)
}

type mockTokenService struct {
	generateFn func(userID uuid.UUID, role string) (string, error)
	validateFn func(tokenString string) (*Claims, error)
}

func (m *mockTokenService) Generate(userID uuid.UUID, role string) (string, error) {
	return m.generateFn(userID, role)
}

func (m *mockTokenService) Validate(tokenString string) (*Claims, error) {
	return m.validateFn(tokenString)
}

type tokenHashArgMatcher struct{ oldHash string }

func (m tokenHashArgMatcher) Match(v driver.Value) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	// token hash is base64(url) of sha256(token) => 32 bytes after decode.
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return false
	}
	if len(b) != 32 {
		return false
	}
	// With extremely low probability it may match the old hash; still allow.
	// (We don't assert inequality here to avoid flakiness.)
	_ = m.oldHash
	return true
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestRefreshTokenService_Generate_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	role := "user"
	userAgent := "ua"
	ip := "127.0.0.1"

	var captured *RefreshToken
	repo := &mockRefreshRepo{
		saveFn: func(ctx context.Context, token *RefreshToken) error {
			captured = token
			return nil
		},
	}

	svc := NewRefreshService(nil, repo, nil, testLogger())

	before := time.Now()
	token, err := svc.Generate(ctx, userID, role, userAgent, ip)
	after := time.Now()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if token == "" {
		t.Fatalf("expected non-empty token")
	}

	raw, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		t.Fatalf("token should be base64: %v", err)
	}
	if len(raw) != 32 {
		t.Fatalf("expected 32 token bytes, got %d", len(raw))
	}

	if captured == nil {
		t.Fatalf("expected repository Save to be called")
	}
	if captured.UserID != userID {
		t.Fatalf("expected userID %s, got %s", userID, captured.UserID)
	}
	if captured.Role != role {
		t.Fatalf("expected role %s, got %s", role, captured.Role)
	}
	if captured.TokenHash != hashToken(token) {
		t.Fatalf("expected token hash %q, got %q", hashToken(token), captured.TokenHash)
	}
	if captured.UserAgent != userAgent {
		t.Fatalf("expected userAgent %q, got %q", userAgent, captured.UserAgent)
	}
	if captured.IP != ip {
		t.Fatalf("expected ip %q, got %q", ip, captured.IP)
	}

	lower := before.Add(RefreshTokenTTL - 2*time.Second)
	upper := after.Add(RefreshTokenTTL + 2*time.Second)
	if captured.ExpiresAt.Before(lower) || captured.ExpiresAt.After(upper) {
		t.Fatalf("unexpected ExpiresAt: %v (expected in [%v, %v])", captured.ExpiresAt, lower, upper)
	}
}

func TestRefreshTokenService_Generate_SaveError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expectedErr := errors.New("db down")
	repo := &mockRefreshRepo{
		saveFn: func(ctx context.Context, token *RefreshToken) error { return expectedErr },
	}

	svc := NewRefreshService(nil, repo, nil, testLogger())
	_, err := svc.Generate(ctx, uuid.New(), "user", "ua", "127.0.0.1")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "failed to save refresh token") {
		t.Fatalf("unexpected error: %v", err)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}
}

func TestRefreshTokenService_Validate_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	tokenString := "some-token"
	expectedToken := &RefreshToken{
		UserID:    uuid.New(),
		Role:      "user",
		TokenHash: hashToken(tokenString),
	}

	repo := &mockRefreshRepo{
		findByHashFn: func(ctx context.Context, tokenHash string) (*RefreshToken, error) {
			if tokenHash != expectedToken.TokenHash {
				t.Fatalf("expected tokenHash %q, got %q", expectedToken.TokenHash, tokenHash)
			}
			return expectedToken, nil
		},
	}

	svc := NewRefreshService(nil, repo, nil, testLogger())
	got, err := svc.Validate(ctx, tokenString)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got != expectedToken {
		t.Fatalf("expected same token instance")
	}
}

func TestRefreshTokenService_Validate_NotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := &mockRefreshRepo{
		findByHashFn: func(ctx context.Context, tokenHash string) (*RefreshToken, error) {
			return nil, sql.ErrNoRows
		},
	}
	svc := NewRefreshService(nil, repo, nil, testLogger())
	_, err := svc.Validate(ctx, "bad-token")
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "refresh token not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRefreshTokenService_Validate_OtherError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expectedErr := errors.New("boom")
	repo := &mockRefreshRepo{
		findByHashFn: func(ctx context.Context, tokenHash string) (*RefreshToken, error) {
			return nil, expectedErr
		},
	}
	svc := NewRefreshService(nil, repo, nil, testLogger())
	_, err := svc.Validate(ctx, "token")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "failed to validate refresh token") {
		t.Fatalf("unexpected error: %v", err)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}
}

func TestRefreshTokenService_Revoke_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	refreshToken := "rt"
	repo := &mockRefreshRepo{
		revokeByHashFn: func(ctx context.Context, tokenHash string) error {
			if tokenHash != hashToken(refreshToken) {
				t.Fatalf("unexpected tokenHash: %q", tokenHash)
			}
			return nil
		},
	}

	svc := NewRefreshService(nil, repo, nil, testLogger())
	if err := svc.Revoke(ctx, refreshToken); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestRefreshTokenService_Revoke_Error(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	refreshToken := "rt"
	expectedErr := errors.New("cannot revoke")
	repo := &mockRefreshRepo{
		revokeByHashFn: func(ctx context.Context, tokenHash string) error {
			return expectedErr
		},
	}

	svc := NewRefreshService(nil, repo, nil, testLogger())
	err := svc.Revoke(ctx, refreshToken)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "failed to revoke refresh token") {
		t.Fatalf("unexpected error: %v", err)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}
}

func TestRefreshTokenService_StartCleanup_CtxCanceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var called atomic.Bool
	repo := &mockRefreshRepo{
		deleteExpiredFn: func(ctx context.Context) error {
			called.Store(true)
			return nil
		},
	}

	svc := NewRefreshService(nil, repo, nil, testLogger())
	svc.StartCleanup(ctx)

	time.Sleep(30 * time.Millisecond)
	if called.Load() {
		t.Fatalf("DeleteExpired should not be called when context is already canceled")
	}
}

func TestRefreshTokenService_Refresh_Success(t *testing.T) {
	t.Parallel()

	// Create sqlmock-backed DB so that auth.Refresh can execute
	// internal/auth/repository.go queries and Scan results.
	rawDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	db := sqlx.NewDb(rawDB, "sqlmock")
	txProvider := database.NewTxProvider(db)

	logger := testLogger()

	userID := uuid.New()
	role := "user"
	oldRefreshToken := "old-refresh-token"
	oldTokenHash := hashToken(oldRefreshToken)
	userAgent := "ua"
	ip := "127.0.0.1"
	expiresAt := time.Now().Add(1 * time.Hour)
	createdAt := time.Now().Add(-1 * time.Hour)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT \\* FROM refresh_tokens").
		WithArgs(oldTokenHash).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "user_id", "role", "token_hash", "expires_at", "created_at", "revoked_at", "user_agent", "ip",
			}).AddRow(
				10, userID, role, oldTokenHash, expiresAt, createdAt, nil, userAgent, ip,
			),
		)

	mock.ExpectExec("UPDATE refresh_tokens SET revoked_at = NOW\\(\\) WHERE token_hash").
		WithArgs(oldTokenHash).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Insert new refresh token with unknown generated token hash/time.
	mock.ExpectQuery("INSERT INTO refresh_tokens").
		WithArgs(
			userID,
			role,
			tokenHashArgMatcher{oldHash: oldTokenHash},
			sqlmock.AnyArg(),
			userAgent,
			ip,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(99))

	mock.ExpectCommit()
	mock.ExpectClose()

	tokenSvc := &mockTokenService{
		generateFn: func(gotUserID uuid.UUID, gotRole string) (string, error) {
			if gotUserID != userID {
				t.Fatalf("expected userID %s, got %s", userID, gotUserID)
			}
			if gotRole != role {
				t.Fatalf("expected role %q, got %q", role, gotRole)
			}
			return "access-token", nil
		},
		validateFn: func(tokenString string) (*Claims, error) {
			return nil, nil
		},
	}

	svc := NewRefreshService(tokenSvc, nil, txProvider, logger)
	accessToken, newRefreshToken, err := svc.Refresh(context.Background(), oldRefreshToken)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if accessToken != "access-token" {
		t.Fatalf("unexpected access token: %q", accessToken)
	}

	if newRefreshToken == "" {
		t.Fatalf("expected non-empty new refresh token")
	}
	raw, err := base64.URLEncoding.DecodeString(newRefreshToken)
	if err != nil {
		t.Fatalf("new refresh token should be base64: %v", err)
	}
	if len(raw) != 32 {
		t.Fatalf("expected 32 token bytes, got %d", len(raw))
	}
	if newHash := hashToken(newRefreshToken); newHash == oldTokenHash {
		t.Fatalf("new token hash unexpectedly equals old token hash")
	}

	if err := rawDB.Close(); err != nil {
		t.Fatalf("failed to close db: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}

}

func TestRefreshTokenService_Refresh_RevokedToken(t *testing.T) {
	t.Parallel()

	rawDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	db := sqlx.NewDb(rawDB, "sqlmock")
	txProvider := database.NewTxProvider(db)

	revokedToken := "already-revoked-token"
	revokedTokenHash := hashToken(revokedToken)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT \\* FROM refresh_tokens").
		WithArgs(revokedTokenHash).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()
	mock.ExpectClose()

	tokenSvc := &mockTokenService{
		generateFn: func(userID uuid.UUID, role string) (string, error) {
			t.Fatal("Generate must not be called for a revoked token")
			return "", nil
		},
	}

	svc := NewRefreshService(tokenSvc, nil, txProvider, testLogger())
	accessToken, newRefreshToken, err := svc.Refresh(context.Background(), revokedToken)

	if err == nil {
		t.Fatal("expected an error for a revoked token, got nil")
	}
	if !strings.Contains(err.Error(), "refresh token not found") {
		t.Fatalf("expected 'refresh token not found' in error, got: %v", err)
	}
	if accessToken != "" {
		t.Fatalf("expected empty access token, got %q", accessToken)
	}
	if newRefreshToken != "" {
		t.Fatalf("expected empty refresh token, got %q", newRefreshToken)
	}

	if err := rawDB.Close(); err != nil {
		t.Fatalf("failed to close db: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
