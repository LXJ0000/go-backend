package main

import (
	"log"
	"log/slog"
	"time"

	"github.com/LXJ0000/go-backend/api/middleware"
	"github.com/LXJ0000/go-backend/api/route"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/internal/event"
	"github.com/LXJ0000/go-backend/internal/repository"
	"github.com/LXJ0000/go-backend/internal/usecase"
	"github.com/LXJ0000/go-backend/internal/usecase/sms/aliyun"
	"github.com/LXJ0000/go-backend/internal/usecase/sms/local"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	feedUsecase "github.com/LXJ0000/go-backend/internal/usecase/feed"

	feedUsecaseHandler "github.com/LXJ0000/go-backend/internal/usecase/feed/handler"
)

func main() {
	_ = godotenv.Load()

	app := bootstrap.App()
	env := app.Env
	db := app.Orm
	redisCache := app.Cache
	localCache := app.LocalCache
	producer := app.Producer
	saramaClient := app.SaramaClient
	smsClient := app.SMSAliyunClient
	minioClient := app.MinioClient
	doubaoChat := app.DoubaoChat
	timeout := time.Duration(env.ContextTimeout) * time.Minute // 接口超时时间

	server := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// Middleware
	server.Use(middleware.CORSMiddleware())
	server.Use(middleware.RateLimitMiddleware(env))
	server.Use(middleware.PrometheusMiddleware())
	apiCache := middleware.NewAPICacheMiddleware(localCache)
	jwtAuth := middleware.JwtAuthMiddleware(app.Env.AccessTokenSecret, redisCache)

	// wire 复用对象
	// Repository
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
	// Usecase
	codeUc := usecase.NewCodeUsecase(codeRepo, aliyun.NewService(env.SMSAppID, env.SMSSignName, smsClient))
	commentUc := usecase.NewCommentUsecase(commentRepo, timeout)
	fileUc := usecase.NewFileUsecase(fileRepo, timeout, env.LocalStaticPath, env.UrlStaticPath, minioClient)
	interactionUc := usecase.NewInteractionUsecase(interactionRepo, timeout)
	postRankUc := usecase.NewPostRankUsecase(timeout, interactionRepo, postRepo, postRankRepo, doubaoChat)
	postUc := usecase.NewPostUsecase(postRepo, timeout, producer, postRankUc, doubaoChat)
	relationUc := usecase.NewRelationUsecase(relationRepo, userRepo, timeout)
	refreshTokenUc := usecase.NewRefreshTokenUsecase(userRepo, timeout)
	tagUc := usecase.NewTagUsecase(tagRepo, timeout)
	taskUc := usecase.NewTaskUsecase(taskRepo, timeout)
	userUc := usecase.NewUserUsecase(userRepo, timeout)
	sync2OpenIMUc := usecase.NewSync2OpenIMUsecase(env.OpenIMServerDoamin)
	// Feed handler and usecase
	feedLikeHdl := feedUsecaseHandler.NewFeedLikeHandler(feedRepo)
	feedPostHandler := feedUsecaseHandler.NewFeedPostHandler(feedRepo, relationUc)
	feedFollowHandler := feedUsecaseHandler.NewFeedFollowHandler(feedRepo)
	handlerMap := map[string]domain.FeedHandler{
		domain.FeedLikeEvent:   feedLikeHdl,
		domain.FeedPostEvent:   feedPostHandler,
		domain.FeedFollowEvent: feedFollowHandler,
	}
	feedUc := feedUsecase.NewFeedUsecase(handlerMap, relationUc, feedRepo)

	// producer and consumer
	consumer := event.NewBatchSyncReadEventConsumer(saramaClient, interactionRepo)
	if saramaClient != nil {
		if err := consumer.Start(); err != nil {
			slog.Error("OMG！消费者启动失败")
			log.Fatal(err)
		}
	}
	// local code service to debug
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

	// 组装对象
	obj := &route.App{
		Env:            env,
		CodeUc:         codeUc,
		CommentUc:      commentUc,
		FeedUc:         feedUc,
		FileUs:         fileUc,
		InterUc:        interactionUc,
		LocalCodeUc:    localCodeUc,
		PostUc:         postUc,
		RefreshTokenUc: refreshTokenUc,
		RelationUc:     relationUc,
		Sync2OpenIMUc:  sync2OpenIMUc,
		TaskUc:         taskUc,
		TagUc:          tagUc,
		UserUc:         userUc,
		JwtAuth:        jwtAuth,
		ApiCache:       apiCache,
	}

	route.Setup(server, obj)

	if err := server.Run(env.ServerAddr); err != nil {
		log.Fatal(err)
	}
}
