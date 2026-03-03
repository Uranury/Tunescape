package apperrors

import (
	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

func GenHTTPError(c *gin.Context, code int, message string, details any) {
	response := HTTPError{code, message, details}
	accept := c.Request.Header.Get("Accept")
	switch accept {
	case gin.MIMEJSON:
		c.JSON(code, response)
	case gin.MIMEXML:
		c.XML(code, response)
	case gin.MIMEYAML:
		c.YAML(code, response)
	default:
		c.JSON(code, response)
	}
}
