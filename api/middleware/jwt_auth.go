package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/pkg/cache"

	"github.com/LXJ0000/go-backend/utils/tokenutil"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(secret string, cache cache.RedisCache) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			authorized, err := tokenutil.IsAuthorized(authToken, secret)
			if authorized {
				userID, ssid, err := tokenutil.ExtractIDFromToken(authToken, secret)
				if err != nil {
					c.JSON(http.StatusUnauthorized, domain.ErrorResp("Not authorized", err))
					c.Abort()
					return
				}
				exist, err := cache.Exist(context.Background(), domain.UserLogoutKey(ssid))
				if err != nil {
					// TODO 降级策略
					c.JSON(http.StatusUnauthorized, domain.ErrorResp("Not authorized", err))
					c.Abort()
					return
				}
				if exist != 0 {
					c.JSON(http.StatusUnauthorized, domain.ErrorResp("Not authorized", err))
					c.Abort()
					return
				}
				c.Set(domain.XUserID, userID)
				c.Set(domain.UserSessionID, ssid)
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, domain.ErrorResp("Not authorized", err))
			c.Abort()
			return
		}
		c.JSON(http.StatusUnauthorized, domain.ErrorResp("Not authorized", nil))
		c.Abort()
	}
}
