package spotify

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
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

func generateRandomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
