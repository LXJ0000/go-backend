package main

import (
	"time"

	route "github.com/LXJ0000/go-backend/api/route"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/gin-gonic/gin"
)

func main() {

	app := bootstrap.App()

	env := app.Env

	//db := app.Mongo.Database(env.DBName)
	//defer app.CloseDBConnection()

	db := app.Orm
	cache := app.Cache

	timeout := time.Duration(env.ContextTimeout) * time.Hour // TODO

	server := gin.Default()

	//route.Setup(env, timeout, db, server, orm)
	route.Setup(env, timeout, db, cache, server)

	_ = server.Run(env.ServerAddress)
}
