package main

import (
	"time"

	"github.com/LXJ0000/go-backend/api/middleware"
	"github.com/LXJ0000/go-backend/api/route"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	app := bootstrap.App()

	env := app.Env

	db := app.Orm
	cache := app.Cache
	localCache := app.LocalCache

	producer := app.Producer
	saramaClient := app.SaramaClient

	smsClient := app.SMSAliyunClient

	minioClient := app.MinioClient

	daobaoChat := app.DoubaoChat

	// TODO wire 前置到这里
	// cron := bootstrap.NewCron()
	// cron.Start()
	// defer func() {
	// 	// 优雅退出
	// 	ctx := cron.Stop()
	// 	<-ctx.Done()
	// }()

	timeout := time.Duration(env.ContextTimeout) * time.Minute // 接口超时时间

	server := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	server.Use(middleware.CORSMiddleware())
	server.Use(middleware.RateLimitMiddleware(env))
	server.Use(middleware.PrometheusMiddleware())
	apiCache := middleware.NewAPICacheMiddleware(cache)
	route.Setup(env, timeout, db, cache, localCache, server, producer, saramaClient, smsClient, minioClient, daobaoChat, apiCache)

	_ = server.Run(env.ServerAddr)
}
