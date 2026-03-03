package validation

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
	"net/http"
)

var validate = validator.New()

// BindAndValidate is a helper that binds and validates request bodies in one call
func BindAndValidate[T any](c *gin.Context) (*T, bool) {
	var req T

	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.GenHTTPError(c, http.StatusBadRequest, "invalid request body", nil)
		return nil, false
	}

	if !validateRequest(c, &req) {
		return nil, false
	}
	return &req, true
}

// BindAndValidateURI is a helper that binds and validates path parameters in one call
func BindAndValidateURI[T any](c *gin.Context) (*T, bool) {
	var req T
	if err := c.ShouldBindUri(&req); err != nil {
		apperrors.GenHTTPError(c, http.StatusBadRequest, "invalid path parameter", nil)
		return nil, false
	}

	if !validateRequest(c, &req) {
		return nil, false
	}
	return &req, true
}

// BindAndValidateQuery is a helper that binds and validates query parameters in one call
func BindAndValidateQuery[T any](c *gin.Context) (*T, bool) {
	var req T
	if err := c.ShouldBindQuery(&req); err != nil {
		apperrors.GenHTTPError(c, http.StatusBadRequest, "invalid query parameter", nil)
		return nil, false
	}
	if !validateRequest(c, &req) {
		return nil, false
	}
	return &req, true
}

func validateRequest[T any](c *gin.Context, req *T) bool {
	if err := validate.Struct(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			validationErrors := make(map[string]string)
			for _, e := range ve {
				validationErrors[e.Field()] = validationErrorMessage(e)
			}
			apperrors.GenHTTPError(c, http.StatusBadRequest, "validation failed", validationErrors)
			return false
		}

		apperrors.GenHTTPError(c, http.StatusUnprocessableEntity, "validation failed", nil)
		return false
	}
	return true
}

func validationErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "this field is required"
	case "min":
		return "value is too short"
	case "max":
		return "value is too long"
	case "email":
		return "invalid email format"
	default:
		return "invalid value"
	}
}
