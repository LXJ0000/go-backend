package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/gin-gonic/gin"
)

func NewAPICacheMiddleware(cache cache.RedisCache) func(timeout time.Duration) gin.HandlerFunc {
	return func(timeout time.Duration) gin.HandlerFunc {
		return func(c *gin.Context) {
			reqBody, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody)) // 回写 c.Request.Body
			// TODO
			// 1 uri + method + reqBody(md5) 作为 key
			// 2 从缓存中获取数据 get key
			// 3 如果有数据则直接返回
			// 4 如果没有数据则执行下面的代码
			// 5 执行完毕后将结果写入缓存 setnx
		}
	}
}
