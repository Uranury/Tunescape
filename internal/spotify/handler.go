package spotify

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/auth"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
)

type Handler struct {
	svc            Service
	authRefreshSvc auth.RefreshTokenService
	logger         *slog.Logger
	secureCookie   bool
	frontendURL    string
}

func NewHandler(
	svc Service,
	authRefreshSvc auth.RefreshTokenService,
	logger *slog.Logger,
	secureCookie bool,
	frontendURL string,
) *Handler {
	return &Handler{
		svc:            svc,
		authRefreshSvc: authRefreshSvc,
		logger:         logger,
		secureCookie:   secureCookie,
		frontendURL:    frontendURL,
	}
}

// @Summary Spotify OAuth login
// @Description Redirects to Spotify authorization endpoint and sets oauth_state cookie.
// @Tags spotify
// @Success 302 "Found (redirect to Spotify authorization endpoint)"
// @Router /auth/spotify/login [get]
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
	c.Redirect(http.StatusFound, h.svc.AuthURL(state))
}

// @Summary Spotify OAuth callback
// @Description Handles Spotify OAuth callback, links Spotify account to the current user, and redirects back to the frontend.
// @Tags spotify
// @Param code query string true "Spotify authorization code"
// @Param state query string true "Spotify state"
// @Param oauth_state header string true "oauth_state cookie value"
// @Param refresh_token header string true "refresh_token cookie value"
// @Success 302 "Found (redirect to frontend with connected=1 or error=...)"
// @Router /auth/spotify/callback [get]
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

	if err := h.svc.ConnectAccount(c.Request.Context(), existingToken.UserID, code); err != nil {
		if errors.Is(err, apperrors.ErrSpotifyIDTaken) {
			errRedirect("spotify_already_linked")
			return
		}
		h.logger.Error("failed to connect spotify account", "err", err, "user_id", existingToken.UserID)
		errRedirect("db_error")
		return
	}

	c.Redirect(http.StatusFound, h.frontendURL+"?connected=1")
}

// @Summary      Disconnect Spotify account
// @Description  Removes the Spotify connection and tokens for the authenticated user.
// @Tags         spotify
// @Produce      json
// @Success      204
// @Failure      401  {object}  apperrors.HTTPError
// @Failure      500  {object}  apperrors.HTTPError
// @Router       /me/spotify [delete]
func (h *Handler) DisconnectHandler(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
		return
	}
	userID, ok := val.(uuid.UUID)
	if !ok {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
		return
	}

	if err := h.svc.Disconnect(c.Request.Context(), userID); err != nil {
		h.logger.Error("failed to disconnect spotify", "user_id", userID, "error", err)
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "failed to disconnect spotify", nil)
		return
	}

	c.Status(http.StatusNoContent)
}

func generateRandomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
