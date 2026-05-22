package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	findByEmailFn              func(ctx context.Context, email string) (*User, error)
	findByIDFn                 func(ctx context.Context, userID uuid.UUID) (*User, error)
	createFn                   func(ctx context.Context, u *User) error
	findDisplayNameFn          func(ctx context.Context, userID uuid.UUID) (string, error)
	findDisplayNamesByIDsFn    func(ctx context.Context, userIDs []string) (map[string]string, error)
	connectSpotifyFn           func(ctx context.Context, userID uuid.UUID, spotifyID *string, avatarURL, country, product *string) error
	clearSpotifyFn             func(ctx context.Context, userID uuid.UUID) error
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*User, error) {
	if m.findByEmailFn != nil {
		return m.findByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m *mockUserRepo) FindByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, userID)
	}
	return nil, nil
}

func (m *mockUserRepo) Create(ctx context.Context, u *User) error {
	if m.createFn != nil {
		return m.createFn(ctx, u)
	}
	return nil
}

func (m *mockUserRepo) FindDisplayName(ctx context.Context, userID uuid.UUID) (string, error) {
	if m.findDisplayNameFn != nil {
		return m.findDisplayNameFn(ctx, userID)
	}
	return "", nil
}

func (m *mockUserRepo) FindDisplayNamesByIDs(ctx context.Context, userIDs []string) (map[string]string, error) {
	if m.findDisplayNamesByIDsFn != nil {
		return m.findDisplayNamesByIDsFn(ctx, userIDs)
	}
	return nil, nil
}

func (m *mockUserRepo) ConnectSpotify(ctx context.Context, userID uuid.UUID, spotifyID *string, avatarURL, country, product *string) error {
	if m.connectSpotifyFn != nil {
		return m.connectSpotifyFn(ctx, userID, spotifyID, avatarURL, country, product)
	}
	return nil
}

func (m *mockUserRepo) ClearSpotify(ctx context.Context, userID uuid.UUID) error {
	if m.clearSpotifyFn != nil {
		return m.clearSpotifyFn(ctx, userID)
	}
	return nil
}

func (m *mockUserRepo) FindByDisplayName(ctx context.Context, displayName string) (*User, error) {
	return nil, nil
}

func (m *mockUserRepo) FindAvatarURLsByIDs(_ context.Context, _ []string) (map[string]*string, error) {
	return map[string]*string{}, nil
}

func (m *mockUserRepo) FindAll(ctx context.Context) ([]User, error) {
	return nil, nil
}

// TestUserService_ValidateCredentials_Success tests successful credential validation
func TestUserService_ValidateCredentials_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	email := "user@example.com"
	password := "correct_password"
	userID := uuid.New()

	// Create a hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	storedUser := &User{
		ID:       userID,
		Email:    email,
		Password: string(hashedPassword),
	}

	repo := &mockUserRepo{
		findByEmailFn: func(ctx context.Context, e string) (*User, error) {
			if e != email {
				t.Fatalf("expected email %q, got %q", email, e)
			}
			return storedUser, nil
		},
	}

	svc := NewService(repo)
	user, err := svc.ValidateCredentials(ctx, email, password)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected non-nil user")
	}
	if user.ID != userID {
		t.Fatalf("expected user ID %s, got %s", userID, user.ID)
	}
	if user.Email != email {
		t.Fatalf("expected email %q, got %q", email, user.Email)
	}
}

// TestUserService_ValidateCredentials_WrongPassword tests rejection of wrong password
func TestUserService_ValidateCredentials_WrongPassword(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	email := "user@example.com"
	correctPassword := "correct_password"
	wrongPassword := "wrong_password"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
	storedUser := &User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hashedPassword),
	}

	repo := &mockUserRepo{
		findByEmailFn: func(ctx context.Context, e string) (*User, error) {
			return storedUser, nil
		},
	}

	svc := NewService(repo)
	_, err := svc.ValidateCredentials(ctx, email, wrongPassword)

	if err == nil {
		t.Fatalf("expected error for wrong password")
	}
	if !errors.Is(err, apperrors.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

// TestUserService_ValidateCredentials_UserNotFound tests rejection when user not found
func TestUserService_ValidateCredentials_UserNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	email := "nonexistent@example.com"

	repo := &mockUserRepo{
		findByEmailFn: func(ctx context.Context, e string) (*User, error) {
			return nil, sql.ErrNoRows
		},
	}

	svc := NewService(repo)
	_, err := svc.ValidateCredentials(ctx, email, "any_password")

	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, apperrors.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

// TestUserService_ValidateCredentials_RepositoryError tests error propagation
func TestUserService_ValidateCredentials_RepositoryError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbErr := errors.New("database error")

	repo := &mockUserRepo{
		findByEmailFn: func(ctx context.Context, e string) (*User, error) {
			return nil, dbErr
		},
	}

	svc := NewService(repo)
	_, err := svc.ValidateCredentials(ctx, "user@example.com", "password")

	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected database error to be wrapped, got %v", err)
	}
}

