package route

import (
	"log"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/internal/event"
	"github.com/LXJ0000/go-backend/internal/repository"
	"github.com/LXJ0000/go-backend/internal/usecase"
	feedUsecase "github.com/LXJ0000/go-backend/internal/usecase/feed"
	feedUsecaseHandler "github.com/LXJ0000/go-backend/internal/usecase/feed/handler"
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
	NewSignupRouter(env, timeout, db, redisCache, publicRouter)
	NewLoginRouter(env, timeout, db, redisCache, publicRouter)
	NewRefreshTokenRouter(env, timeout, db, redisCache, publicRouter)

	protectedRouter := server.Group("/api")
	// Middleware to verify AccessToken
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret, redisCache))
	// All Private APIs

	// 复用对象
	userRepo := repository.NewUserRepository(db, redisCache)
	relationRepo := repository.NewRelationRepository(db)
	feedRepo := repository.NewFeedRepository(db)
	postRepo := repository.NewPostRepository(db, redisCache)
	interactionRepo := repository.NewInteractionRepository(db, redisCache)
	postRankRepo := repository.NewPostRankRepository(localCache, redisCache)

	interactionUc := usecase.NewInteractionUsecase(interactionRepo, timeout)
	postRankUc := usecase.NewPostRankUsecase(interactionRepo, postRepo, postRankRepo, timeout)
	postUc := usecase.NewPostUsecase(postRepo, timeout, producer, postRankUc)
	relationUc := usecase.NewRelationUsecase(relationRepo, userRepo, timeout)
	userUc := usecase.NewUserUsecase(userRepo, timeout)

	feedLikeHdl := feedUsecaseHandler.NewFeedLikeHandler(feedRepo)
	feedPostHandler := feedUsecaseHandler.NewFeedPostHandler(feedRepo, relationUc)
	feedFollowHandler := feedUsecaseHandler.NewFeedFollowHandler(feedRepo)
	handlerMap := map[string]domain.FeedHandler{
		domain.FeedLikeEvent:   feedLikeHdl,
		domain.FeedPostEvent:   feedPostHandler,
		domain.FeedFollowEvent: feedFollowHandler,
	}

	feedUc := feedUsecase.NewFeedUsecase(handlerMap, relationUc, feedRepo)

	//生产消费
	consumer := event.NewBatchSyncReadEventConsumer(saramaClient, interactionRepo)
	if err := consumer.Start(); err != nil {
		slog.Error("OMG！消费者启动失败")
		log.Fatal(err)
	}

	// User
	NewUserRouter(env, timeout, db, redisCache, userUc, relationUc, postUc, protectedRouter)
	// Task
	NewTaskRouter(env, timeout, db, protectedRouter)
	// Post
	NewPostRouter(postUc, interactionUc, feedUc, userUc, protectedRouter)
	// Comment
	NewCommentRouter(env, timeout, db, protectedRouter)
	// Ralation
	NewRelationRouter(env, timeout, db, redisCache, protectedRouter)
	// File
	NewFileRouter(env, timeout, db, protectedRouter)
	// Tag
	NewTagRouter(env, timeout, db, redisCache, protectedRouter)
	// Feed
	NewFeedRouter(feedUc, protectedRouter)
}
