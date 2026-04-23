package snapshot

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		switch {
		case errors.Is(err, apperrors.ErrSpotifyNotConnected):
			apperrors.GenHTTPError(c, http.StatusUnprocessableEntity, apperrors.ErrSpotifyNotConnected.Error(), nil)
		case errors.Is(err, apperrors.ErrUpstreamUnavailable):
			apperrors.GenHTTPError(c, http.StatusBadGateway, "Spotify is temporarily unavailable, please try again", nil)
		default:
			apperrors.GenHTTPError(c, http.StatusInternalServerError, "failed to create snapshot", nil)
		}
		return
	}

	c.JSON(http.StatusCreated, snap)
}

// @Summary      List snapshots
// @Description  Returns all snapshots for the authenticated user, ordered by most recent first.
// @Tags         snapshots
// @Produce      json
// @Success      200  {array}   SnapshotSummary
// @Failure      401  {object}  apperrors.HTTPError
// @Failure      500  {object}  apperrors.HTTPError
// @Router       /me/snapshots [get]
func (h *Handler) ListSnapshots(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
		return
	}

	summaries, err := h.svc.ListSnapshots(c.Request.Context(), userID)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "failed to list snapshots", nil)
		return
	}

	c.JSON(http.StatusOK, summaries)
}

// @Summary      Get snapshot by ID
// @Description  Returns a single snapshot with its full track list. Returns 404 if the snapshot does not exist or belongs to another user.
// @Tags         snapshots
// @Produce      json
// @Param        id   path      string  true  "Snapshot UUID"
// @Success      200  {object}  Snapshot
// @Failure      400  {object}  apperrors.HTTPError  "Invalid UUID"
// @Failure      401  {object}  apperrors.HTTPError
// @Failure      404  {object}  apperrors.HTTPError
// @Failure      500  {object}  apperrors.HTTPError
// @Router       /me/snapshots/{id} [get]
func (h *Handler) GetSnapshot(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
		return
	}

	snapshotID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusBadRequest, "invalid snapshot id", nil)
		return
	}

	snap, err := h.svc.GetSnapshot(c.Request.Context(), userID, snapshotID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.GenHTTPError(c, http.StatusNotFound, apperrors.ErrNotFound.Error(), nil)
			return
		}
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "failed to get snapshot", nil)
		return
	}

	c.JSON(http.StatusOK, snap)
}