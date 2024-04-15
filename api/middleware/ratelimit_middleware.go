package middleware

import (
	"fmt"
	"github.com/LXJ0000/go-backend/redis"
	"github.com/LXJ0000/go-backend/script"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

func RateLimitMiddleware(cache redis.Cache, window time.Duration, rate int) gin.HandlerFunc {
	return func(c *gin.Context) {
		isLimit, err := cache.LuaWithReturnBool(c, script.LuaSliceWindow, []string{fmt.Sprintf("ip:%s", c.ClientIP())}, window.Milliseconds(), rate, time.Now().UnixMilli())
		if err != nil {
			slog.Warn("Redis 操作失败 limit", "error", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if isLimit {
			slog.Warn("To many requests from ", "ip", c.ClientIP())
			c.AbortWithStatus(http.StatusServiceUnavailable)
			return
		}
		c.Next()
	}
}
