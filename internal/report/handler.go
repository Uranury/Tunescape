package report

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

// @Summary      Download PDF report
// @Description  Generates and streams a PDF report containing the user's top tracks, leaderboard rankings, and summary stats.
// @Tags         report
// @Produce      application/pdf
// @Success      200  {file}    binary
// @Failure      401  {object}  apperrors.HTTPError
// @Failure      404  {object}  apperrors.HTTPError  "User not found"
// @Failure      500  {object}  apperrors.HTTPError
// @Security     BearerAuth
// @Router       /me/report [get]
func (h *Handler) GetReport(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
		return
	}
	data, err := h.svc.GenerateReport(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.GenHTTPError(c, http.StatusNotFound, "user not found", nil)
			return
		}
		apperrors.GenHTTPError(c, http.StatusInternalServerError, "failed to generate report", nil)
		return
	}
	c.Header("Content-Disposition", `attachment; filename="tunescape-report.pdf"`)
	c.Data(http.StatusOK, "application/pdf", data)
}