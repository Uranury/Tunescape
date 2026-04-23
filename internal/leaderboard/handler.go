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

// @Summary      Get global leaderboard
// @Description  Returns the top-ranked users for a given audio feature. Valid features: valence, energy, danceability, acousticness.
// @Tags         leaderboards
// @Produce      json
// @Param        feature  path   string  true   "Audio feature"  Enums(valence, energy, danceability, acousticness)
// @Param        limit    query  int     false  "Number of entries to return (default 10)"
// @Param        offset   query  int     false  "Number of entries to skip (default 0)"
// @Success      200  {object}  LeaderboardResponse
// @Failure      400  {object}  apperrors.HTTPError  "Invalid feature"
// @Router       /leaderboards/{feature} [get]
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
