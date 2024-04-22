package main

import (
	"github.com/LXJ0000/go-backend/api/middleware"
	"time"

	route "github.com/LXJ0000/go-backend/api/route"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/gin-gonic/gin"
)

func main() {

	app := bootstrap.App()

	env := app.Env

	db := app.Orm
	cache := app.Cache

	producer := app.Producer
	saramaClient := app.SaramaClient

	timeout := time.Duration(env.ContextTimeout) * time.Second // TODO

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())
	server.Use(middleware.RateLimitMiddleware())
	server.Use(middleware.PrometheusMiddleware())
	route.Setup(env, timeout, db, cache, server, producer, saramaClient)

	_ = server.Run(env.ServerAddress)
}
