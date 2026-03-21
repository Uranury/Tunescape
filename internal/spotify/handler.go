package spotify

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"log/slog"
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
	"gitlab.com/Uranury/tunescape/internal/auth"
	"gitlab.com/Uranury/tunescape/internal/user"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
)

type Handler struct {
	client         *Client
	svc            *Service
	authRefreshSvc auth.RefreshTokenService
	userRepo       user.Repository
	logger         *slog.Logger
	secureCookie   bool
	frontendURL    string
}

func NewHandler(
	client *Client,
	svc *Service,
	authRefreshSvc auth.RefreshTokenService,
	userRepo user.Repository,
	logger *slog.Logger,
	secureCookie bool,
	frontendURL string,
) *Handler {
	return &Handler{
		client:         client,
		svc:            svc,
		authRefreshSvc: authRefreshSvc,
		userRepo:       userRepo,
		logger:         logger,
		secureCookie:   secureCookie,
		frontendURL:    frontendURL,
	}
}

func (h *Handler) LoginHandler(c *gin.Context) {
	state, err := generateRandomState()
	if err != nil {
		h.logger.Error("failed to generate oauth state", "err", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		MaxAge:   300,
		Path:     "/",
		Secure:   h.secureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	c.Redirect(http.StatusFound, h.client.oauth2Cfg.AuthCodeURL(state))
}

func (h *Handler) CallbackHandler(c *gin.Context) {
	errRedirect := func(msg string) {
		c.Redirect(http.StatusFound, h.frontendURL+"?error="+msg)
	}

	code := c.Query("code")
	if code == "" {
		errRedirect("missing_code")
		return
	}

	state := c.Query("state")
	stateCookie, err := c.Request.Cookie("oauth_state")
	if err != nil || subtle.ConstantTimeCompare([]byte(stateCookie.Value), []byte(state)) != 1 {
		errRedirect("invalid_state")
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Secure:   h.secureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	refreshCookie, err := c.Request.Cookie("refresh_token")
	if err != nil {
		errRedirect("not_logged_in")
		return
	}

	existingToken, err := h.authRefreshSvc.Validate(c.Request.Context(), refreshCookie.Value)
	if err != nil {
		errRedirect("session_expired")
		return
	}

	token, err := h.client.oauth2Cfg.Exchange(c.Request.Context(), code)
	if err != nil {
		h.logger.Error("failed to exchange oauth code", "err", err)
		errRedirect("spotify_exchange_failed")
		return
	}

	profile, err := h.client.getMe(c.Request.Context(), token.AccessToken)
	if err != nil {
		h.logger.Error("failed to fetch spotify profile", "err", err)
		errRedirect("spotify_profile_failed")
		return
	}

	spotifyID := profile.ID
	var avatarURL, country, product *string
	if len(profile.Images) > 0 {
		avatarURL = &profile.Images[0].URL
	}
	if profile.Country != "" {
		country = &profile.Country
	}
	if profile.Product != "" {
		product = &profile.Product
	}

	if err := h.userRepo.ConnectSpotify(c.Request.Context(), existingToken.UserID, &spotifyID, avatarURL, country, product); err != nil {
		if errors.Is(err, apperrors.ErrSpotifyIDTaken) {
			errRedirect("spotify_already_linked")
			return
		}
		h.logger.Error("failed to connect spotify to user", "err", err, "user_id", existingToken.UserID)
		errRedirect("db_error")
		return
	}

	if err := h.svc.UpsertTokens(c.Request.Context(), existingToken.UserID, token.AccessToken, token.RefreshToken, token.Expiry); err != nil {
		h.logger.Error("failed to upsert spotify tokens", "err", err)
		errRedirect("db_error")
		return
	}

	c.Redirect(http.StatusFound, h.frontendURL+"?connected=1")
}

func generateRandomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
