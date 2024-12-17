package route

import (
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"

	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/gin-gonic/gin"
)

type App struct {
	Env *bootstrap.Env

	CodeUc         domain.CodeUsecase
	CommentUc      domain.CommentUsecase
	FeedUc         domain.FeedUsecase
	FileUs         domain.FileUsecase
	InterUc        domain.InteractionUseCase
	LocalCodeUc    domain.CodeUsecase
	PostUc         domain.PostUsecase
	RefreshTokenUc domain.RefreshTokenUsecase
	RelationUc     domain.RelationUsecase
	Sync2OpenIMUc  domain.Sync2OpenIMUsecase
	TaskUc         domain.TaskUsecase
	TagUc          domain.TagUsecase
	UserUc         domain.UserUsecase

	JwtAuth  gin.HandlerFunc
	ApiCache func(timeout time.Duration) gin.HandlerFunc
}

func Setup(server *gin.Engine, app *App) {

	server.Static(app.Env.UrlStaticPath, app.Env.LocalStaticPath)
	server.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "success")
	})

	// Public APIs
	publicRouter := server.Group("/api")
	// Private APIs
	protectedRouter := server.Group("/api")
	protectedRouter.Use(app.JwtAuth)

	// User
	NewUserRouter(app.Env, app.UserUc, app.RelationUc, app.PostUc, app.CodeUc, app.LocalCodeUc, app.Sync2OpenIMUc, app.FileUs, publicRouter, protectedRouter) // TODO 替换成 codeUc
	// Task
	NewTaskRouter(app.Env, app.TaskUc, protectedRouter)
	// Post
	NewPostRouter(app.PostUc, app.InterUc, app.FeedUc, app.UserUc, app.CommentUc, app.FileUs, app.ApiCache, protectedRouter)
	// Comment
	NewCommentRouter(app.Env, app.CommentUc, app.UserUc, app.InterUc, protectedRouter)
	// Relation
	NewRelationRouter(app.Env, app.RelationUc, protectedRouter)
	// File
	NewFileRouter(app.Env, app.FileUs, protectedRouter)
	// Tag
	NewTagRouter(app.Env, app.TagUc, protectedRouter)
	// Feed
	NewFeedRouter(app.FeedUc, protectedRouter)
	// Intr
	NewIntrRouter(app.InterUc, app.FeedUc, protectedRouter)
	// RefreshToken
	NewRefreshTokenRouter(app.Env, app.UserUc, app.RefreshTokenUc, publicRouter)

}
