package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// @Summary      Get current user profile
// @Description  Returns the authenticated user's display name, email, avatar, and Spotify connection status.
// @Tags         user
// @Produce      json
// @Success      200  {object}  ProfileResponse
// @Failure      401  {object}  apperrors.HTTPError
// @Failure      404  {object}  apperrors.HTTPError
// @Router       /me/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
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

	profile, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.GenHTTPError(c, http.StatusNotFound, "user not found", nil)
			return
		}
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "internal error", nil)
		return
	}

	c.JSON(http.StatusOK, profile)
}
