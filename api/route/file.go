package route

import (
	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/internal/repository"
	"github.com/LXJ0000/go-backend/internal/usecase"
	"github.com/LXJ0000/go-backend/pkg/orm"
	"github.com/gin-gonic/gin"
	"time"
)

func NewFileRouter(env *bootstrap.Env, timeout time.Duration,
	orm orm.Database, group *gin.RouterGroup) {
	repo := repository.NewFileRepository(orm)
	useCase := usecase.NewFileUsecase(repo, timeout, env.LocalStaticPath, env.UrlStaticPath)
	col := &controller.FileController{FileUsecase: useCase}
	group.POST("/file/upload", col.Upload)
	group.GET("/file/list", col.FileList)
}