// TestUserService_Create_Success tests successful user creation
func TestUserService_Create_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	email := "newuser@example.com"
	password := "secure_password"
	displayName := "New User"

	repo := &mockUserRepo{
		createFn: func(ctx context.Context, u *User) error {
			u.ID = uuid.New()
			return nil
		},
	}

	svc := NewService(repo)
	user, err := svc.Create(ctx, email, password, displayName)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected non-nil user")
	}

	if user.Email != email {
		t.Fatalf("expected email %q, got %q", email, user.Email)
	}
	if user.DisplayName != displayName {
		t.Fatalf("expected display name %q, got %q", displayName, user.DisplayName)
	}
	if user.Role != "user" {
		t.Fatalf("expected role %q, got %q", "user", user.Role)
	}

	// Verify password is hashed
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		t.Fatalf("password not properly hashed: %v", err)
	}
}

// TestUserService_Create_RepositoryError tests error handling during creation
func TestUserService_Create_RepositoryError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	createErr := errors.New("email already exists")

	repo := &mockUserRepo{
		createFn: func(ctx context.Context, u *User) error {
			return createErr
		},
	}

	svc := NewService(repo)
	_, err := svc.Create(ctx, "user@example.com", "password", "displayname")

	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, createErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
}

// TestUserService_GetProfile_Success tests successful profile retrieval
func TestUserService_GetProfile_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	email := "user@example.com"
	displayName := "Test User"
	avatarURL := "https://example.com/avatar.jpg"
	spotifyID := "spotify_user_123"

	storedUser := &User{
		ID:          userID,
		Email:       email,
		DisplayName: displayName,
		AvatarURL:   &avatarURL,
		SpotifyID:   &spotifyID,
	}

	repo := &mockUserRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*User, error) {
			if id != userID {
				t.Fatalf("expected user ID %s, got %s", userID, id)
			}
			return storedUser, nil
		},
	}

	svc := NewService(repo)
	profile, err := svc.GetProfile(ctx, userID)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if profile == nil {
		t.Fatal("expected non-nil profile")
	}

	if profile.Email != email {
		t.Fatalf("expected email %q, got %q", email, profile.Email)
	}
	if profile.DisplayName != displayName {
		t.Fatalf("expected display name %q, got %q", displayName, profile.DisplayName)
	}
	if profile.AvatarURL == nil || *profile.AvatarURL != avatarURL {
		t.Fatalf("expected avatar URL %q, got %v", avatarURL, profile.AvatarURL)
	}
	if !profile.SpotifyConnected {
		t.Fatalf("expected spotify connected to be true")
	}
}

// TestUserService_GetProfile_NoSpotify tests profile for user without Spotify connection
func TestUserService_GetProfile_NoSpotify(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()

	storedUser := &User{
		ID:          userID,
		Email:       "user@example.com",
		DisplayName: "Test User",
		SpotifyID:   nil,
	}

	repo := &mockUserRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*User, error) {
			return storedUser, nil
		},
	}

	svc := NewService(repo)
	profile, err := svc.GetProfile(ctx, userID)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if profile.SpotifyConnected {
		t.Fatalf("expected spotify connected to be false")
	}
}

// TestUserService_GetProfile_UserNotFound tests error handling for missing user
func TestUserService_GetProfile_UserNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	notFoundErr := sql.ErrNoRows

	repo := &mockUserRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*User, error) {
			return nil, notFoundErr
		},
	}

	svc := NewService(repo)
	_, err := svc.GetProfile(ctx, userID)

	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, notFoundErr) {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
}

// TestUserService_GetProfile_RepositoryError tests error propagation
func TestUserService_GetProfile_RepositoryError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	dbErr := errors.New("database connection failed")

	repo := &mockUserRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*User, error) {
			return nil, dbErr
		},
	}

	svc := NewService(repo)
	_, err := svc.GetProfile(ctx, userID)

	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected database error, got %v", err)
	}
}
