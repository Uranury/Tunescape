package analytics

import (
	"errors"
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

// @Summary      Get music taste based on top tracks
// @Description  Reads the authenticated user's latest snapshot, fetches audio features
//
//	from Reccobeats for those tracks, stores them, and returns aggregated
//	averages (danceability, valence, energy, etc.).
//	Returns 404 if the user has no snapshots yet.
//
// @Tags         analytics
// @Produce      json
// @Success      200  {object}  MusicTasteResponse
// @Failure      401  {object}  apperrors.HTTPError
// @Failure      404  {object}  apperrors.HTTPError  "No snapshot found for user"
// @Failure      500  {object}  apperrors.HTTPError
// @Router       /analytics/me/top-tracks [get]
func (h *Handler) GetMusicTaste(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
		return
	}

	result, err := h.svc.GetMusicTaste(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNoSnapshot) {
			apperrors.GenHTTPError(c, http.StatusNotFound, apperrors.ErrNoSnapshot.Error(), nil)
			return
		}
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "failed to get music taste", nil)
		return
	}

	c.JSON(http.StatusOK, result)
}
