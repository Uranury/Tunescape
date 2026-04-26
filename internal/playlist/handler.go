package playlist

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

// @Summary      Create Spotify playlist from latest snapshot
// @Description  Creates a private Spotify playlist containing all tracks from the
//
//	authenticated user's most recent snapshot and returns the playlist
//	metadata including a direct URL and embed URL.
//	Returns 404 if the user has no snapshots yet.
//
// @Tags         playlists
// @Produce      json
// @Success      201  {object}  Response
// @Failure      401  {object}  apperrors.HTTPError
// @Failure      404  {object}  apperrors.HTTPError  "No snapshot found — create one first"
// @Failure      422  {object}  apperrors.HTTPError  "Spotify account not connected"
// @Failure      502  {object}  apperrors.HTTPError  "Spotify temporarily unavailable"
// @Failure      500  {object}  apperrors.HTTPError
// @Router       /me/playlists/top-tracks [post]
func (h *Handler) CreateFromSnapshot(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
		return
	}

	result, err := h.svc.CreateFromLatestSnapshot(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNoSnapshot):
			apperrors.GenHTTPError(c, http.StatusNotFound, "no snapshot found, create one first", nil)
		case errors.Is(err, apperrors.ErrSpotifyNotConnected):
			apperrors.GenHTTPError(c, http.StatusUnprocessableEntity, apperrors.ErrSpotifyNotConnected.Error(), nil)
		case errors.Is(err, apperrors.ErrUpstreamUnavailable):
			apperrors.GenHTTPError(c, http.StatusBadGateway, "Spotify is temporarily unavailable, please try again", nil)
		default:
			apperrors.GenHTTPError(c, http.StatusInternalServerError, "failed to create playlist", nil)
		}
		return
	}

	c.JSON(http.StatusCreated, result)
}
