package route

import (
	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/orm"
	"github.com/LXJ0000/go-backend/repository"
	"github.com/LXJ0000/go-backend/usecase"
	"github.com/gin-gonic/gin"
	"time"
)

func NewPostRouter(env *bootstrap.Env, timeout time.Duration, orm orm.Database, group *gin.RouterGroup) {
	repo := repository.NewPostRepository(orm)
	col := &controller.PostController{
		PostUsecase: usecase.NewPostUsecase(repo, timeout),
	}
	group.POST("/post", col.CreateOrPublish)
	group.GET("/post", col.Info)
	group.GET("/post/publish", col.ReaderList)
	group.GET("/post/private", col.WriterList)
}
