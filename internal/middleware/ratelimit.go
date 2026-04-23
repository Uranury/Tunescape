package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gitlab.com/Uranury/tunescape/pkg/apperrors"
)

type RateLimiter struct {
	rds       *redis.Client
	limit     int64
	window    time.Duration
	keyPrefix string
}

func NewRateLimiter(rds *redis.Client, limit int64, window time.Duration) *RateLimiter {
	return &RateLimiter{
		rds:       rds,
		limit:     limit,
		window:    window,
		keyPrefix: "ratelimit:",
	}
}

func (rl *RateLimiter) PerIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("%s%s", rl.keyPrefix, ip)
		now := time.Now()
		nowMs := now.UnixMilli()
		windowMs := rl.window.Milliseconds()
		member := fmt.Sprintf("%d", now.UnixNano())

		ctx := c.Request.Context()
		allowed, err := rl.rds.Eval(ctx, `
local key = KEYS[1]
local nowMs = tonumber(ARGV[1])
local windowMs = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local member = ARGV[4]
local minScore = nowMs - windowMs

redis.call("ZREMRANGEBYSCORE", key, "-inf", minScore)
local count = redis.call("ZCARD", key)
if count >= limit then
  redis.call("PEXPIRE", key, windowMs)
  return 0
end

redis.call("ZADD", key, nowMs, member)
redis.call("PEXPIRE", key, windowMs)
return 1
`, []string{key}, nowMs, windowMs, rl.limit, member).Int()
		if err != nil {
			apperrors.GenHTTPError(c, http.StatusServiceUnavailable, "rate limiter unavailable", nil)
			c.Abort()
			return
		}

		if allowed == 0 {
			apperrors.GenHTTPError(c, http.StatusTooManyRequests, "rate limit exceeded", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
