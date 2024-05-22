package route

import (
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"

	"github.com/IBM/sarama"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/internal/event"
	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/LXJ0000/go-backend/pkg/orm"

	"github.com/LXJ0000/go-backend/api/middleware"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/gin-gonic/gin"
)

func Setup(env *bootstrap.Env, timeout time.Duration,
	db orm.Database, redisCache cache.RedisCache, localCache cache.LocalCache,
	server *gin.Engine,
	producer event.Producer, saramaClient sarama.Client) {

	server.Static(env.UrlStaticPath, env.LocalStaticPath)
	server.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "success")
	})
	server.POST("/upload", func(c *gin.Context) {
		// 单个文件
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, domain.ErrorResp("Upload File Fail", err))
			return
		}

		file.Filename = fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))
		dst := filepath.Join("assets/file", file.Filename)
		// 上传文件到指定的目录
		if err := c.SaveUploadedFile(file, dst); err != nil {
			slog.Error("save file fail", "err", err.Error())
			c.JSON(http.StatusInternalServerError, domain.ErrorResp("Upload File Fail", err))
			return
		}
		c.JSON(http.StatusOK, domain.SuccessResp(fmt.Sprintf("'%s' uploaded!", file.Filename)))
	})

	publicRouter := server.Group("/api")
	// All Public APIs
	NewSignupRouter(env, timeout, db, publicRouter)
	NewLoginRouter(env, timeout, db, publicRouter)
	NewRefreshTokenRouter(env, timeout, db, publicRouter)

	protectedRouter := server.Group("/api")
	// Middleware to verify AccessToken
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	// All Private APIs
	// User
	NewUserRouter(env, timeout, db, protectedRouter)
	// Task
	NewTaskRouter(env, timeout, db, protectedRouter)
	// Post
	NewPostRouter(env, timeout, db, redisCache, localCache, protectedRouter, producer, saramaClient)
	// Comment
	NewCommentRouter(env, timeout, db, protectedRouter)
	// Ralation
	NewRelationRouter(env, timeout, db, redisCache, protectedRouter)
}
