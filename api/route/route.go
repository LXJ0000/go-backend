package route

import (
	"time"

	"github.com/IBM/sarama"
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
	// File
	NewFileRouter(env, timeout, db, protectedRouter)
}
