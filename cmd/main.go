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

	timeout := time.Duration(env.ContextTimeout) * time.Second // TODO

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())
	server.Use(middleware.RateLimitMiddleware())
	route.Setup(env, timeout, db, cache, server)

	_ = server.Run(env.ServerAddress)
}
