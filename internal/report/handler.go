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