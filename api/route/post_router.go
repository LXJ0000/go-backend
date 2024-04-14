package route

import (
	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/orm"
	"github.com/LXJ0000/go-backend/redis"
	"github.com/LXJ0000/go-backend/repository"
	"github.com/LXJ0000/go-backend/usecase"
	"github.com/gin-gonic/gin"
	"time"
)

func NewPostRouter(env *bootstrap.Env, timeout time.Duration, orm orm.Database, cache redis.Cache, group *gin.RouterGroup) {
	repoPost := repository.NewPostRepository(orm, cache)
	repoInteraction := repository.NewInteractionRepository(orm, cache)
	col := &controller.PostController{
		PostUsecase:        usecase.NewPostUsecase(repoPost, timeout),
		InteractionUseCase: usecase.NewInteractionUsecase(repoInteraction, timeout),
	}
	group.POST("/post", col.CreateOrPublish)
	group.GET("/post", col.Info)
	group.GET("/post/publish", col.ReaderList)
	group.GET("/post/private", col.WriterList)
	group.POST("/post/like", col.Like)
}
