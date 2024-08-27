package middleware

import (
	"time"

	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-lib/gin-plugin/ratelimit"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(env *bootstrap.Env) gin.HandlerFunc {
	return ratelimit.New(
		// ratelimit.DefaultConfig(),
		ratelimit.Config{
			RedisAddr:     env.RedisAddr,
			Window:        time.Second,
			Limit:         100,
			RedisPassword: env.RedisPassword,
		},
	)
}
