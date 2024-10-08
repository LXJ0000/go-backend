package main

import (
	"time"

	"github.com/LXJ0000/go-backend/api/middleware"
	"github.com/LXJ0000/go-backend/api/route"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/gin-gonic/gin"
)

func main() {

	app := bootstrap.App()

	env := app.Env

	db := app.Orm
	cache := app.Cache
	localCache := app.LocalCache

	producer := app.Producer
	saramaClient := app.SaramaClient

	cron := app.Cron
	cron.Start()
	defer func() {
		// 优雅退出
		ctx := cron.Stop()
		<-ctx.Done()
	}()

	timeout := time.Duration(env.ContextTimeout) * time.Minute // 接口超时时间

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())
	server.Use(middleware.RateLimitMiddleware(env))
	server.Use(middleware.PrometheusMiddleware())
	route.Setup(env, timeout, db, cache, localCache, server, producer, saramaClient)

	_ = server.Run(env.ServerAddress)
}
