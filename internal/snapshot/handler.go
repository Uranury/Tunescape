package snapshot

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

// @Summary      Create top-tracks snapshot
// @Description  Fetches the authenticated user's top 50 tracks from Spotify,
//               persists them, and returns the new snapshot with its track list.
// @Tags         snapshots
// @Produce      json
// @Success      201  {object}  Snapshot
// @Failure      401  {object}  apperrors.HTTPError
// @Failure      422  {object}  apperrors.HTTPError  "Spotify account not connected"
// @Failure      500  {object}  apperrors.HTTPError
// @Router       /me/snapshots [post]
func (h *Handler) CreateSnapshot(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
		return
	}

	snap, err := h.svc.CreateSnapshot(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrSpotifyNotConnected) {
			apperrors.GenHTTPError(c, http.StatusUnprocessableEntity, apperrors.ErrSpotifyNotConnected.Error(), nil)
			return
		}
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "failed to create snapshot", nil)
		return
	}

	c.JSON(http.StatusCreated, snap)
}