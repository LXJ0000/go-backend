package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/LXJ0000/go-backend/utils/md5util"
	"github.com/gin-gonic/gin"
)

const (
	ApiCacheKeyPrefix = "api_cache:"
)

type cacheWrappedWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (rw *cacheWrappedWriter) Write(body []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(body)
	if err == nil {
		rw.body.Write(body)
	}
	return n, err
}

// NewAPICacheMiddleware 创建一个 API 缓存中间件 用于缓存 API 请求结果 以 uri + method + reqBody 作为 key
// 还没有自测 慎用
func NewAPICacheMiddleware(cache cache.RedisCache) func(timeout time.Duration) gin.HandlerFunc {
	return func(timeout time.Duration) gin.HandlerFunc {
		return func(c *gin.Context) {
			reqBody, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody)) // 回写 c.Request.Body
			// 1 uri + method + reqBody(md5) 作为 key
			rawKey := c.Request.RequestURI + c.Request.Method + string(reqBody)
			key := ApiCacheKeyPrefix + md5util.Md5([]byte(rawKey))
			// 2 从缓存中获取数据 get key
			ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
			value, _ := cache.Get(ctx, key)
			// 3 如果有数据则直接返回
			if value != "" {
				var resp domain.Response
				_ = json.Unmarshal([]byte(value), &resp)
				c.AbortWithStatusJSON(http.StatusOK, resp)
				cancel()
				return
			}
			cancel()
			// 4 如果没有数据则执行下面的代码
			temp := c.Writer
			w := &cacheWrappedWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
			c.Writer = w

			c.Next()
			// 5 执行完毕后将结果写入缓存 setnx
			c.Writer = temp
			ctx, cancel = context.WithTimeout(c.Request.Context(), time.Second)
			defer cancel()
			cache.SetNx(ctx, key, w.body.String(), timeout)
		}
	}
}
