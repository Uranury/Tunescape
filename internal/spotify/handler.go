package spotify

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	client       *Client
	svc          *Service
	secureCookie bool
}

func NewHandler(client *Client, svc *Service, secureCookie bool) *Handler {
	return &Handler{
		client:       client,
		svc:          svc,
		secureCookie: secureCookie,
	}
}

func (h *Handler) LoginHandler(c *gin.Context) {
	state, err := generateRandomState()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		MaxAge:   300,
		Path:     "/",
		Domain:   "",
		Secure:   h.secureCookie,
		HttpOnly: true,
	})
	url := h.client.oauth2Cfg.AuthCodeURL(state)
	c.Redirect(http.StatusFound, url)
}

func (h *Handler) CallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	state := c.Query("state")
	cookie, err := c.Request.Cookie("oauth_state")
	if err != nil || subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(state)) != 1 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// clear the state cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})

	token, err := h.client.oauth2Cfg.Exchange(c.Request.Context(), code)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := h.svc.SaveAccess(c.Request.Context(), token.AccessToken); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err := h.svc.SaveRefresh(c.Request.Context(), token.RefreshToken); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Redirect(http.StatusFound, "/")

	token.Expiry
	// token.AccessToken, token.RefreshToken, token.Expiry
	// now save to DB and redirect to frontend
}

func generateRandomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	state := base64.URLEncoding.EncodeToString(b)
	return state, nil
}
