package leaderboard

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetLeaderboard(c *gin.Context) {
	feature := c.Param("feature")
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)
	if err != nil || offset < 0 {
		offset = 0
	}
	result, err := h.svc.GetLeaderboard(c.Request.Context(), feature, limit, offset)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, result)
}
