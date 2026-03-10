package auth

import "github.com/gin-gonic/gin"

type Handler struct {
	refreshTokenSvc RefreshTokenService
	tokenSvc        TokenService
}

func NewHandler(refreshTokenSvc RefreshTokenService, tokenSvc TokenService) *Handler {
	return &Handler{
		refreshTokenSvc: refreshTokenSvc,
		tokenSvc:        tokenSvc,
	}
}

func (h *Handler) Login(c *gin.Context) {

}
