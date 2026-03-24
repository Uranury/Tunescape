package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/Uranury/tunescape/internal/auth"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
)

type Auth struct {
	authService auth.TokenService
}

func NewAuth(authService auth.TokenService) *Auth {
	return &Auth{authService: authService}
}

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

func (m *Auth) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, err := GetUserRole(c)
		if err != nil {
			apperrors.GenHTTPError(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error(), nil)
			c.Abort()
			return
		}

		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		apperrors.GenHTTPError(c, http.StatusForbidden, apperrors.ErrForbidden.Error(), nil)
		c.Abort()
	}
}

func (m *Auth) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			apperrors.GenHTTPError(c, http.StatusUnauthorized, "invalid token", nil)
			c.Abort()
			return
		}
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		claims, err := m.authService.Validate(tokenString)
		if err != nil {
			apperrors.GenHTTPError(c, http.StatusUnauthorized, "invalid token", nil)
			c.Abort()
			return
		}

		c.Set(string(UserIDKey), claims.UserID)
		c.Set(string(RoleKey), claims.Role)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (uuid.UUID, error) {
	val, exists := c.Get(string(UserIDKey))
	if !exists {
		return uuid.Nil, errors.New("user id not found in context")
	}
	id, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user id has invalid type")
	}
	return id, nil
}

func GetUserRole(c *gin.Context) (string, error) {
	val, exists := c.Get(string(RoleKey))
	if !exists {
		return "", errors.New("role not found in context")
	}
	role, ok := val.(string)
	if !ok {
		return "", errors.New("role has invalid type")
	}
	return role, nil
}
