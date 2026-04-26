package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/user"
	"golang.org/x/oauth2"
)

type mockSpotifyRepo struct {
	upsertFn      func(ctx context.Context, userID uuid.UUID, accessToken, refreshToken string, expiresAt time.Time) error
	getByUserIDFn func(ctx context.Context, userID uuid.UUID) (*Token, error)
}

func (m *mockSpotifyRepo) UpsertTokens(ctx context.Context, userID uuid.UUID, accessToken, refreshToken string, expiresAt time.Time) error {
	return m.upsertFn(ctx, userID, accessToken, refreshToken, expiresAt)
}

func (m *mockSpotifyRepo) DeleteByUserID(_ context.Context, _ uuid.UUID) error { return nil }

func (m *mockSpotifyRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*Token, error) {
	if m.getByUserIDFn != nil {
		return m.getByUserIDFn(ctx, userID)
	}
	return nil, nil
}

type mockUserRepo struct {
	connectSpotifyFn func(ctx context.Context, userID uuid.UUID, spotifyID *string, avatarURL, country, product *string) error
	createFn         func(ctx context.Context, u *user.User) error
	findByEmailFn    func(ctx context.Context, email string) (*user.User, error)
}

func (m *mockUserRepo) ConnectSpotify(ctx context.Context, userID uuid.UUID, spotifyID *string, avatarURL, country, product *string) error {
	return m.connectSpotifyFn(ctx, userID, spotifyID, avatarURL, country, product)
}

func (m *mockUserRepo) Create(ctx context.Context, u *user.User) error { return m.createFn(ctx, u) }
func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return m.findByEmailFn(ctx, email)
}
func (m *mockUserRepo) ClearSpotify(_ context.Context, _ uuid.UUID) error { return nil }
func (m *mockUserRepo) FindByID(_ context.Context, _ uuid.UUID) (*user.User, error) {
	return nil, nil
}
func (m *mockUserRepo) FindDisplayName(_ context.Context, _ uuid.UUID) (string, error) {
	return "", nil
}
func (m *mockUserRepo) FindDisplayNamesByIDs(_ context.Context, _ []string) (map[string]string, error) {
	return map[string]string{}, nil
}

func (m *mockUserRepo) FindAll(_ context.Context) ([]user.User, error) {
	return nil, nil
}

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestSpotifyService_AuthURL(t *testing.T) {
	t.Parallel()

	codeState := "state123"

	c := &Client{
		oauth2Cfg: &oauth2.Config{
			ClientID:     "cid",
			ClientSecret: "sec",
			RedirectURL:  "http://localhost/callback",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "http://example.com/auth",
				TokenURL: "http://example.com/token",
			},
		},
	}

	repo := &mockSpotifyRepo{}
	userRepo := &mockUserRepo{}

	svc := NewService(repo, userRepo, c, nil, slog.Default())
	got := svc.AuthURL(codeState)

	want := c.oauth2Cfg.AuthCodeURL(codeState)
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestSpotifyService_UpsertTokens(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	accessToken := "at"
	refreshToken := "rt"
	expiresAt := time.Now().Add(1 * time.Hour)

	called := false
	repo := &mockSpotifyRepo{
		upsertFn: func(ctx context.Context, gotUserID uuid.UUID, gotAccessToken, gotRefreshToken string, gotExpiresAt time.Time) error {
			called = true
			if gotUserID != userID {
				t.Fatalf("expected userID %s, got %s", userID, gotUserID)
			}
			if gotAccessToken != accessToken {
				t.Fatalf("expected accessToken %q, got %q", accessToken, gotAccessToken)
			}
			if gotRefreshToken != refreshToken {
				t.Fatalf("expected refreshToken %q, got %q", refreshToken, gotRefreshToken)
			}
			if !gotExpiresAt.Equal(expiresAt) {
				t.Fatalf("expected expiresAt %v, got %v", expiresAt, gotExpiresAt)
			}
			return nil
		},
	}

	svc := NewService(repo, &mockUserRepo{}, &Client{}, nil, slog.Default())
	if err := svc.UpsertTokens(ctx, userID, accessToken, refreshToken, expiresAt); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !called {
		t.Fatalf("expected UpsertTokens to call repository")
	}
}

