package route

import (
	"github.com/IBM/sarama"
	"github.com/LXJ0000/go-backend/cache"
	"github.com/LXJ0000/go-backend/event"
	"github.com/LXJ0000/go-backend/orm"
	"time"

	"github.com/LXJ0000/go-backend/api/middleware"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/gin-gonic/gin"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db orm.Database, cache cache.Cache, gin *gin.Engine,
	producer event.Producer, saramaClient sarama.Client) {

	publicRouter := gin.Group("/api")
	// All Public APIs
	NewSignupRouter(env, timeout, db, publicRouter)
	NewLoginRouter(env, timeout, db, publicRouter)
	NewRefreshTokenRouter(env, timeout, db, publicRouter)

	protectedRouter := gin.Group("/api")
	// Middleware to verify AccessToken
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	// All Private APIs
	// User
	NewUserRouter(env, timeout, db, protectedRouter)
	// Task
	NewTaskRouter(env, timeout, db, protectedRouter)
	// Post
	NewPostRouter(env, timeout, db, cache, protectedRouter, producer, saramaClient)
}
