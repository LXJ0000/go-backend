package middleware

import (
	"time"

	"github.com/LXJ0000/go-lib/gin-plugin/ratelimit"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware() gin.HandlerFunc {
	return ratelimit.New(
		// ratelimit.DefaultConfig(),
		ratelimit.Config{
			RedisAddr: "127.0.0.1:6379",
			Window:    time.Second,
			Limit:     1000,
		},
	)
}