func TestSpotifyService_ConnectAccount_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	authCode := "code123"

	httpClient := &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			switch {
			case req.URL.Host == "token.test" && req.URL.Path == "/token":
				if err := req.ParseForm(); err != nil {
					return nil, err
				}
				if code := req.PostForm.Get("code"); code != authCode {
					return nil, fmt.Errorf("expected code %q, got %q", authCode, code)
				}
				body := `{"access_token":"access-token","refresh_token":"refresh-token","expires_in":3600}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(body)),
					Request:    req,
				}, nil
			case req.URL.Host == "api.spotify.com" && req.URL.Path == "/v1/me":
				if got := req.Header.Get("Authorization"); got != "Bearer access-token" {
					return nil, fmt.Errorf("unexpected Authorization header: %q", got)
				}
				me := map[string]any{
					"id":           "spotify-id",
					"email":        "x@y.z",
					"display_name": "Disp Name",
					"country":      "US",
					"product":      "Premium",
					"images":       []map[string]any{{"url": "https://img/avatar.png"}},
				}
				b, _ := json.Marshal(me)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(string(b))),
					Request:    req,
				}, nil
			default:
				return nil, fmt.Errorf("unexpected external request: %s %s", req.URL.Host, req.URL.Path)
			}
		}),
	}

	c := &Client{
		httpClient: httpClient,
		oauth2Cfg: &oauth2.Config{
			ClientID:     "cid",
			ClientSecret: "sec",
			RedirectURL:  "http://localhost/callback",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "http://token.test/auth",
				TokenURL: "http://token.test/token",
			},
		},
	}

	var (
		mu           sync.Mutex
		gotSpotifyID *string
		gotAvatarURL *string
		gotCountry   *string
		gotProduct   *string
		upsertCalled bool
		gotAccess    string
		gotRefresh   string
		gotExpiresAt time.Time
	)

	repo := &mockSpotifyRepo{
		upsertFn: func(ctx context.Context, gotUserID uuid.UUID, accessToken, refreshToken string, expiresAt time.Time) error {
			mu.Lock()
			defer mu.Unlock()
			upsertCalled = true
			if gotUserID != userID {
				t.Fatalf("expected userID %s, got %s", userID, gotUserID)
			}
			gotAccess = accessToken
			gotRefresh = refreshToken
			gotExpiresAt = expiresAt
			return nil
		},
	}

	userRepo := &mockUserRepo{
		connectSpotifyFn: func(ctx context.Context, gotUserID uuid.UUID, spotifyID *string, avatarURL, country, product *string) error {
			mu.Lock()
			defer mu.Unlock()
			if gotUserID != userID {
				t.Fatalf("expected userID %s, got %s", userID, gotUserID)
			}
			gotSpotifyID = spotifyID
			gotAvatarURL = avatarURL
			gotCountry = country
			gotProduct = product
			return nil
		},
	}

	svc := NewService(repo, userRepo, c, nil, slog.Default())
	before := time.Now()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	if err := svc.ConnectAccount(ctx, userID, authCode); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	after := time.Now()

	mu.Lock()
	defer mu.Unlock()

	if gotSpotifyID == nil || *gotSpotifyID != "spotify-id" {
		t.Fatalf("unexpected spotifyID: %#v", gotSpotifyID)
	}
	if gotAvatarURL == nil || *gotAvatarURL != "https://img/avatar.png" {
		t.Fatalf("unexpected avatarURL: %#v", gotAvatarURL)
	}
	if gotCountry == nil || *gotCountry != "US" {
		t.Fatalf("unexpected country: %#v", gotCountry)
	}
	if gotProduct == nil || *gotProduct != "Premium" {
		t.Fatalf("unexpected product: %#v", gotProduct)
	}

	if !upsertCalled {
		t.Fatalf("expected UpsertTokens to be called")
	}
	if gotAccess != "access-token" {
		t.Fatalf("unexpected access token: %q", gotAccess)
	}
	if gotRefresh != "refresh-token" {
		t.Fatalf("unexpected refresh token: %q", gotRefresh)
	}

	lower := before.Add(3600*time.Second - 10*time.Second)
	upper := after.Add(3600*time.Second + 10*time.Second)
	if gotExpiresAt.Before(lower) || gotExpiresAt.After(upper) {
		t.Fatalf("unexpected expiresAt: %v (expected in [%v, %v])", gotExpiresAt, lower, upper)
	}
}

func TestSpotifyService_ConnectAccount_TokenExchangeError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	authCode := "bad-code"

	httpClient := &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Host == "token.test" && req.URL.Path == "/token" {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Header:     http.Header{"Content-Type": []string{"text/plain"}},
					Body:       io.NopCloser(strings.NewReader("bad request")),
					Request:    req,
				}, nil
			}
			return nil, fmt.Errorf("should not call /me when token exchange fails")
		}),
	}

	c := &Client{
		httpClient: httpClient,
		oauth2Cfg: &oauth2.Config{
			ClientID:     "cid",
			ClientSecret: "sec",
			RedirectURL:  "http://localhost/callback",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "http://token.test/auth",
				TokenURL: "http://token.test/token",
			},
		},
	}

	svc := NewService(&mockSpotifyRepo{}, &mockUserRepo{}, c, nil, slog.Default())
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	err := svc.ConnectAccount(ctx, userID, authCode)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "exchange oauth code") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSpotifyService_ConnectAccount_GetMeError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	authCode := "code123"

	httpClient := &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			switch {
			case req.URL.Host == "token.test" && req.URL.Path == "/token":
				body := `{"access_token":"access-token","refresh_token":"refresh-token","expires_in":3600}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(body)),
					Request:    req,
				}, nil
			case req.URL.Host == "api.spotify.com" && req.URL.Path == "/v1/me":
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"error":{"status":401,"message":"Unauthorized"}}`)),
					Request:    req,
				}, nil
			default:
				return nil, fmt.Errorf("unexpected request: %s %s", req.URL.Host, req.URL.Path)
			}
		}),
	}

	c := &Client{
		httpClient: httpClient,
		oauth2Cfg: &oauth2.Config{
			ClientID:     "cid",
			ClientSecret: "sec",
			RedirectURL:  "http://localhost/callback",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "http://token.test/auth",
				TokenURL: "http://token.test/token",
			},
		},
	}

	repo := &mockSpotifyRepo{
		upsertFn: func(_ context.Context, _ uuid.UUID, _, _ string, _ time.Time) error {
			t.Fatal("UpsertTokens must not be called when /v1/me fails")
			return nil
		},
	}
	userRepo := &mockUserRepo{
		connectSpotifyFn: func(_ context.Context, _ uuid.UUID, _ *string, _, _, _ *string) error {
			t.Fatal("ConnectSpotify must not be called when /v1/me fails")
			return nil
		},
	}

	svc := NewService(repo, userRepo, c, nil, slog.Default())
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	err := svc.ConnectAccount(ctx, userID, authCode)

	if err == nil {
		t.Fatal("expected error when /v1/me fails, got nil")
	}
	if !strings.Contains(err.Error(), "fetch spotify profile") {
		t.Fatalf("expected 'fetch spotify profile' in error, got: %v", err)
	}
}

func TestParseTimeRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input   string
		want    TimeRange
		wantErr bool
	}{
		{"short_term", ShortTerm, false},
		{"medium_term", MediumTerm, false},
		{"long_term", LongTerm, false},
		{"", MediumTerm, false},
		{"invalid", "", true},
		{"MEDIUM_TERM", "", true},
		{"all_time", "", true},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseTimeRange(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for input %q, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func TestSpotifyService_GetTopTracks_TimeRangeInURL(t *testing.T) {
	t.Parallel()

	for _, timeRange := range []TimeRange{ShortTerm, MediumTerm, LongTerm} {
		tr := timeRange
		t.Run(string(tr), func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			userID := uuid.New()
			accessToken := "test-token"

			tracks := []map[string]any{
				{"id": "s1", "name": "Track 1", "popularity": 80},
			}
			body, _ := json.Marshal(map[string]any{"items": tracks})

			var capturedURL string
			httpClient := &http.Client{
				Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
					capturedURL = req.URL.String()
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body:       io.NopCloser(strings.NewReader(string(body))),
						Request:    req,
					}, nil
				}),
			}

			c := &Client{
				httpClient: httpClient,
				oauth2Cfg:  &oauth2.Config{},
			}

			repo := &mockSpotifyRepo{
				getByUserIDFn: func(_ context.Context, _ uuid.UUID) (*Token, error) {
					return &Token{
						AccessToken: accessToken,
						ExpiresAt:   time.Now().Add(time.Hour),
					}, nil
				},
			}

			svc := NewService(repo, &mockUserRepo{}, c, nil, slog.Default())
			result, err := svc.GetTopTracks(ctx, userID, 50, tr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 {
				t.Fatalf("expected 1 track, got %d", len(result))
			}
			wantParam := "time_range=" + string(tr)
			if !strings.Contains(capturedURL, wantParam) {
				t.Fatalf("expected URL to contain %q, got %q", wantParam, capturedURL)
			}
		})
	}
}
