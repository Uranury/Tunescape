package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/Uranury/tunescape/internal/user"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
	"gitlab.com/Uranury/tunescape/pkg/validation"
)

type Handler struct {
	refreshTokenSvc RefreshTokenService
	tokenSvc        TokenService
	userSvc         user.Service
	secureCookie    bool
}

func NewHandler(refreshTokenSvc RefreshTokenService, tokenSvc TokenService, userSvc user.Service, secureCookie bool) *Handler {
	return &Handler{
		refreshTokenSvc: refreshTokenSvc,
		tokenSvc:        tokenSvc,
		userSvc:         userSvc,
		secureCookie:    secureCookie,
	}
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type signupRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	DisplayName string `json:"display_name" validate:"required,min=1"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (h *Handler) Login(c *gin.Context) {
	req, ok := validation.BindAndValidate[loginRequest](c)
	if !ok {
		return
	}

	u, err := h.userSvc.ValidateCredentials(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidCredentials) {
			apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
			return
		}
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "internal error", nil)
		return
	}

	accessToken, err := h.tokenSvc.Generate(u.ID, u.Role)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "internal error", nil)
		return
	}

	refreshToken, err := h.refreshTokenSvc.Generate(c.Request.Context(), u.ID, u.Role, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "internal error", nil)
		return
	}

	setRefreshCookie(c, refreshToken, h.secureCookie)
	c.JSON(http.StatusOK, tokenResponse{AccessToken: accessToken})
}

func (h *Handler) Signup(c *gin.Context) {
	req, ok := validation.BindAndValidate[signupRequest](c)
	if !ok {
		return
	}

	u, err := h.userSvc.Create(c.Request.Context(), req.Email, req.Password, req.DisplayName)
	if err != nil {
		if errors.Is(err, apperrors.ErrEmailTaken) {
			apperrors.GenHTTPError(c, http.StatusConflict, apperrors.ErrEmailTaken.Error(), nil)
			return
		}
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "internal error", nil)
		return
	}

	accessToken, err := h.tokenSvc.Generate(u.ID, u.Role)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "internal error", nil)
		return
	}

	refreshToken, err := h.refreshTokenSvc.Generate(c.Request.Context(), u.ID, u.Role, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "internal error", nil)
		return
	}

	setRefreshCookie(c, refreshToken, h.secureCookie)
	c.JSON(http.StatusCreated, tokenResponse{AccessToken: accessToken})
}

func (h *Handler) Refresh(c *gin.Context) {
	cookie, err := c.Request.Cookie("refresh_token")
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
		return
	}

	accessToken, newRefreshToken, err := h.refreshTokenSvc.Refresh(c.Request.Context(), cookie.Value)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
		return
	}

	setRefreshCookie(c, newRefreshToken, h.secureCookie)
	c.JSON(http.StatusOK, tokenResponse{AccessToken: accessToken})
}

func (h *Handler) Logout(c *gin.Context) {
	cookie, err := c.Request.Cookie("refresh_token")
	if err != nil {
		c.Status(http.StatusNoContent)
		return
	}

	_ = h.refreshTokenSvc.Revoke(c.Request.Context(), cookie.Value)

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	})
	c.Status(http.StatusNoContent)
}

func setRefreshCookie(c *gin.Context, token string, secure bool) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		MaxAge:   int(RefreshTokenTTL.Seconds()),
		Path:     "/",
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
