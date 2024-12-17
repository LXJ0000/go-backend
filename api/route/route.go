package route

import (
	"log"
	"log/slog"
	"time"

	"github.com/LXJ0000/go-backend/internal/usecase/sms/aliyun"
	"github.com/LXJ0000/go-backend/internal/usecase/sms/local"
	sms "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"

	"github.com/IBM/sarama"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/internal/event"
	"github.com/LXJ0000/go-backend/internal/repository"
	"github.com/LXJ0000/go-backend/internal/usecase"
	feedUsecase "github.com/LXJ0000/go-backend/internal/usecase/feed"
	feedUsecaseHandler "github.com/LXJ0000/go-backend/internal/usecase/feed/handler"
	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/LXJ0000/go-backend/pkg/chat"
	"github.com/LXJ0000/go-backend/pkg/file"
	"github.com/LXJ0000/go-backend/pkg/orm"

	"github.com/LXJ0000/go-backend/api/middleware"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/gin-gonic/gin"
)

func Setup(env *bootstrap.Env, timeout time.Duration,
	db orm.Database, redisCache *cache.RedisCache, localCache *cache.RistrettoCache,
	server *gin.Engine,
	producer event.Producer, saramaClient sarama.Client,
	smsClient *sms.Client,
	minioClient file.FileStorage,
	doubaoChat chat.Chat,
	apiCache func(timeout time.Duration) gin.HandlerFunc,
) {

	server.Static(env.UrlStaticPath, env.LocalStaticPath)
	server.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "success")
	})

	publicRouter := server.Group("/api")
	// All Public APIs
	NewRefreshTokenRouter(env, timeout, db, redisCache, publicRouter)

	protectedRouter := server.Group("/api")
	// Middleware to verify AccessToken
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret, redisCache))
	// All Private APIs

	// 复用对象
	codeRepo := repository.NewCodeRepository(redisCache)
	commentRepo := repository.NewCommentRepository(db)
	feedRepo := repository.NewFeedRepository(db)
	fileRepo := repository.NewFileRepository(db)
	interactionRepo := repository.NewInteractionRepository(db, redisCache)
	postRankRepo := repository.NewPostRankRepository(localCache, redisCache)
	postRepo := repository.NewPostRepository(db, redisCache)
	relationRepo := repository.NewRelationRepository(db)
	tagRepo := repository.NewTagRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	userRepo := repository.NewUserRepository(db, redisCache)

	codeUc := usecase.NewCodeUsecase(codeRepo, aliyun.NewService(env.SMSAppID, env.SMSSignName, smsClient))
	commentUc := usecase.NewCommentUsecase(commentRepo, timeout)
	fileUc := usecase.NewFileUsecase(fileRepo, timeout, env.LocalStaticPath, env.UrlStaticPath, minioClient)
	interactionUc := usecase.NewInteractionUsecase(interactionRepo, timeout)
	postRankUc := usecase.NewPostRankUsecase(timeout, interactionRepo, postRepo, postRankRepo, doubaoChat)
	postUc := usecase.NewPostUsecase(postRepo, timeout, producer, postRankUc, doubaoChat)
	relationUc := usecase.NewRelationUsecase(relationRepo, userRepo, timeout)
	tagUc := usecase.NewTagUsecase(tagRepo, timeout)
	taskUc := usecase.NewTaskUsecase(taskRepo, timeout)
	userUc := usecase.NewUserUsecase(userRepo, timeout)
	sync2OpenIMUc := usecase.NewSync2OpenIMUsecase(env.OpenIMServerDoamin)

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
	if saramaClient != nil {
		if err := consumer.Start(); err != nil {
			slog.Error("OMG！消费者启动失败")
			log.Fatal(err)
		}
	}

	localCodeService := local.NewService()
	localCodeUc := usecase.NewCodeUsecase(codeRepo, localCodeService)

	// Cron
	cron := bootstrap.NewCron(timeout, postRankUc)
	cron.Start()
	defer func() {
		// 优雅退出
		ctx := cron.Stop()
		<-ctx.Done()
	}()

	// User
	NewUserRouter(env, userUc, relationUc, postUc, codeUc, localCodeUc, sync2OpenIMUc, fileUc, publicRouter, protectedRouter) // TODO 替换成 codeUc
	// Task
	NewTaskRouter(env, taskUc, protectedRouter)
	// Post
	NewPostRouter(postUc, interactionUc, feedUc, userUc, commentUc, fileUc, apiCache, protectedRouter)
	// Comment
	NewCommentRouter(env, commentUc, userUc, interactionUc, protectedRouter)
	// Relation
	NewRelationRouter(env, relationUc, protectedRouter)
	// File
	NewFileRouter(env, fileUc, protectedRouter)
	// Tag
	NewTagRouter(env, tagUc, protectedRouter)
	// Feed
	NewFeedRouter(feedUc, protectedRouter)
	// Intr
	NewIntrRouter(interactionUc, feedUc, protectedRouter)
}
