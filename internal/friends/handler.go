package friends

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/middleware"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
	"gitlab.com/Uranury/tunescape/pkg/validation"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

type sendRequestBody struct {
	ReceiverID string `json:"receiver_id" validate:"required,uuid"`
}

// @Summary      Send a friend request
// @Tags         friends
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body sendRequestBody true "receiver"
// @Success      201
// @Failure      400  {object}  apperrors.HTTPError
// @Failure      401  {object}  apperrors.HTTPError
// @Failure      409  {object}  apperrors.HTTPError
// @Failure      500  {object}  apperrors.HTTPError
// @Router       /friends/requests [post]
func (h *Handler) SendRequest(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
		return
	}

	body, ok := validation.BindAndValidate[sendRequestBody](c)
	if !ok {
		return
	}

	receiverID, err := uuid.Parse(body.ReceiverID)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusBadRequest, "invalid receiver_id", nil)
		return
	}

	if err := h.svc.SendRequest(c.Request.Context(), userID, receiverID); err != nil {
		h.handleError(c, err)
		return
	}
	c.Status(http.StatusCreated)
}

// @Summary      List incoming friend requests
// @Tags         friends
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  IncomingRequest
// @Failure      401  {object}  apperrors.HTTPError
// @Failure      500  {object}  apperrors.HTTPError
// @Router       /friends/requests [get]
func (h *Handler) ListIncoming(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
		return
	}

	reqs, err := h.svc.ListIncoming(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, reqs)
}

// @Summary      Accept a friend request
// @Tags         friends
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Request ID"
// @Success      200
// @Failure      401  {object}  apperrors.HTTPError
// @Failure      403  {object}  apperrors.HTTPError
// @Failure      404  {object}  apperrors.HTTPError
// @Failure      500  {object}  apperrors.HTTPError
// @Router       /friends/requests/{id}/accept [post]
func (h *Handler) AcceptRequest(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusBadRequest, "invalid request id", nil)
		return
	}

	if err := h.svc.AcceptRequest(c.Request.Context(), id, userID); err != nil {
		h.handleError(c, err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary      Reject a friend request
// @Tags         friends
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Request ID"
// @Success      204
// @Failure      401  {object}  apperrors.HTTPError
// @Failure      404  {object}  apperrors.HTTPError
// @Failure      500  {object}  apperrors.HTTPError
// @Router       /friends/requests/{id}/reject [post]
func (h *Handler) RejectRequest(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusBadRequest, "invalid request id", nil)
		return
	}

	if err := h.svc.RejectRequest(c.Request.Context(), id, userID); err != nil {
		h.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary      List friends
// @Tags         friends
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  FriendProfile
// @Failure      401  {object}  apperrors.HTTPError
// @Failure      500  {object}  apperrors.HTTPError
// @Router       /friends [get]
func (h *Handler) ListFriends(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
		return
	}

	friends, err := h.svc.ListFriends(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, friends)
}

// @Summary      Compare music tastes with a friend
// @Description  Returns averaged audio feature scores and a 0–100 compatibility score. Returns 422 if the friend has no listening data yet.
// @Tags         friends
// @Produce      json
// @Security     BearerAuth
// @Param        friend_id path string true "Friend's user ID"
// @Success      200  {object}  TasteComparison
// @Failure      401  {object}  apperrors.HTTPError
// @Failure      403  {object}  apperrors.HTTPError
// @Failure      422  {object}  apperrors.HTTPError
// @Failure      500  {object}  apperrors.HTTPError
// @Router       /friends/{friend_id}/compare [get]
func (h *Handler) CompareTastes(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
		return
	}

	friendID, err := uuid.Parse(c.Param("friend_id"))
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusBadRequest, "invalid friend_id", nil)
		return
	}

	result, err := h.svc.CompareTastes(c.Request.Context(), userID, friendID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      Get a friend's playlists
// @Description  Returns playlists created via Tunescape. Returns 422 if the friend has not connected Spotify yet.
// @Tags         friends
// @Produce      json
// @Security     BearerAuth
// @Param        friend_id path string true "Friend's user ID"
// @Success      200  {array}  gitlab_com_Uranury_tunescape_internal_playlist.Playlist
// @Failure      401  {object}  apperrors.HTTPError
// @Failure      403  {object}  apperrors.HTTPError
// @Failure      422  {object}  apperrors.HTTPError
// @Failure      500  {object}  apperrors.HTTPError
// @Router       /friends/{friend_id}/playlists [get]
func (h *Handler) GetFriendPlaylists(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
		return
	}

	friendID, err := uuid.Parse(c.Param("friend_id"))
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusBadRequest, "invalid friend_id", nil)
		return
	}

	playlists, err := h.svc.GetFriendPlaylists(c.Request.Context(), userID, friendID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, playlists)
}

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperrors.ErrCannotAddSelf):
		apperrors.GenHTTPError(c, http.StatusBadRequest, err.Error(), nil)
	case errors.Is(err, apperrors.ErrRequestAlreadySent), errors.Is(err, apperrors.ErrAlreadyFriends):
		apperrors.GenHTTPError(c, http.StatusConflict, err.Error(), nil)
	case errors.Is(err, apperrors.ErrRequestNotFound):
		apperrors.GenHTTPError(c, http.StatusNotFound, err.Error(), nil)
	case errors.Is(err, apperrors.ErrNotFriends):
		apperrors.GenHTTPError(c, http.StatusForbidden, err.Error(), nil)
	case errors.Is(err, apperrors.ErrSpotifyNotConnected), errors.Is(err, apperrors.ErrNoSnapshot):
		apperrors.GenHTTPError(c, http.StatusUnprocessableEntity,
			"this friend has not connected Spotify yet or has no listening data", nil)
	case errors.Is(err, apperrors.ErrForbidden):
		apperrors.GenHTTPError(c, http.StatusForbidden, apperrors.ErrForbidden.Error(), nil)
	default:
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "internal server error", nil)
	}
}
