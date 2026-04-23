package trends

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/Uranury/tunescape/internal/middleware"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// @Summary      Get listening trends
// @Description  Returns a time-ordered series of audio feature averages, one point per snapshot.
// @Tags         trends
// @Produce      json
// @Success      200  {object}  TrendsResponse
// @Failure      401  {object}  apperrors.HTTPError
// @Failure      500  {object}  apperrors.HTTPError
// @Security     BearerAuth
// @Router       /me/trends [get]
func (h *Handler) GetTrends(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
		return
	}
	result, err := h.svc.GetTrends(c.Request.Context(), userID)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "failed to get trends", nil)
		return
	}
	c.JSON(http.StatusOK, result)
}